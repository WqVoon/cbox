package volume

import (
	"github.com/wqvoon/cbox/pkg/log"
	"github.com/wqvoon/cbox/pkg/utils"
	"golang.org/x/sys/unix"
)

type Volume struct {
	// ReadOnly      bool
	HostPath      string `json:"host_path"`
	ContainerPath string `json:"container_path"`
}

func (v *Volume) Mount() {
	utils.CreateDirIfNotExist(v.ContainerPath)

	if err := unix.Mount(v.HostPath, v.ContainerPath, "bind", unix.MS_BIND, ""); err != nil {
		log.Errorf("failed to bind mount %q to %q, err: %v\n",
			v.HostPath, v.ContainerPath, err)
	}
}

func (v *Volume) Unmount() {
	if err := unix.Unmount(v.ContainerPath, 0); err != nil {
		log.Errorf("failed to unmount %q, err: %v\n", v.ContainerPath, err)
	}
}
