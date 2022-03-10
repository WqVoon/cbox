package container

import (
	"fmt"

	"github.com/wqvoon/cbox/pkg/image"
	"github.com/wqvoon/cbox/pkg/log"
	"github.com/wqvoon/cbox/pkg/rootdir"
)

// func MountFSByRawCopy(manifest image.Manifest, containerID string) {
// 	containerMntPath := rootdir.GetContainerMountPath(containerID)

// 	for _, layerPath := range manifest.GetLayerPaths() {
// 		utils.CopyDir(layerPath, containerMntPath)
// 	}
// }

type Container struct {
	rootPath string

	ID         string
	Name       string
	Env        []string
	Entrypoint []string
	Image      *image.Image
	// TODO: 可能需要补充 namespace、pid 等运行时内容
}

func CreateContainer(img *image.Image, name string) *Container {
	containerID := newContainerID()

	createContainerLayout(containerID)

	idx := GetContainerIdx()
	if idx.Has(name) {
		log.Errorln("container name has exists, try another plz")
	}

	idx[name] = &ContainerEntity{
		ContainerID: containerID,
		ImageHash:   img.Hash,
	}
	idx.Save()

	return &Container{
		rootPath: rootdir.GetContainerLayoutPath(containerID),

		ID:         containerID,
		Name:       name,
		Env:        img.Config.Config.Env,
		Entrypoint: img.Config.Config.Cmd,
		Image:      img,
	}
}

func (c *Container) Start() {
	log.TODO()
}

func (c *Container) String() string {
	return fmt.Sprintf(`
Container(%s):
	ID: %s
	Env: %v
	Entrypoint: %v
`,
		c.Name, c.ID, c.Env, c.Entrypoint,
	)
}
