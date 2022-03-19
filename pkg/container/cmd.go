package container

import (
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/wqvoon/cbox/pkg/log"
	"github.com/wqvoon/cbox/pkg/rootdir"
	runtimeCmd "github.com/wqvoon/cbox/pkg/runtime/cmd"
	runtimeUtils "github.com/wqvoon/cbox/pkg/runtime/utils"
	"github.com/wqvoon/cbox/pkg/storage/driver"
	"github.com/wqvoon/cbox/pkg/utils"
	"golang.org/x/sys/unix"
)

func (c *Container) Start() {
	if runtimeUtils.GetContainerInfo(c.ID).IsRunning() {
		log.Errorln("can not start a running container")
		return
	}

	if os.Geteuid() != 0 {
		log.Errorln("only root user can start a container")
	}

	containerMntPoint := rootdir.GetContainerMountPath(c.ID)
	driver.D.Mount(containerMntPoint, c.Image.Layers...)

	runtimeCmd.Run(c.ID)

	log.Println("container started")
}

func (c *Container) Exec(input ...string) {
	if !runtimeUtils.GetContainerInfo(c.ID).IsRunning() {
		log.Errorln("can not exec a not-running container")
		return
	}

	if os.Geteuid() != 0 {
		log.Errorln("only root user can exec a container")
	}

	var name string
	var args []string
	if len(input) > 0 {
		name, args = utils.ParseCmd(input...)
	} else {
		name, args = utils.ParseCmd(c.Entrypoint...)
	}

	enterNamespace(c.ID)

	// 需要保证在 ExtractCmdFromOSArgs 前进行 Env 的处理，这样得到的 cmd 才是正确的
	os.Clearenv()
	for _, oneEnv := range c.Env {
		envPair := strings.Split(oneEnv, "=")
		key, val := envPair[0], envPair[1]
		os.Setenv(key, val)
	}

	// 这里需要保证在 ExtractCmdFromOSArgs 前进行 chroot，这样得到的 cmd 才是正确的
	unix.Chroot(rootdir.GetContainerMountPath(c.ID))

	cmd := exec.Command(name, args...)
	{
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Dir = "/"
		cmd.Env = c.Env
	}

	if err := cmd.Run(); err != nil {
		// 这里不进行 Errorln，用于避免 Ctrl+C 后 Ctrl+D 引起的常见错误
		log.Println("an error may have occurred while running container:", err)
	}
}

func (c *Container) Stop() {
	if !runtimeUtils.GetContainerInfo(c.ID).IsRunning() {
		log.Errorln("can not stop a not-running container")
		return
	}

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

	info := runtimeUtils.GetContainerInfo(c.ID)
	if err := info.GetProcess().Kill(); err != nil {
		log.Errorln("failed to kill runtime process ,err:", err)
	}
	info.SavePid(runtimeUtils.STOPPED_PID)

	log.Println("container stopped")
}

func (c *Container) Delete() {
	if runtimeUtils.GetContainerInfo(c.ID).IsRunning() {
		log.Errorln("can not delete a running container")
		return
	}

	GetContainerIdx().DeleteByName(c.Name)
	if err := os.RemoveAll(c.rootPath); err != nil {
		log.Errorf("faild to remove container %q, err: %v\n", c.Name, err)
	}
	// TODO: 后面要处理更多的运行时副作用

	log.Println("container deleted")
}
