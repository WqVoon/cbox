package cmd

import (
	"os"
	"os/exec"

	"github.com/wqvoon/cbox/pkg/log"
	"golang.org/x/sys/unix"
)

func Run(containerID string, name string, args []string) {
	exePath := "/proc/self/exe"

	cmd := &exec.Cmd{
		Path: exePath,
		Args: append([]string{exePath, containerID, name}, args...),

		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,

		// 由于 golang 环境对 mount namespace 使用 setns 会报错，所以先不创建新的 ns
		// 依靠 chroot 来做隔离，宿主机依然可见容器内的挂载操作
		SysProcAttr: &unix.SysProcAttr{
			Cloneflags: unix.CLONE_NEWPID |
				unix.CLONE_NEWUTS |
				unix.CLONE_NEWIPC,
		},
	}

	if err := cmd.Run(); err != nil {
		log.Errorln("faild to run runtime, err:", err)
	}
}
