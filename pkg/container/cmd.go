package container

import (
	"os"
	"path"

	"github.com/wqvoon/cbox/pkg/log"
	"github.com/wqvoon/cbox/pkg/rootdir"
	"github.com/wqvoon/cbox/pkg/runtime/cmd"
	"github.com/wqvoon/cbox/pkg/storage/driver"
	"github.com/wqvoon/cbox/pkg/utils"
	"golang.org/x/sys/unix"
)

func (c *Container) Start(input ...string) {
	if os.Geteuid() != 0 {
		log.Errorln("only root user can start a container")
	}

	containerMntPoint := rootdir.GetContainerMountPath(c.ID)
	driver.D.Mount(containerMntPoint, c.Image.Layers...)

	var name string
	var args []string
	if len(input) > 0 {
		name, args = utils.ParseCmd(input...)
	} else {
		name, args = utils.ParseCmd(c.Entrypoint...)
	}

	cmd.Run(c.ID, name, args)
}

func (c *Container) Stop() {
	mntPath := rootdir.GetContainerMountPath(c.ID)

	procPath := path.Join(mntPath, "proc")
	if err := unix.Unmount(procPath, 0); err != nil {
		log.Errorln("faild to unmount proc, err:", err)
	}

	namespaces := []string{"/pid", "/uts", "/ipc"}
	nsPath := rootdir.GetContainerNSPath(c.ID)
	for _, ns := range namespaces {
		dstPath := path.Join(nsPath, ns)
		if err := unix.Unmount(dstPath, 0); err != nil {
			log.Errorf("failed to unmount %q, err: %v\n", dstPath, err)
		}
	}
}

func (c *Container) Delete() {
	// TODO: 先检测是否执行过 Stop

	GetContainerIdx().DeleteByName(c.Name)
	if err := os.RemoveAll(c.rootPath); err != nil {
		log.Errorf("faild to remove container %q, err: %v\n", c.Name, err)
	}
	// TODO: 后面要处理更多的运行时副作用
}
