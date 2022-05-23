package runtime

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/wqvoon/cbox/pkg/container"
	"github.com/wqvoon/cbox/pkg/log"
	"github.com/wqvoon/cbox/pkg/network/dns"
	"github.com/wqvoon/cbox/pkg/rootdir"
	"github.com/wqvoon/cbox/pkg/runtime/info"
	runtimeUtils "github.com/wqvoon/cbox/pkg/runtime/utils"
	"github.com/wqvoon/cbox/pkg/utils"
	"golang.org/x/sys/unix"
)

// 这个函数只能在 RuntimeMode 下使用
func Handle() {
	c := container.GetContainerByID(flag.Arg(0))

	// 创建容器对应的子 CGroup，并将当前pid加入
	runtimeUtils.SetupCGroup(c.ID)
	runtimeUtils.JoinCurrentProcessToCGroup(os.Getpid(), c.ID)

	// 这里就不 Error 了，仅做个提醒，也不是啥大事
	if err := unix.Sethostname([]byte(c.ID)); err != nil {
		log.Println("faild to set hostname, err:", err)
	}

	runtimeUtils.UpdateEnv(c.Env)

	containerInfo := info.GetContainerInfo(c.ID)

	{ // 配置容器 dns，如果 /etc 目录存在但 /etc/resolv.conf 文件不存在，那么创建该文件
		hostDnsFilePath := dns.GetDNSFilePath()
		containerDnsFilePath := rootdir.GetContainerDNSConfigPath(c.ID)
		containerEtcPath := filepath.Dir(containerDnsFilePath)
		if !utils.PathIsExist(containerDnsFilePath) && utils.PathIsExist(containerEtcPath) {
			utils.CopyFile(hostDnsFilePath, containerDnsFilePath)
			containerInfo.SaveDNSFilePath(hostDnsFilePath)
		}
	}

	for _, v := range containerInfo.Volumes {
		v.Mount()
	}

	if err := unix.Chroot(rootdir.GetContainerMountPath(c.ID)); err != nil {
		log.Errorln("failed to chroot, err:", err)
	}

	if err := unix.Chdir("/"); err != nil {
		log.Errorln("failed to chdir, err:", err)
	}

	// TODO: Mount 的第一个参数如果留空则宿主机上会因为解析错误而读不到这条记录，也许可以利用下
	utils.CreateDirIfNotExist("/proc")
	if err := unix.Mount("cbox-proc", "/proc", "proc", 0, ""); err != nil {
		log.Errorln("faild to mount /proc, err:", err)
	}

	enterLoop(c)
}

// 让 runtime 进入循环，从而保持后台运行的状态，根据是否有 healthCheckTask，会进入健康检查循环或 pause
func enterLoop(c *container.Container) {
	healthCheckTask := c.Image.HealthCheckTask
	filePath := rootdir.GetContainerHealthCheckInfoPath(c.ID, true)

	if healthCheckTask != nil && healthCheckTask.IsValid() {
		log.Println("runtime start to check health")
		healthCheckTask.Start(func([]byte) {
			// 成功时执行的回调
			if utils.PathIsExist(filePath) {
				os.Remove(filePath)
			}
		}, func(e error, info []byte) {
			// 失败时执行的回调
			var err error
			defer func() {
				if err != nil {
					c.Stop()
					log.Errorln("failed to write unhealthy reason, err:", err)
				}
			}()

			file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
			if err != nil {
				return
			}
			defer file.Close()

			_, err = file.Write(info)
			if err != nil {
				return
			}
		})
	} else {
		log.Println("no valid health check task, so runtime just pause")
		for {
			unix.Pause()
		}
	}
}
