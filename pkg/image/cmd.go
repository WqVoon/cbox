package image

import (
	"path"
	"path/filepath"

	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/wqvoon/cbox/pkg/log"
	"github.com/wqvoon/cbox/pkg/rootdir"
	"github.com/wqvoon/cbox/pkg/utils"
)

func Pull(nameTag *utils.NameTag) {
	img, err := crane.Pull(nameTag.String())
	if err != nil {
		log.Errorln("failed to pull image, err:", err)
	}

	var imgHash string
	if manifest, err := img.Manifest(); err != nil {
		log.Errorln("failed to get manifest, err:", err)
	} else {
		imgHash = manifest.Config.Digest.Hex[:12]
	}

	if changed := GetImageIdx().Update(nameTag, imgHash); !changed {
		log.Println("image", nameTag, "has exists")
		return
	}

	tarballPath := rootdir.GetImageTarballPath(imgHash)
	if utils.PathIsExist(tarballPath) {
		// 比如 images/idx.json 被修改过，或者 centos:latest 和 centos:xxx 有可能指向的是同一个镜像
		// 这个时候如果对应的 tarball 已经存在，那么不需要进行 Save
		log.Println("use cached tarball for", nameTag)
	} else {
		utils.WriteFileIfNotExist(tarballPath, nil)

		if err := crane.SaveLegacy(img, nameTag.String(), tarballPath); err != nil {
			log.Errorln("failed to save image, err:", err)
		}
	}

	imageLayoutPath := rootdir.GetImageLayoutPath(imgHash)
	if utils.PathIsExist(imageLayoutPath) {
		// 原因同上，直接跳过解包环节
		log.Println("use cached image layout for", nameTag)
		return
	}

	utils.Untar(tarballPath, imageLayoutPath)

	manifest := GetManifestByHash(imgHash)
	for _, layer := range manifest.Layers {
		layerTarPath := path.Join(imageLayoutPath, layer)

		dirName := filepath.Dir(layer)
		layerFsPath := rootdir.GetImageFsPath(imgHash, dirName)

		utils.Untar(layerTarPath, layerFsPath)
	}

	log.Println("downloaded image for", nameTag)
}
