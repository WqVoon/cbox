package volume

import (
	"os"

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
	stat, err := os.Stat(v.HostPath)
	if err != nil {
		log.Errorf("failed to stat %q, err: %v\n", v.HostPath, err)
	}

	if stat.IsDir() {
		utils.CreateDirIfNotExist(v.ContainerPath)
	} else {
		utils.WriteFileIfNotExist(v.ContainerPath, nil)
	}

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
