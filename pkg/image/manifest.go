package image

import (
	"github.com/wqvoon/cbox/pkg/log"
	"github.com/wqvoon/cbox/pkg/rootdir"
	"github.com/wqvoon/cbox/pkg/utils"
)

type ManifestType struct {
	rootPath string
	filePath string

	Config   string
	RepoTags []string
	Layers   []string
}

type ManifestList []*ManifestType

func GetManifestByNameTag(nameTag *utils.NameTag) *ManifestType {
	hash := GetImageIdx().GetImageHash(nameTag)

	return GetManifestByHash(hash)
}

func GetManifestByHash(hash string) *ManifestType {
	var lst ManifestList

	manifestFilePath := rootdir.GetManifestPath(hash)
	utils.GetObjFromJsonFile(manifestFilePath, &lst)

	if len(lst) != 1 {
		log.Errorln("unsupported length of manifest")
	}

	manifest := lst[0]
	manifest.rootPath = rootdir.GetImageLayoutPath(hash)
	manifest.filePath = rootdir.GetManifestPath(hash)

	return manifest
}
