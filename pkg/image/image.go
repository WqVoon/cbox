package image

import (
	"log"

	"github.com/wqvoon/cbox/pkg/rootdir"
	"github.com/wqvoon/cbox/pkg/utils"
)

// image name -> image hash
type ImageEntity map[string]string
type ImageIdx map[string]ImageEntity

type Manifest []*struct {
	ConfigPath string `json:"config"`
	RepoTags   []string
	Layers     []string
	Config     *ImageConfig `json:"-"`
}

type ImageConfigDetail struct {
	Env []string `json:"Env"`
	Cmd []string `json:"Cmd"`
}
type ImageConfig struct {
	Config ImageConfigDetail `json:"config"`
}

func GetIdx() ImageIdx {
	var ret ImageIdx

	idxFilePath := rootdir.GetImageIdxPath()
	utils.GetObjFromJsonFile(idxFilePath, &ret)

	return ret
}

func (idx ImageIdx) GetManifest(nameTag *utils.NameTag) Manifest {
	var ret Manifest

	entity, isIn := idx[nameTag.Name]
	if !isIn {
		log.Fatalln("no such image in imageIdx:", nameTag)
	}

	hash, isIn := entity[nameTag.Tag]
	if !isIn {
		log.Fatalln("no such image in imageIdx:", nameTag)
	}

	manifestFilePath := rootdir.GetManifestPath(hash)
	utils.GetObjFromJsonFile(manifestFilePath, &ret)

	for _, oneManifest := range ret {
		absImageConfigPath := rootdir.GetImageConfigPath(hash, oneManifest.ConfigPath)

		oneManifest.Config = &ImageConfig{}
		utils.GetObjFromJsonFile(absImageConfigPath, oneManifest.Config)
	}

	return ret
}
