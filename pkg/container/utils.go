package container

import (
	"crypto/rand"
	"fmt"
	"os"
	"path"

	"github.com/wqvoon/cbox/pkg/log"
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

	containerNSPath := rootdir.GetContainerNSPath(containerID)
	utils.CreateDirWithExclusive(containerNSPath)

	namespaces := []string{"ipc", "uts", "pid"}
	pathPrefix := rootdir.GetContainerNSPath(containerID)
	for _, ns := range namespaces {
		fullPath := path.Join(pathPrefix, ns)
		utils.WriteFileIfNotExist(fullPath, nil)
	}

	infoPath := rootdir.GetContainerInfoPath(containerID)
	utils.WriteFileIfNotExist(infoPath, []byte("{}"))
}

// 让当前进程及其子进程进入 containerID 对应的 Container 的命名空间中
func enterNamespace(containerID string) {
	nsPathPrefix := rootdir.GetContainerNSPath(containerID)
	// nsPathPrefix := "/proc/114195/ns"
	ipcFd, ipcErr := os.Open(nsPathPrefix + "/ipc")
	pidFd, pidErr := os.Open(nsPathPrefix + "/pid")
	utsFd, utsErr := os.Open(nsPathPrefix + "/uts")
	if ipcErr != nil || pidErr != nil || utsErr != nil {
		log.Errorln("failed to open namespace")
	}

	if err := unix.Setns(int(ipcFd.Fd()), unix.CLONE_NEWIPC); err != nil {
		log.Errorln("faild to enter ipc ns, err:", err)
	}
	if err := unix.Setns(int(pidFd.Fd()), unix.CLONE_NEWPID); err != nil {
		log.Errorln("faild to enter pid ns, err:", err)
	}
	if err := unix.Setns(int(utsFd.Fd()), unix.CLONE_NEWUTS); err != nil {
		log.Errorln("faild to enter uts ns, err:", err)
	}

	ipcFd.Close()
	pidFd.Close()
	utsFd.Close()
}
