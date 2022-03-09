package image

import (
	"log"
	"path"
	"strings"

	"github.com/wqvoon/cbox/pkg/rootdir"
	"github.com/wqvoon/cbox/pkg/utils"
)

// image name -> image hash
type ImageEntity map[string]string
type ImageIdx map[string]ImageEntity

type Manifest []*struct {
	RootPath   string
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
		oneManifest.RootPath = rootdir.GetImagePath(hash)

		absImageConfigPath := rootdir.GetImageConfigPath(hash, oneManifest.ConfigPath)

		oneManifest.Config = &ImageConfig{}
		utils.GetObjFromJsonFile(absImageConfigPath, oneManifest.Config)

		for idx, layerPath := range oneManifest.Layers {
			const tail = "/layer.tar"
			const tailLen = len(tail)

			if !strings.HasSuffix(layerPath, tail) {
				log.Fatalln("error layer format:", layerPath)
			}

			oneManifest.Layers[idx] = layerPath[:len(layerPath)-tailLen]
		}
	}

	return ret
}

func (manifest Manifest) GetLayerPaths() []string {
	if len(manifest) != 1 {
		log.Fatalln("unsupported length of manifest")
	}

	oneManifest := manifest[0]
	ret := make([]string, 0, len(oneManifest.Layers))

	for _, layerPath := range oneManifest.Layers {
		ret = append(ret, path.Join(oneManifest.RootPath, layerPath, "fs")+"/")
	}

	return ret
}
