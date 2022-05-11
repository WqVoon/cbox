package utils

import (
	"fmt"
	"os"
	"syscall"

	"github.com/wqvoon/cbox/pkg/log"
	"golang.org/x/sys/unix"
)

// 检查进程是否存在
func ProcessIsRunning(pid int) bool {
	// unix 系统下一定不会返回错误，所以不检查
	p, _ := os.FindProcess(pid)

	err := p.Signal(syscall.Signal(0))
	return err != os.ErrProcessDone
}

// 进入 pid 对应的进程的命名空间中
func EnterNamespaceByPid(pid int) {
	if !ProcessIsRunning(pid) {
		log.Errorln("can not enter namespace, process", pid, "has finished")
	}

	nsPathPrefix := fmt.Sprintf("/proc/%d/ns/", pid)
	targetNs := map[string]int{
		"ipc": unix.CLONE_NEWIPC,
		"pid": unix.CLONE_NEWPID,
		"uts": unix.CLONE_NEWUTS,
	}

	for nsName, nsType := range targetNs {
		nsFullPath := nsPathPrefix + nsName
		fd, err := os.Open(nsFullPath)
		if err != nil {
			log.Errorf("failed to open namespace %s, err: %s", nsName, err)
		}

		if err := unix.Setns(int(fd.Fd()), nsType); err != nil {
			log.Errorln("failed to enter namespace %s, err: %s", nsName, err)
		}

		fd.Close()
	}
}
