package container

import (
	"fmt"

	"github.com/wqvoon/cbox/pkg/image"
	"github.com/wqvoon/cbox/pkg/rootdir"
	"github.com/wqvoon/cbox/pkg/utils"
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
}

func CreateContainer(img *image.Image, name string) *Container {
	containerID := newContainerID()

	createContainerLayout(containerID)
	/*
		TODO: 写入一个 json 文件用于记录
		格式形如: {
			<containerName>: {
				"container_id": <id>,
				"image_id": <id>
			}
		}
		可以根据容器名找到容器和其镜像
	*/

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
	utils.TODO()
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
