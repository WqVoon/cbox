package container

import (
	"crypto/rand"
	"fmt"
	"os"
	"os/exec"

	"github.com/wqvoon/cbox/pkg/rootdir"
	"github.com/wqvoon/cbox/pkg/utils"
	"golang.org/x/sys/unix"
)

func newContainerID() string {
	randBytes := make([]byte, 6)
	rand.Read(randBytes)
	return fmt.Sprintf("%02x%02x%02x%02x%02x%02x",
		randBytes[0], randBytes[1], randBytes[2],
		randBytes[3], randBytes[4], randBytes[5])
}

func createContainerLayout(containerID string) {
	containerMntPath := rootdir.GetContainerMountPath(containerID)
	utils.CreateDirWithExclusive(containerMntPath)

	infoPath := rootdir.GetContainerInfoPath(containerID)
	utils.WriteFileIfNotExist(infoPath, []byte("{}"))
}

func StartContainersByName(names ...string) {
	for _, name := range names {
		GetContainerByName(name).Start()
	}
}

func StopContainersByName(names ...string) {
	for _, name := range names {
		GetContainerByName(name).Stop()
	}
}

func DeleteContainersByName(names ...string) {
	for _, name := range names {
		GetContainerByName(name).Delete()
	}
}

// 为 Container 创建一个 cmd 对象，优先使用 input 作为命令，否则使用容器的 entrypoint
func getCmdForContainer(c *Container, input ...string) *exec.Cmd {
	var name string
	var args []string
	if len(input) > 0 {
		name, args = utils.ParseCmd(input...)
	} else {
		name, args = utils.ParseCmd(c.Entrypoint...)
	}

	return &exec.Cmd{
		Path:   name,
		Args:   append([]string{name}, args...),
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Dir:    "/",
		Env:    c.Env,
		SysProcAttr: &unix.SysProcAttr{
			Chroot: rootdir.GetContainerMountPath(c.ID),
		},
	}
}
