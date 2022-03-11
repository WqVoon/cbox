package image

import (
	"github.com/wqvoon/cbox/pkg/log"
	"github.com/wqvoon/cbox/pkg/rootdir"
	"github.com/wqvoon/cbox/pkg/utils"
)

// image tag -> image hash
type ImageEntity map[string]string

// image name -> ImageEntity
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
		log.Errorln("no such image in imageIdx:", nameTag)
	}

	hash, isIn := entity[nameTag.Tag]
	if !isIn {
		log.Errorln("no such image in imageIdx:", nameTag)
	}

	return hash
}

func (idx ImageIdx) GetImageNameTag(wantHash string) *utils.NameTag {
	// TODO: 当前的实现很低效
	for name, entity := range idx {
		for tag, hash := range entity {
			if hash != wantHash {
				continue
			}

			return &utils.NameTag{
				Name: name,
				Tag:  tag,
			}
		}
	}

	log.Errorln("no such image in imageIdx:", wantHash)
	return nil
}
