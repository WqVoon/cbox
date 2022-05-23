package utils

import (
	"os"
	"os/exec"

	"github.com/wqvoon/cbox/pkg/cgroups"
	"github.com/wqvoon/cbox/pkg/config"
	"github.com/wqvoon/cbox/pkg/flags"
	"github.com/wqvoon/cbox/pkg/log"
)

// var runtimeModePattern = regexp.MustCompile(`^/proc/\d+/exe$`)
const runtimeModePattern = "/proc/self/exe"

// TODO: 除了命令的模式外，还需要检测 pid、挂载点、命名空间等特征
func IsRuntimeMode() bool {
	return len(os.Args) > 1 && runtimeModePattern == os.Args[0]
}

// 从 os.Args 中提取出 exec.Cmd 所需的 path 及 args 参数
func ExtractCmdFromOSArgs() (path string, args []string) {
	switch length := len(os.Args); length {
	case 1, 2:
		log.Errorln("error length of runtimeMode args:", length)

	case 3:
		path, args = os.Args[2], os.Args[2:]

	default:
		path, args = os.Args[3], os.Args[3:]
	}

	var err error
	path, err = exec.LookPath(path)
	if err != nil {
		log.Errorln("failed to lookPath, err:", err)
	}

	return
}

// 按照 flags 来设置 containerID 对应的 CGroup
func SetupCGroup(containerID string) {
	if !config.GetCgroupConfig().Enable {
		return
	}

	cpuCGroup := cgroups.Cpu.GetOrCreateSubCGroup(containerID)
	cpuCGroup.SetCPULimit(flags.GetCPULimit())
	cpuCGroup.SetNotifyOnRelease(true)

	memCGroup := cgroups.Mem.GetOrCreateSubCGroup(containerID)
	memCGroup.SetMemLimit(flags.GetMemLimit())
	memCGroup.SetNotifyOnRelease(true)

	taskLimit := flags.GetTaskLimit()
	if taskLimit == -1 { // 如果等于 -1，说明用户没有做限制
		return
	}

	pidCGroup := cgroups.Pid.GetOrCreateSubCGroup(containerID)
	// 这里加 8 是因为要包括 runtime 进程与 exec 的宿主机进程各占 4 个线程
	pidCGroup.SetTaskLimit(taskLimit + 8)
	pidCGroup.SetNotifyOnRelease(true)
}

// 将 pid 对应的进程加入到 containerID 对应的 CGroup 中
func JoinCurrentProcessToCGroup(pid int, containerID string) {
	if !config.GetCgroupConfig().Enable {
		return
	}

	cpuCGroup := cgroups.Cpu.GetOrCreateSubCGroup(containerID)
	cpuCGroup.JoinProcessToSelf(pid)

	memCGroup := cgroups.Mem.GetOrCreateSubCGroup(containerID)
	memCGroup.JoinProcessToSelf(pid)

	pidCGroup := cgroups.Pid.GetOrCreateSubCGroup(containerID)
	pidCGroup.JoinProcessToSelf(pid)
}
