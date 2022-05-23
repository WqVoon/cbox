package container

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/wqvoon/cbox/pkg/log"
	"github.com/wqvoon/cbox/pkg/rootdir"
	runtimeCmd "github.com/wqvoon/cbox/pkg/runtime/cmd"
	runtimeInfo "github.com/wqvoon/cbox/pkg/runtime/info"
	runtimeUtils "github.com/wqvoon/cbox/pkg/runtime/utils"
	"github.com/wqvoon/cbox/pkg/storage/driver"
	"github.com/wqvoon/cbox/pkg/utils"
	"golang.org/x/sys/unix"
)

func (c *Container) Start() {
	info := runtimeInfo.GetContainerInfo(c.ID)
	{
		if info.IsRunning() {
			log.Errorln("can not start a running container")
			return
		}
		if info.GetStorageDriver() != driver.D {
			log.Errorf("can not use %s to start %s", driver.D, c.Name)
		}
	}

	if os.Geteuid() != 0 {
		log.Errorln("only root user can start a container")
	}

	containerMntPoint := rootdir.GetContainerMountPath(c.ID)
	driver.D.Mount(containerMntPoint, c.Image.Layers...)

	runtimeCmd.Run(c.ID)

	log.Printf("container %s started\n", c.Name)
}

func (c *Container) Exec(input ...string) {
	containerInfo := runtimeInfo.GetContainerInfo(c.ID)
	if !containerInfo.IsRunning() {
		log.Errorln("can not exec a not-running container")
		return
	}

	if os.Geteuid() != 0 {
		log.Errorln("only root user can exec a container")
	}

	var name string
	var args []string
	if len(input) > 0 {
		name, args = utils.ParseCmd(input...)
	} else {
		name, args = utils.ParseCmd(c.Entrypoint...)
	}

	utils.EnterNamespaceByPid(containerInfo.Pid)

	// 需要保证在 exec.Command 前进行 Env 的处理，这样得到的 cmd 才是正确的
	runtimeUtils.UpdateEnv(c.Env)

	runtimeUtils.JoinCurrentProcessToCGroup(os.Getpid(), c.ID)

	// 这里需要保证在 exec.Command 前进行 chroot，这样得到的 cmd 才是正确的
	unix.Chroot(rootdir.GetContainerMountPath(c.ID))

	cmd := exec.Command(name, args...)
	{
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Dir = "/"
		cmd.Env = c.Env
	}

	if err := cmd.Run(); err != nil {
		// 这里不进行 Errorln，用于避免 Ctrl+C 后 Ctrl+D 引起的常见错误
		log.Println("an error may have occurred while running container:", err)
	}
}

func (c *Container) Stop() {
	info := runtimeInfo.GetContainerInfo(c.ID)
	{
		if !info.IsRunning() {
			log.Errorln("can not stop a not-running container")
			return
		}
		if info.GetStorageDriver() != driver.D {
			log.Errorf("can not use %s to stop %s", driver.D, c.Name)
		}
	}

	runtimeCmd.Stop(info)

	mntPath := rootdir.GetContainerMountPath(info.ContainerID)
	driver.D.UnMount(mntPath)

	log.Printf("container %s stopped\n", c.Name)
}

func (c *Container) Delete() {
	if runtimeInfo.GetContainerInfo(c.ID).IsRunning() {
		log.Errorln("can not delete a running container")
		return
	}

	runtimeInfo.GetImageInfo(c.Image.Hash).MarkReleasedBy(c.ID)

	GetContainerIdx().DeleteByName(c.Name)
	if err := os.RemoveAll(c.rootPath); err != nil {
		log.Errorf("faild to remove container %q, err: %v\n", c.Name, err)
	}
	// TODO: 后面要处理更多的运行时副作用

	log.Printf("container %s deleted\n", c.Name)
}

// 从宿主机复制文件/文件夹到容器内，from 是宿主机路径，to 是容器路径
// from 可以是相对路径，to 需要是绝对路径
func (c *Container) CopyFromHost(from, to string) {
	if !runtimeInfo.GetContainerInfo(c.ID).IsRunning() {
		log.Errorln("can not copy file to a not running container")
	}

	if !utils.PathIsExist(from) {
		log.Errorln("path", from, "is not exists")
	}

	if !filepath.IsAbs(to) {
		log.Errorln("path", to, "is not abs path")
	}

	fullDstPath := rootdir.GetContainerMountPath(c.ID) + to

	utils.CopyDirContent(from, fullDstPath)

	log.Println("copy done")
}

// 从容器复制文件/文件夹到宿主机，from 是容器路径，to 是宿主机路径
// from 需要是绝对路径，to 可以是相对路径
func (c *Container) CopyToHost(from, to string) {
	if !runtimeInfo.GetContainerInfo(c.ID).IsRunning() {
		log.Errorln("can not copy file to a not running container")
	}

	if !filepath.IsAbs(from) {
		log.Errorln("path", from, "is not abs path")
	}

	fullSrcPath := rootdir.GetContainerMountPath(c.ID) + from

	if !utils.PathIsExist(fullSrcPath) {
		log.Errorf("path `%s` is not exists for container `%s`\n", from, c.Name)
	}

	utils.CopyDirContent(fullSrcPath, to)

	log.Println("copy done")
}

// 展示容器相关的详细信息
func (c *Container) Inspect() {
	info := runtimeInfo.GetContainerInfo(c.ID)

	log.Println("inspection of container", c.Name)
	{
		log.Println("- id:", c.ID)
		log.Println("- name:", c.Name)
		log.Println("- image:", c.Image.NameTag)
		log.Println("- is healthy:", info.IsHealthy())
		log.Println("- is running:", info.IsRunning())

		if info.IsRunning() {
			log.Println("- runtime pid:", info.Pid)
		}

		log.Println("- storage driver:", info.StorageDriver)
		log.Println("- dns file path:", info.DNSFilePath)
		log.Println("- container layout path:", c.rootPath)
		log.Println("- entrypoint:", c.Entrypoint)

		if len(c.Env) > 0 {
			log.Println("- env:")
			for idx, kvPair := range c.Env {
				splitedKvPair := strings.SplitN(kvPair, "=", 2)
				log.Printf("  - env #%d:\n", idx)
				log.Println("    - key:", splitedKvPair[0])
				log.Println("    - val:", splitedKvPair[1])
			}
		}

		if len(info.Volumes) > 0 {
			log.Println("- volumes:")
			containerPathPrefix := len(rootdir.GetContainerMountPath(c.ID))

			for idx, v := range info.Volumes {
				log.Printf("  - volume #%d:\n", idx)
				log.Println("    - host path:", v.HostPath)
				log.Println("    - container path:", v.ContainerPath[containerPathPrefix:])
			}
		}
	}
}
