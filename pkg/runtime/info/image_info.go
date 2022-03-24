package info

import (
	"github.com/wqvoon/cbox/pkg/rootdir"
	"github.com/wqvoon/cbox/pkg/utils"
)

func GetImageInfo(imageHash string) *ImageInfo {
	infoPath := rootdir.GetImageInfoPath(imageHash)
	infoObj := &ImageInfo{}

	utils.WriteFileIfNotExist(infoPath, []byte("{}"))

	utils.GetObjFromJsonFile(infoPath, infoObj)
	infoObj.filePath = infoPath

	return infoObj
}

// 表示 Image 运行时的一些信息
type ImageInfo struct {
	// info 文件相对于 rootdir 的路径
	filePath string

	// 表示当前 Image 被哪些 Container 所使用，此值非空时对应的 Image 不可删除
	UsedBy []string `json:"used_by"`
}

// 当前 Image 是否可以删除
func (info *ImageInfo) CanBeDeleted() bool {
	return !info.BeUsed()
}

// 当前 Image 是否在使用中
func (info *ImageInfo) BeUsed() bool {
	return len(info.UsedBy) != 0
}

// 标记当前容器被某个容器使用了
func (info *ImageInfo) MarkUsedBy(containerID string) {
	for _, item := range info.UsedBy {
		if item == containerID {
			return
		}
	}

	info.UsedBy = append(info.UsedBy, containerID)
	info.save()
}

// 标记被哪个容器释放了
func (info *ImageInfo) MarkReleasedBy(containerID string) {
	usedBy := info.UsedBy

	for i := 0; i < len(usedBy); i++ {
		if usedBy[i] == containerID {
			usedBy = append(usedBy[:i], usedBy[i+1:]...)
			i--
		}
	}

	info.UsedBy = usedBy
	info.save()
}

func (info *ImageInfo) save() {
	utils.SaveObjToJsonFile(info.filePath, info)
}
