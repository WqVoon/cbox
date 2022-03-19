package container

import (
	"github.com/wqvoon/cbox/pkg/log"
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

func (c ContainerIdx) DeleteByName(name string) {
	if !c.Has(name) {
		log.Errorln("no such container in idx:", name)
	}

	delete(c, name)
	c.Save()
}

func (c ContainerIdx) Has(containerName string) bool {
	_, has := c[containerName]

	return has
}

func (c ContainerIdx) GetByName(name string) (string, *ContainerEntity) {
	if entity, isIn := c[name]; isIn {
		return name, entity
	}

	log.Errorln("no such container in containerIdx:", name)
	return "", nil
}

// TODO：可以像 Docker 一样做前缀匹配
func (c ContainerIdx) GetByID(id string) (string, *ContainerEntity) {
	for name, entity := range c {
		if entity.ContainerID == id {
			return name, entity
		}
	}

	log.Errorln("no such container in containerIdx:", id)
	return "", nil
}

// 遍历所有的记录，如果 fn 返回 false 那么提前结束遍历
func (c ContainerIdx) Range(fn func(string, *ContainerEntity) bool) {
	for name, entity := range c {
		if !fn(name, entity) {
			break
		}
	}
}
