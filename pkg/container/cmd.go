package container

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/wqvoon/cbox/pkg/cgroups"
	"github.com/wqvoon/cbox/pkg/config"
	"github.com/wqvoon/cbox/pkg/log"
	"github.com/wqvoon/cbox/pkg/network/address"
	"github.com/wqvoon/cbox/pkg/rootdir"
	runtimeCmd "github.com/wqvoon/cbox/pkg/runtime/cmd"
	runtimeInfo "github.com/wqvoon/cbox/pkg/runtime/info"
	runtimeUtils "github.com/wqvoon/cbox/pkg/runtime/utils"
	"github.com/wqvoon/cbox/pkg/storage/driver"
	"github.com/wqvoon/cbox/pkg/utils"
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

	// 这里应该直接退出，因为如果开启了 cgroup feature 但不能加入 pid cgroup，那么后面的工作无意义
	if !runtimeUtils.CanJoinTaskToPidCGroup(c.ID) {
		log.Errorln("can not exec container, err: task limit")
	}

	utils.EnterNamespaceByPid(containerInfo.Pid)

	cmd := getCmdForContainer(c, input...)
	{
		if err := cmd.Start(); err != nil {
			// 这里进行 Errorln，因为这里会触发 LookPathError 之类的错误，此时应该直接退出
			log.Errorln("failed to exec container, err:", err)
		}

		runtimeUtils.JoinProcessToCGroup(cmd.Process.Pid, c.ID)

		if err := cmd.Wait(); err != nil {
			// 这里不进行 Errorln，用于避免 Ctrl+C 后 Ctrl+D 引起的常见错误
			log.Println("an error may have occurred while running container:", err)
		}
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
	info := runtimeInfo.GetContainerInfo(c.ID)
	if info.IsRunning() {
		log.Errorln("can not delete a running container")
		return
	}

	runtimeInfo.GetImageInfo(c.Image.Hash).MarkReleasedBy(c.ID)

	GetContainerIdx().DeleteByName(c.Name)
	if err := os.RemoveAll(c.rootPath); err != nil {
		log.Errorf("faild to remove container %q, err: %v\n", c.Name, err)
	}

	address.ReleaseIPByString(info.IP)
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
	isRunning := info.IsRunning()

	log.Println("inspection of container", c.Name)
	{
		log.Println("- id:", c.ID)
		log.Println("- name:", c.Name)
		log.Println("- image:", c.Image.NameTag)
		log.Println("- is healthy:", info.IsHealthy())
		log.Println("- is running:", isRunning)

		if isRunning {
			log.Println("- runtime pid:", info.Pid)
		}

		if address.IsValidIPv4(info.IP) {
			log.Println("- ip:", info.IP)
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

		if isRunning && config.GetCgroupConfig().Enable {
			log.Println("- resources:")

			memCGroup := cgroups.Mem.GetOrCreateSubCGroup(c.ID)
			{
				log.Printf("  - memory limit: %d byte\n", memCGroup.GetMemLimit())
				log.Printf("  - memory usage: %d byte\n", memCGroup.GetMemUsage())
			}

			pidCGroup := cgroups.Pid.GetOrCreateSubCGroup(c.ID)
			{
				limitVal := pidCGroup.GetLowerTaskLimit(cgroups.Pid)

				if limitVal == cgroups.TaskNoLimit {
					log.Println("  - task limit: no limit")
				} else {
					log.Printf("  - task limit: %d\n", limitVal)
				}

				log.Printf("  - task count: %d\n", pidCGroup.GetCurrentTaskNum())
			}
		}
	}
}
