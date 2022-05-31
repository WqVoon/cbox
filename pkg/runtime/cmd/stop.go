package cmd

import (
	"path"

	"github.com/wqvoon/cbox/pkg/log"
	"github.com/wqvoon/cbox/pkg/rootdir"
	runtimeInfo "github.com/wqvoon/cbox/pkg/runtime/info"
	runtimeUtils "github.com/wqvoon/cbox/pkg/runtime/utils"
	"github.com/wqvoon/cbox/pkg/utils"
	"golang.org/x/sys/unix"
)

func Stop(info *runtimeInfo.ContainerInfo) {
	if err := info.GetProcess().Kill(); err != nil {
		log.Errorln("failed to kill runtime process ,err:", err)
	}

	// UnMount 必须在 Kill 之后，否则会报 device busy（至少对于 Overlay2 来说）
	utils.WaitForStop(info.Pid)
	info.MarkStop()

	mntPath := rootdir.GetContainerMountPath(info.ContainerID)

	procPath := path.Join(mntPath, "proc")
	if err := unix.Unmount(procPath, 0); err != nil {
		log.Errorln("faild to unmount proc, err:", err)
	}

	for _, v := range info.Volumes {
		v.Unmount()
	}

	runtimeUtils.DeleteCGroupForContainer(info.ContainerID)
}
