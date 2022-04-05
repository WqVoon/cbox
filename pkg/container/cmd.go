package container

import (
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/wqvoon/cbox/pkg/log"
	"github.com/wqvoon/cbox/pkg/rootdir"
	runtimeCmd "github.com/wqvoon/cbox/pkg/runtime/cmd"
	runtimeInfo "github.com/wqvoon/cbox/pkg/runtime/info"
	"github.com/wqvoon/cbox/pkg/storage/driver"
	"github.com/wqvoon/cbox/pkg/utils"
	"golang.org/x/sys/unix"
)

func (c *Container) Start() {
	info := runtimeInfo.GetContainerInfo(c.ID)
	{
		if info.IsRunning() {
			log.Errorln("can not start a running container")
			return
		}
		if info.GetStorageDriver() != driver.D {
			log.Errorf("can not use %s to start %s", driver.D, c.Name)
		}
	}

	if os.Geteuid() != 0 {
		log.Errorln("only root user can start a container")
	}

	containerMntPoint := rootdir.GetContainerMountPath(c.ID)
	driver.D.Mount(containerMntPoint, c.Image.Layers...)

	runtimeCmd.Run(c.ID)

	log.Printf("container %s started\n", c.Name)
}

func (c *Container) Exec(input ...string) {
	if !runtimeInfo.GetContainerInfo(c.ID).IsRunning() {
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
	info := runtimeInfo.GetContainerInfo(c.ID)
	{
		if !info.IsRunning() {
			log.Errorln("can not stop a not-running container")
			return
		}
		if info.GetStorageDriver() != driver.D {
			log.Errorf("can not use %s to stop %s", driver.D, c.Name)
		}
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
	driver.D.UnMount(mntPath)

	log.Printf("container %s stopped\n", c.Name)
}

func (c *Container) Delete() {
	if runtimeInfo.GetContainerInfo(c.ID).IsRunning() {
		log.Errorln("can not delete a running container")
		return
	}

	runtimeInfo.GetImageInfo(c.Image.Hash).MarkReleasedBy(c.ID)

	GetContainerIdx().DeleteByName(c.Name)
	if err := os.RemoveAll(c.rootPath); err != nil {
		log.Errorf("faild to remove container %q, err: %v\n", c.Name, err)
	}
	// TODO: 后面要处理更多的运行时副作用

	log.Printf("container %s deleted\n", c.Name)
}
