package image

import (
	"log"
	"path"

	"github.com/wqvoon/cbox/pkg/rootdir"
	"github.com/wqvoon/cbox/pkg/utils"
)

// image name -> image hash
type ImageEntity map[string]string
type ImageIdx map[string]ImageEntity

type Manifest []struct {
	Config   string
	RepoTags []string
	Layers   []string
}

func GetIdx() ImageIdx {
	var ret ImageIdx

	idxFilePath := path.Join(rootdir.GetPath(), "images", "images.json")
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

	manifestFilePath := path.Join(rootdir.GetPath(), "images", hash, "manifest.json")
	utils.GetObjFromJsonFile(manifestFilePath, &ret)

	return ret
}
