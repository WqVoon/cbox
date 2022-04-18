package info

import (
	"os"
	"path"
	"strings"

	"github.com/wqvoon/cbox/pkg/log"
	"github.com/wqvoon/cbox/pkg/rootdir"
	"github.com/wqvoon/cbox/pkg/storage/driver"
	"github.com/wqvoon/cbox/pkg/storage/volume"
	"github.com/wqvoon/cbox/pkg/utils"
)

// 在容器停止时会将 containerInfo.Pid 设置为此值
const STOPPED_PID = 0

func GetContainerInfo(containerID string) *ContainerInfo {
	infoPath := rootdir.GetContainerInfoPath(containerID)
	infoObj := &ContainerInfo{}

	utils.GetObjFromJsonFile(infoPath, infoObj)
	infoObj.filePath = infoPath
	infoObj.containerID = containerID

	return infoObj
}

// 表示 Container 运行时的一些信息
type ContainerInfo struct {
	// info 文件相对于 rootdir 的路径
	filePath string
	// 该 Info 对象对应的 containerID
	containerID string

	// runtime 进程的 pid，被 runtime.Run 写入
	Pid int `json:"pid"`

	// 采用了哪个 StorageDriver，在 Container 创建时确定，不可更改
	StorageDriver string `json:"storage_driver"`

	// runtime 传递过来的宿主机 dns 文件路径
	DNSFilePath string `json:"dns_file"`

	Volumes []*volume.Volume `json:"volumes"`
}

// 判断 Container 是否在运行中
func (ci *ContainerInfo) IsRunning() bool {
	return ci.Pid != STOPPED_PID
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

func (ci *ContainerInfo) GetStorageDriver() driver.Interface {
	if len(ci.StorageDriver) == 0 {
		log.Errorln("invalid StorageDriver")
	}

	return driver.GetDriverByName(ci.StorageDriver)
}

func (ci *ContainerInfo) SavePid(pid int) {
	ci.Pid = pid
	ci.save()
}

// 标记容器已退出
func (ci *ContainerInfo) MarkStop() {
	ci.SavePid(STOPPED_PID)
}

// 这个方法仅应该被 container.CreateContainer 使用
func (ci *ContainerInfo) SaveStorageDriver(driverName string) {
	ci.StorageDriver = driverName
	ci.save()
}

func (ci *ContainerInfo) SaveDNSFilePath(filePath string) {
	ci.DNSFilePath = filePath
	ci.save()
}

func (ci *ContainerInfo) SaveVolumes(vs []*volume.Volume) {
	mntPath := rootdir.GetContainerMountPath(ci.containerID)

	// 保证容器路径指向的是宿主机视角下的绝对路径
	for _, v := range vs {
		if !strings.HasPrefix(v.ContainerPath, mntPath) {
			v.ContainerPath = path.Join(mntPath, v.ContainerPath)
		}
	}

	ci.Volumes = vs
	ci.save()
}

// 保存整个对象到 info 文件中
func (ci *ContainerInfo) save() {
	utils.SaveObjToJsonFile(ci.filePath, ci)
}
