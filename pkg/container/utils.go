package container

import (
	"crypto/rand"
	"fmt"
	"path"

	"github.com/wqvoon/cbox/pkg/rootdir"
	"github.com/wqvoon/cbox/pkg/utils"
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
}
