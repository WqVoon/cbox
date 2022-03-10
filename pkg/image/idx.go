package image

import (
	"log"

	"github.com/wqvoon/cbox/pkg/rootdir"
	"github.com/wqvoon/cbox/pkg/utils"
)

// image name -> image hash
type ImageEntity map[string]string
type ImageIdx map[string]ImageEntity

func GetImageIdx() ImageIdx {
	var ret ImageIdx

	idxFilePath := rootdir.GetImageIdxPath()
	utils.GetObjFromJsonFile(idxFilePath, &ret)

	return ret
}

func (idx ImageIdx) GetImageHash(nameTag *utils.NameTag) string {
	entity, isIn := idx[nameTag.Name]
	if !isIn {
		log.Fatalln("no such image in imageIdx:", nameTag)
	}

	hash, isIn := entity[nameTag.Tag]
	if !isIn {
		log.Fatalln("no such image in imageIdx:", nameTag)
	}

	return hash
}
