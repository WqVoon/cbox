package network

import (
	"github.com/wqvoon/cbox/pkg/log"
	"golang.org/x/sys/unix"
)

// 创建一个新的网络命名空间并返回其 fd，当前进程依然留在原命名空间中
// 原理在于切换当前进程的命名空间并不会影响已打开的文件描述符
// 所以当前进程就可以根据这些打开的描述符在不同命名空间中反复横跳
// fd 的关闭由调用方接管
func CreateNamespace() int {
	sourceNS, err := unix.Open("/proc/self/ns/net", unix.O_RDONLY, 0)
	if err != nil {
		log.Errorln("failed to get source namespace, err:", err)
	}

	unix.Unshare(unix.CLONE_NEWNET)

	targetNS, err := unix.Open("/proc/self/ns/net", unix.O_RDONLY, 0)
	if err != nil {
		log.Errorln("failed to get target namespace, err:", err)
	}

	if err := unix.Setns(sourceNS, unix.CLONE_NEWNET); err != nil {
		log.Errorln("failed to go back source namespace, err:", err)
	}

	return targetNS
}

func EnterNamespaceByFd(fd int) {
	if err := unix.Setns(fd, unix.CLONE_NEWNET); err != nil {
		log.Errorln("failed to enter ns, err:", err)
	}
}
