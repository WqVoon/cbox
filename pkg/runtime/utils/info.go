package utils

import (
	"os"

	"github.com/wqvoon/cbox/pkg/log"
	"github.com/wqvoon/cbox/pkg/rootdir"
	"github.com/wqvoon/cbox/pkg/utils"
)

// 在容器停止时会将 containerInfo.Pid 设置为此值
const STOPPED_PID = -1

func GetContainerInfo(containerID string) *ContainerInfo {
	infoPath := rootdir.GetContainerInfoPath(containerID)
	infoObj := &ContainerInfo{}

	utils.GetObjFromJsonFile(infoPath, infoObj)
	infoObj.filePath = infoPath

	return infoObj
}

// 表示 Container 运行时的一些信息
type ContainerInfo struct {
	// info 文件相对于 rootdir 的路径
	filePath string

	// runtime 进程的 pid，被 runtime.Run 写入
	Pid int `json:"pid"`
}

// 获取 pid 对应的 Process 对象
func (ci *ContainerInfo) GetProcess() *os.Process {
	if ci.Pid == STOPPED_PID {
		log.Errorln("can not get process from stopped containers")
	}

	p, err := os.FindProcess(ci.Pid)
	if err != nil {
		log.Errorln("failed to get process from container info, pid:", ci.Pid)
	}

	return p
}

func (ci *ContainerInfo) SavePid(pid int) {
	ci.Pid = pid
	ci.save()
}

// 保存整个对象到 info 文件中
func (ci *ContainerInfo) save() {
	utils.SaveObjToJsonFile(ci.filePath, ci)
}
