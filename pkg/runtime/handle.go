package runtime

import (
	"flag"
	"os"
	"strings"

	"github.com/wqvoon/cbox/pkg/container"
	"github.com/wqvoon/cbox/pkg/flags"
	"github.com/wqvoon/cbox/pkg/log"
	"github.com/wqvoon/cbox/pkg/rootdir"
	"github.com/wqvoon/cbox/pkg/runtime/info"
	"github.com/wqvoon/cbox/pkg/utils"
	"golang.org/x/sys/unix"
)

// 这个函数只能在 RuntimeMode 下使用
func Handle() {
	c := container.GetContainerByID(flag.Arg(0))

	// 这里就不 Error 了，仅做个提醒，也不是啥大事
	if err := unix.Sethostname([]byte(c.ID)); err != nil {
		log.Println("faild to set hostname, err:", err)
	}

	// TODO: 待补充其他的 ns
	namespaces := []string{"/pid", "/uts", "/ipc"}
	srcPathPrefix := "/proc/self/ns"
	dstPathPrefix := rootdir.GetContainerNSPath(c.ID)
	for _, ns := range namespaces {
		src := srcPathPrefix + ns
		dst := dstPathPrefix + ns
		if err := unix.Mount(src, dst, "", unix.MS_BIND, ""); err != nil {
			log.Errorf("failed to bind mount %q to %q, err: %v\n", src, dst, err)
		}
	}

	os.Clearenv()
	for _, oneEnv := range c.Env {
		envPair := strings.Split(oneEnv, "=")
		key, val := envPair[0], envPair[1]
		os.Setenv(key, val)
	}

	containerInfo := info.GetContainerInfo(c.ID)
	dnsFilePath := flags.GetDNSFilePath()
	if dnsFilePath != "" {
		utils.CopyFile(dnsFilePath, rootdir.GetContainerDNSConfigPath(c.ID))
		containerInfo.SaveDNSFilePath(dnsFilePath)
	}

	unix.Chroot(rootdir.GetContainerMountPath(c.ID))

	// TODO: Mount 的第一个参数如果留空则宿主机上会因为解析错误而读不到这条记录，也许可以利用下
	utils.CreateDirIfNotExist("/proc")
	if err := unix.Mount("cbox-proc", "/proc", "proc", 0, ""); err != nil {
		log.Errorln("faild to mount /proc, err:", err)
	}

	// TODO: 暂时使用 for + pause 的方式来减少资源消耗，后面尝试其他方法
	for {
		unix.Pause()
	}
}
