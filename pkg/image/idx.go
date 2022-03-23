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

// 尝试更新，如果确实更新了那么返回 true，如果现有的值等同于 imageHash 参数，那么返回 false
func (idx ImageIdx) Update(nameTag *utils.NameTag, newHash string) bool {
	if entity, isIn := idx[nameTag.Name]; isIn {
		oldHash := entity[nameTag.Tag]
		if oldHash == newHash {
			return false
		}

		entity[nameTag.Tag] = newHash
	} else {
		idx[nameTag.Name] = ImageEntity{nameTag.Tag: newHash}
	}

	utils.SaveObjToJsonFile(rootdir.GetImageIdxPath(), idx)
	return true
}

func (idx ImageIdx) Range(fn func(repo, tag, hash string) bool) {
	for repo, entity := range idx {
		for tag, hash := range entity {
			if !fn(repo, tag, hash) {
				return
			}
		}
	}
}
