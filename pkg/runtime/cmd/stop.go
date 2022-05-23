package cmd

import (
	"path"
	"time"

	"github.com/wqvoon/cbox/pkg/log"
	"github.com/wqvoon/cbox/pkg/rootdir"
	runtimeInfo "github.com/wqvoon/cbox/pkg/runtime/info"
	runtimeUtils "github.com/wqvoon/cbox/pkg/runtime/utils"
	"golang.org/x/sys/unix"
)

func Stop(info *runtimeInfo.ContainerInfo) {
	mntPath := rootdir.GetContainerMountPath(info.ContainerID)

	procPath := path.Join(mntPath, "proc")
	if err := unix.Unmount(procPath, 0); err != nil {
		log.Errorln("faild to unmount proc, err:", err)
	}

	for _, v := range info.Volumes {
		v.Unmount()
	}

	if err := info.GetProcess().Kill(); err != nil {
		log.Errorln("failed to kill runtime process ,err:", err)
	}
	info.MarkStop()

	// UnMount 必须在 Kill 之后，否则会报 device busy（至少对于 Overlay2 来说）
	// TODO: 这里简单等待100ms，后面整个更稳妥的办法确保进程退出后再执行 UnMount
	time.Sleep(100 * time.Millisecond)

	runtimeUtils.DeleteCGroupForContainer(info.ContainerID)
}
