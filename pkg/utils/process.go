package utils

import (
	"fmt"
	"os"
	"syscall"
	"time"

	"github.com/wqvoon/cbox/pkg/log"
	"golang.org/x/sys/unix"
)

const waitInterval = 100 * time.Millisecond

// 检查进程是否存在
func ProcessIsRunning(pid int) bool {
	// unix 系统下一定不会返回错误，所以不检查
	p, _ := os.FindProcess(pid)

	err := p.Signal(syscall.Signal(0))
	return err != os.ErrProcessDone
}

// 等待进程启动
func WaitForStart(pid int) {
	for {
		if ProcessIsRunning(pid) {
			return
		}
		time.Sleep(waitInterval)
	}
}

// 等待进程退出
func WaitForStop(pid int) {
	for {
		if !ProcessIsRunning(pid) {
			return
		}
		time.Sleep(waitInterval)
	}
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
		"net": unix.CLONE_NEWNET,
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
