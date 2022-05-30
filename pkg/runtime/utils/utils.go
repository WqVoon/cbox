package utils

import (
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"

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
	if taskLimit == cgroups.TaskNoLimit {
		return
	}

	pidCGroup := cgroups.Pid.GetOrCreateSubCGroup(containerID)
	pidCGroup.SetTaskLimit(taskLimit)
	pidCGroup.SetNotifyOnRelease(true)
}

// 是否可以加入一个 Task 到 name 对应的 CGroup 中，如果没有开启 cgroup feature，那么永远返回 true
func CanJoinTaskToPidCGroup(name string) bool {
	if !config.GetCgroupConfig().Enable {
		return true
	}

	pidCGroup := cgroups.Pid.GetOrCreateSubCGroup(name)
	return cgroups.Pid.CanJoinTask() && pidCGroup.CanJoinTask()
}

// 获取 containerID 对应容器的几个 cgroup.procs 文件
// 之所以要这样做，是因为若 `cbox exec` 想让容器的 PATH 环境变量生效就需要使用 chroot
// 而一旦使用 chroot，那么就无法直接调用 JoinProcessToSelf 方法
// 所以在 chroot 之前拿到几个 cgroup.procs 文件的句柄，然后后面直接写入进程 id 到句柄中
func GetCGroupProcsFiles(containerID string) []*os.File {
	if !config.GetCgroupConfig().Enable {
		return nil
	}

	fileName := "cgroup.procs"

	pidProcsFilePath := path.Join(cgroups.Pid.GetOrCreateSubCGroup(containerID).GetDirPath(), fileName)
	pidFile, err := os.OpenFile(pidProcsFilePath, os.O_WRONLY, 0700)
	if err != nil {
		log.Errorln("failed to open pidProcsFile, err:", err)
	}

	cpuProcsFilePath := path.Join(cgroups.Cpu.GetOrCreateSubCGroup(containerID).GetDirPath(), fileName)
	cpuFile, err := os.OpenFile(cpuProcsFilePath, os.O_WRONLY, 0700)
	if err != nil {
		log.Errorln("failed to open cpuProcsFile, err:", err)
	}

	memProcsFilePath := path.Join(cgroups.Mem.GetOrCreateSubCGroup(containerID).GetDirPath(), fileName)
	memFile, err := os.OpenFile(memProcsFilePath, os.O_WRONLY, 0700)
	if err != nil {
		log.Errorln("failed to open memProcsFile, err:", err)
	}

	return []*os.File{pidFile, cpuFile, memFile}
}

// 将 pid 对应的进程写入到 files 中，这里的 files 实际是通过 GetCGroupProcsFiles 获得的
func JoinProcessToCGroup(pidStr int, files []*os.File) {
	if !config.GetCgroupConfig().Enable {
		return
	}

	pidNum := strconv.Itoa(pidStr)

	for _, fd := range files {
		if _, err := fd.WriteString(pidNum); err != nil {
			log.Errorln("failed to write", pidStr, "to cgroups, err:", err)
		}
		fd.Close() // 写完后直接关闭即可
	}
}

// 为 containerID 对应的容器删除 CGroups
func DeleteCGroupForContainer(containerID string) {
	if !config.GetCgroupConfig().Enable {
		return
	}

	cgroups.Cpu.DeleteSubCGroup(containerID)
	cgroups.Mem.DeleteSubCGroup(containerID)
	cgroups.Pid.DeleteSubCGroup(containerID)
}

// 将当前进程的环境变量设置为 env 指定的内容，env 的格式如镜像的 config 文件中 ENV 字段的格式
func UpdateEnv(env []string) {
	os.Clearenv()
	for _, oneEnv := range env {
		envPair := strings.SplitN(oneEnv, "=", 2)
		key, val := envPair[0], envPair[1]
		os.Setenv(key, val)
	}
}
