package container

import (
	"os"
	"os/exec"

	"golang.org/x/sys/unix"

	"github.com/wqvoon/cbox/pkg/log"
	"github.com/wqvoon/cbox/pkg/rootdir"
	"github.com/wqvoon/cbox/pkg/storage/driver"
	"github.com/wqvoon/cbox/pkg/utils"
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

	cmd := &exec.Cmd{
		Path: name,
		Args: args,
		Dir:  "/",

		// TODO: 这部分的赋值应该可选
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,

		// TODO：这个 Env 放在这里是有问题的，因为本次的 Cmd 在被执行时的 Env 还来自于宿主机
		// 有可能会报路径不存在的错误
		Env: c.Env,
		SysProcAttr: &unix.SysProcAttr{
			Chroot: containerMntPoint,

			Cloneflags: unix.CLONE_NEWPID |
				unix.CLONE_NEWNS |
				unix.CLONE_NEWUTS |
				unix.CLONE_NEWIPC,
		},
	}

	if err := cmd.Run(); err != nil {
		log.Errorln("faild to start container, err:", err)
	}
}

func (c *Container) Stop() {
	log.TODO()
}

func (c *Container) Delete() {
	log.TODO()
}
