package image

import (
	"log"

	"github.com/wqvoon/cbox/pkg/rootdir"
	"github.com/wqvoon/cbox/pkg/utils"
)

type ManifestType struct {
	rootPath string

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
		log.Fatalln("unsupported length of manifest")
	}

	manifest := lst[0]
	manifest.rootPath = rootdir.GetManifestPath(hash)

	return manifest
}
