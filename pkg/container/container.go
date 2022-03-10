package container

import (
	"crypto/rand"
	"fmt"

	"github.com/wqvoon/cbox/pkg/rootdir"
	"github.com/wqvoon/cbox/pkg/utils"
)

func NewContainerID() string {
	randBytes := make([]byte, 6)
	rand.Read(randBytes)
	return fmt.Sprintf("%02x%02x%02x%02x%02x%02x",
		randBytes[0], randBytes[1], randBytes[2],
		randBytes[3], randBytes[4], randBytes[5])
}

func CreateContainerRootDir(containerID string) {
	containerMntPath := rootdir.GetContainerMountPath(containerID)

	utils.CreateDirWithExclusive(containerMntPath)
}

// func MountFSByRawCopy(manifest image.Manifest, containerID string) {
// 	containerMntPath := rootdir.GetContainerMountPath(containerID)

// 	for _, layerPath := range manifest.GetLayerPaths() {
// 		utils.CopyDir(layerPath, containerMntPath)
// 	}
// }
