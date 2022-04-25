package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/wqvoon/cbox/pkg/flags"
	"github.com/wqvoon/cbox/pkg/log"
	"github.com/wqvoon/cbox/pkg/rootdir"
	runtimeInfo "github.com/wqvoon/cbox/pkg/runtime/info"
	"golang.org/x/sys/unix"
)

// 启动 containerID 对应的容器
// 创建一个 runtime 进程来设置容器的运行时状态
func Run(containerID string) {
	exePath := "/proc/self/exe"
	rootdirFlag := fmt.Sprintf("--root_dir=%s", rootdir.GetRootDirPath())
	dnsFilePath := fmt.Sprintf("--dns_file_path=%s", flags.GetDNSFilePath())

	cmd := &exec.Cmd{
		Path: exePath,

		// 这里目前真正有效的只有 exePath 和 containerID，后面的内容只是帮助调试
		Args: []string{exePath, rootdirFlag, dnsFilePath, containerID, "/* cbox's runtime */"},

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

	// 由于 runtime 进程应该运行在后台，所以这里使用 Start
	if err := cmd.Start(); err != nil {
		log.Errorln("faild to start runtime, err:", err)
	}

	runtimeInfo.GetContainerInfo(containerID).SavePid(cmd.Process.Pid)
}
