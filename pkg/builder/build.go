package builder

import (
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/wqvoon/cbox/pkg/image"
	"github.com/wqvoon/cbox/pkg/log"
	"github.com/wqvoon/cbox/pkg/network/dns"
	"github.com/wqvoon/cbox/pkg/rootdir"
	"github.com/wqvoon/cbox/pkg/utils"
	"golang.org/x/sys/unix"
)

const (
	imageConfigFilename = "config.json"
	layerName           = "layer/"        // 仅支持单层镜像
	cmdsFileName        = "/cbox-cmds.sh" // BuildTask.Cmds 的文本形式，会被 execCmds 函数生成并执行
)

func (bt *BuildTask) Build() {
	log.Println("start build image")

	bt.preflight()
	bt.prepareBaseImage()
	bt.generateImageHash()
	bt.generateLayout()
	bt.generateTempDNSFile()
	bt.execCmds()
	bt.updateImageIdx()

	log.Println("build image", bt.ImageNameTag, "done")
}

// 对 BuildTask 的各个字段做合法性验证，并在合法的情况下为一些内部使用的字段赋值
func (bt *BuildTask) preflight() {
	// 解析失败的话会退出
	log.Println("check base image")
	bt.BaseImageNameTag = utils.GetNameTag(bt.BaseImageStr)

	log.Println("check image name")
	// TODO: 这里需要判重
	bt.ImageNameTag = utils.GetNameTag(bt.ImageNameStr)

	log.Println("check copy path")
	for _, copyPair := range bt.Copy {
		if !utils.PathIsExist(copyPair.Src) {
			log.Errorln("path", copyPair.Src, "is not exists")
		}

		if !filepath.IsAbs(copyPair.Dst) {
			log.Errorln("path", copyPair.Dst, "is not abs path")
		}
	}

	log.Println("check entrypoint")
	if len(bt.Entrypint) == 0 {
		log.Errorln("entrypoint must be defined")
	}

	log.Println("check health-check-task")
	if bt.HealthCheckTask != nil && !bt.HealthCheckTask.IsValid() {
		log.Errorln("invalid health-check-task")
	}
}

// 准备基础镜像，如果不存在则进行拉取
func (bt *BuildTask) prepareBaseImage() {
	nameTag := bt.BaseImageNameTag
	log.Println("prepare base image", nameTag)
	// image.Pull(nameTag)
	bt.BaseImage = image.GetImage(nameTag)
}

// 生成镜像对应的哈希
func (bt *BuildTask) generateImageHash() {
	const prefix = "cbox"
	randBytes := make([]byte, 4)
	rand.Read(randBytes)
	// TODO: 这里需要判重
	bt.ImageHash = fmt.Sprintf("%s%02x%02x%02x%02x",
		prefix,
		randBytes[0], randBytes[1],
		randBytes[2], randBytes[3])

	log.Println("generate image hash", bt.ImageHash)
}

// 生成镜像的 layout，包括 manifest、image config、fs
func (bt *BuildTask) generateLayout() {
	log.Println("generate image layout")
	imageLayoutpath := rootdir.GetImageLayoutPath(bt.ImageHash)
	utils.CreateDirWithExclusive(imageLayoutpath)

	log.Println("generate manifest file")
	manifest := image.ManifestList{
		{
			Config:   imageConfigFilename,
			RepoTags: []string{bt.ImageNameTag.String()},
			Layers:   []string{layerName},
		},
	}
	utils.SaveObjToJsonFile(rootdir.GetManifestPath(bt.ImageHash), manifest)

	log.Println("generate config file")
	config := &image.ImageConfig{
		Config: image.ImageConfigDetail{
			Env:             bt.parseEnv(),
			Cmd:             bt.Entrypint,
			HealthCheckTask: bt.HealthCheckTask,
		},
	}
	utils.SaveObjToJsonFile(rootdir.GetImageConfigPath(bt.ImageHash, imageConfigFilename), config)

	bt.generateFS()
}

