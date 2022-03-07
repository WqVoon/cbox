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

func (idx ImageIdx) GetManifest(nameTag string) Manifest {
	var ret Manifest

	splitedParam := strings.Split(nameTag, ":")
	if len(splitedParam) != 2 {
		log.Fatalln("error format of nameTag, should be `name:tag`")
	}
	name, tag := splitedParam[0], splitedParam[1]

	entity, isIn := idx[name]
	if !isIn {
		log.Fatalf("no such image in imageIdx: %s:%s\n", name, tag)
	}

	hash, isIn := entity[tag]
	if !isIn {
		log.Fatalf("no such image in imageIdx: %s:%s\n", name, tag)
	}

	manifestFilePath := path.Join(rootdir.GetPath(), "images", hash, "manifest.json")
	utils.GetObjFromJsonFile(manifestFilePath, &ret)

	return ret
}
