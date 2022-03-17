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

		SysProcAttr: &unix.SysProcAttr{
			Cloneflags: unix.CLONE_NEWPID |
				unix.CLONE_NEWNS |
				unix.CLONE_NEWUTS |
				unix.CLONE_NEWIPC,
		},
	}

	if err := cmd.Run(); err != nil {
		log.Errorln("faild to run runtime, err:", err)
	}
}