// 生成镜像的 fs，先复制基础镜像的内容，再根据 Copy 字段复制宿主机的内容到镜像内
func (bt *BuildTask) generateFS() {
	log.Println("generate image fs from", bt.BaseImage.NameTag)

	if len(bt.BaseImage.Layers) != 1 {
		log.Errorln("cbox only supports single layer base image")
	}

	dstPath := rootdir.GetImageFsPath(bt.ImageHash, layerName)
	utils.CreateDirWithExclusive(path.Dir(dstPath))

	for _, layerPath := range bt.BaseImage.Layers {
		utils.CopyDirContent(layerPath, dstPath)
	}

	log.Println("generate image fs from copy field")
	for _, copyPair := range bt.Copy {
		utils.CopyDirContent(copyPair.Src, dstPath+copyPair.Dst)
	}
}

// 生成镜像的环境变量，该环境变量是基础镜像环境变量和 Env 字段的合并
// 如果有环境变量名重复的情况，那么优先使用 Env 字段中的内容
func (bt *BuildTask) parseEnv() []string {
	baseEnv := bt.BaseImage.Config.Config.Env

	envMap := make(map[string]string, len(baseEnv))
	for _, kv := range baseEnv {
		splitedKV := strings.Split(kv, "=")
		if len(splitedKV) != 2 {
			log.Errorln("error format of env", kv)
		}

		k, v := splitedKV[0], splitedKV[1]
		envMap[k] = v
	}

	for k, v := range bt.Env {
		envMap[k] = v
	}

	ret := make([]string, 0, len(envMap))
	for k, v := range envMap {
		ret = append(ret, k+"="+v)
	}
	return ret
}

// 创建一个临时的 /etc/resolv.conf 文件，因为 Cmds 中可能有网络下载的命令
func (bt *BuildTask) generateTempDNSFile() {
	if len(bt.Cmds) == 0 {
		return
	}

	dnsFilePath := dns.GetDNSFilePath()
	log.Println("generate /etc/resolv.conf file from", dnsFilePath)

	content, err := ioutil.ReadFile(dnsFilePath)
	if err != nil {
		log.Errorln("can not read content from", dnsFilePath)
	}

	fsPath := rootdir.GetImageFsPath(bt.ImageHash, layerName)
	utils.WriteFileIfNotExist(fsPath+"/etc/resolv.conf", content)
}

// 以镜像的 fs 作为根目录，执行 bt.Cmds 中的内容
// 当前实现为将 bt.Cmds 中的内容整理成一个 cbox-cmds.sh 文件
// 然后使用宿主机的 sh 来执行它，执行后将文件删除掉
func (bt *BuildTask) execCmds() {
	if len(bt.Cmds) == 0 {
		return
	}

	log.Println("exec image cmds")

	fsPath := rootdir.GetImageFsPath(bt.ImageHash, layerName)
	fullFilePath := fsPath + cmdsFileName

	{ // 创建并写入 cbox-cmds.sh 文件
		bt.Cmds = append(bt.Cmds, Cmd{"rm", cmdsFileName})
		oneLineCmds := make([]string, 0, len(bt.Cmds))
		for _, cmd := range bt.Cmds {
			cmdStr := fmt.Sprintf("chroot %s sh -c '%s' || exit 1",
				fsPath, strings.Join(cmd, " "))

			oneLineCmds = append(oneLineCmds, cmdStr)
		}
		fileContent := strings.Join(oneLineCmds, "\n")

		log.Println(fileContent)
		utils.WriteFileIfNotExist(fullFilePath, []byte(fileContent))
	}

	cmd := exec.Command("sh", fullFilePath)
	{
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Env = bt.parseEnv()
		cmd.SysProcAttr = &unix.SysProcAttr{
			Cloneflags: unix.CLONE_NEWPID |
				unix.CLONE_NEWUTS |
				unix.CLONE_NEWIPC,
		}
	}
	if err := cmd.Run(); err != nil {
		log.Errorln("failed to exec cmds, err:", err)
	}
}

// 将镜像名与哈希的对应关系写入 imageIdx 文件中
func (bt *BuildTask) updateImageIdx() {
	log.Println("update image idx")
	image.GetImageIdx().Update(bt.ImageNameTag, bt.ImageHash)
}
