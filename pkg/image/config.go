package image

import (
	"github.com/wqvoon/cbox/pkg/rootdir"
	"github.com/wqvoon/cbox/pkg/utils"
)

type ImageConfigDetail struct {
	Env []string `json:"Env"`
	Cmd []string `json:"Cmd"`
}
type ImageConfig struct {
	rootPath string

	Config ImageConfigDetail `json:"config"`
}

func (manifest *ManifestType) GetConfigByNameTag(nameTag *utils.NameTag) *ImageConfig {
	hash := GetImageIdx().GetImageHash(nameTag)

	return manifest.GetConfigByHash(hash)
}

func (manifest *ManifestType) GetConfigByHash(hash string) *ImageConfig {
	imageConfigPath := rootdir.GetImageConfigPath(hash, manifest.Config)

	config := &ImageConfig{}
	config.rootPath = imageConfigPath
	utils.GetObjFromJsonFile(imageConfigPath, config)

	return config
}
