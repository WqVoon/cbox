package utils

import (
	"os"
	"syscall"
)

// 检查进程是否存在
func ProcessIsRunning(pid int) bool {
	// unix 系统下一定不会返回错误，所以不检查
	p, _ := os.FindProcess(pid)

	err := p.Signal(syscall.Signal(0))
	return err != os.ErrProcessDone
}
