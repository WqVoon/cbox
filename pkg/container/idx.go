package container

import (
	"github.com/wqvoon/cbox/pkg/rootdir"
	"github.com/wqvoon/cbox/pkg/utils"
)

type ContainerEntity struct {
	ContainerID string `json:"container_id"`
	ImageHash   string `json:"image_hash"`
}
type ContainerIdx map[string]*ContainerEntity

func GetContainerIdx() ContainerIdx {
	var ret ContainerIdx

	idxFilePath := rootdir.GetContainerIdxPath()
	utils.GetObjFromJsonFile(idxFilePath, &ret)

	return ret
}

func (c ContainerIdx) Save() {
	idxFilePath := rootdir.GetContainerIdxPath()

	utils.SaveObjToJsonFile(idxFilePath, c)
}

func (c ContainerIdx) Has(containerName string) bool {
	_, has := c[containerName]

	return has
}
