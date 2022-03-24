package container

import (
	"fmt"

	"github.com/wqvoon/cbox/pkg/image"
	"github.com/wqvoon/cbox/pkg/log"
	"github.com/wqvoon/cbox/pkg/rootdir"
	runtimeInfo "github.com/wqvoon/cbox/pkg/runtime/info"
	"github.com/wqvoon/cbox/pkg/storage/driver"
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
	// TODO: 可能需要补充 namespace、pid 等运行时内容
}

func CreateContainer(img *image.Image, name string) *Container {
	// TODO：
	//  后面可以试着抽象出接口，然后把这个 CreateContainer 方法直接给 Image 对象
	//  当前会有 image 和 container 循环引用的问题
	containerID := newContainerID()

	idx := GetContainerIdx()
	if idx.Has(name) {
		log.Errorln("container name has exists, try another plz")
	}

	idx[name] = &ContainerEntity{
		ContainerID: containerID,
		ImageHash:   img.Hash,
	}
	idx.Save()

	createContainerLayout(containerID)

	runtimeInfo.GetContainerInfo(containerID).SaveStorageDriver(driver.D.String())
	runtimeInfo.GetImageInfo(img.Hash).MarkUsedBy(containerID)

	log.Printf("create container %s(%s)\n", name, containerID)

	return getContainerHelper(img, containerID, name)
}

// TODO: 新增一个 GetContainer 方法，依次按 name、id 搜索
func GetContainerByName(name string) *Container {
	name, entity := GetContainerIdx().GetByName(name)
	img := image.GetImageFromLocalByHash(entity.ImageHash)

	return getContainerHelper(img, entity.ContainerID, name)
}

func GetContainerByID(id string) *Container {
	name, entity := GetContainerIdx().GetByID(id)
	img := image.GetImageFromLocalByHash(entity.ImageHash)

	return getContainerHelper(img, id, name)
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

// 仅 container 内部使用，根据已知信息帮助创建 Container 对象
// 原因同 image.getImageHelper 方法
func getContainerHelper(img *image.Image, containerID, containerName string) *Container {
	return &Container{
		rootPath: rootdir.GetContainerLayoutPath(containerID),

		ID:         containerID,
		Name:       containerName,
		Env:        img.Config.Config.Env,
		Entrypoint: img.Config.Config.Cmd,
		Image:      img,
	}
}

// 展示所有的 container 信息，类似于 `docker ps`，每个字段占 16 个字符长度
// TODO：后面可以加一些 filter，并且以表格的形式输出
func ListAllContainer() {
	tw := utils.NewTableWriter([]string{
		"container name", "container id", "image", "command", "running", "driver",
	}, 16)

	tw.PrintlnHeader()

	GetContainerIdx().Range(func(name string, entity *ContainerEntity) bool {
		c := GetContainerByName(name)
		info := runtimeInfo.GetContainerInfo(c.ID)

		containerName := c.Name
		containerID := c.ID
		imageName := image.GetImageIdx().GetImageNameTag(entity.ImageHash).String()
		command := fmt.Sprint(c.Entrypoint)
		status := fmt.Sprint(info.IsRunning())
		driver := info.StorageDriver

		tw.PrintlnData(containerName, containerID, imageName, command, status, driver)
		return true
	})
}
