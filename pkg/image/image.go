package image

import (
	"fmt"

	"github.com/wqvoon/cbox/pkg/log"
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

func GetImage(nameTag *utils.NameTag) *Image {
	log.TODO()
	return nil
}

func GetImageFromLocal(nameTag *utils.NameTag) *Image {
	hash := GetImageIdx().GetImageHash(nameTag)

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
