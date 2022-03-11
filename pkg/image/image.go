package image

import (
	"fmt"

	"github.com/wqvoon/cbox/pkg/rootdir"
	"github.com/wqvoon/cbox/pkg/utils"
)

type Image struct {
	// 仅内部使用，用于快速定位当前镜像的 layout 位置
	rootPath string
	// 镜像 layout 的文件夹名
	Hash string
	// 镜像的 nameTag
	NameTag *utils.NameTag
	// manifest.json 中 Config 对应的文件解码后的内容
	Config *ImageConfig
	// manifest.json 解码后的内容
	Manifest *ManifestType
	// Layer fs 对应的文件夹相对于 rootdir 的路径
	Layers []string
}

// TODO: 后面可以做成优先搜索 nameTag，找不到再搜 hash，再找不到后退出
var GetImage = GetImageFromLocalByNameTag

func GetImageFromLocalByNameTag(nameTag *utils.NameTag) *Image {
	hash := GetImageIdx().GetImageHash(nameTag)

	return getImageHelper(nameTag, hash)
}

func GetImageFromLocalByHash(hash string) *Image {
	nameTag := GetImageIdx().GetImageNameTag(hash)

	return getImageHelper(nameTag, hash)
}

func (img *Image) String() string {
	return fmt.Sprintf(`
Image(%s):
	Hash: %s
	ConfigPath: %s
	ManifestPath: %s
	Layers: %v
`,
		img.NameTag, img.Hash, img.Config.rootPath, img.Manifest.rootPath, img.Layers,
	)
}

// 仅 image 内部使用，根据已知信息帮助创建 Image 对象
// 之所以抽离出这个方法，是因为后面可能会以 functional options 模式传参
func getImageHelper(nameTag *utils.NameTag, hash string) *Image {
	manifest := GetManifestByHash(hash)

	return &Image{
		rootPath: rootdir.GetImageLayoutPath(hash),

		Hash:     hash,
		NameTag:  nameTag,
		Manifest: manifest,
		Config:   manifest.GetConfigByHash(hash),
		Layers:   manifest.GetLayerFSPaths(),
	}
}
