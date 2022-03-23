package image

import (
	"os"
	"path"
	"path/filepath"

	"github.com/google/go-containerregistry/pkg/crane"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/wqvoon/cbox/pkg/log"
	"github.com/wqvoon/cbox/pkg/rootdir"
	"github.com/wqvoon/cbox/pkg/utils"
)

func Pull(nameTag *utils.NameTag) {
	log.Println("start download image for", nameTag)

	var img v1.Image
	var imgHash string
	utils.NewTask("getting manifest", func() {
		var err error
		img, err = crane.Pull(nameTag.String())
		if err != nil {
			log.Errorln("failed to pull image, err:", err)
		}

		if manifest, err := img.Manifest(); err != nil {
			log.Errorln("failed to get manifest, err:", err)
		} else {
			imgHash = manifest.Config.Digest.Hex[:12]
		}
	}).Start()

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
		utils.NewTask("downloading tarball", func() {
			utils.WriteFileIfNotExist(tarballPath, nil)
			if err := crane.SaveLegacy(img, nameTag.String(), tarballPath); err != nil {
				log.Errorln("failed to save image, err:", err)
			}
		}).Start()
	}

	imageLayoutPath := rootdir.GetImageLayoutPath(imgHash)
	if utils.PathIsExist(imageLayoutPath) {
		// 原因同上，直接跳过解包环节
		log.Println("use cached image layout for", nameTag)
		return
	}

	utils.NewTask("untaring tarball", func() {
		utils.Untar(tarballPath, imageLayoutPath)

		manifest := GetManifestByHash(imgHash)
		for _, layer := range manifest.Layers {
			layerTarPath := path.Join(imageLayoutPath, layer)

			dirName := filepath.Dir(layer)
			layerFsPath := rootdir.GetImageFsPath(imgHash, dirName)

			utils.Untar(layerTarPath, layerFsPath)
		}
	}).Start()

	log.Println("success download image for", nameTag)
}

// TODO: 同样的，后面可以加一些 filter
func ListAllImage() {
	tw := utils.NewTableWriter([]string{
		"repository", "tag", "image id",
	}, 32)

	tw.PrintlnHeader()

	GetImageIdx().Range(func(repo, tag, hash string) bool {
		tw.PrintlnData(repo, tag, hash)
		return true
	})
}

func DeleteImagesNyName(names ...string) {
	imgIdx := GetImageIdx()

	for _, name := range names {
		// TODO: 检测是否使用中
		nameTag := utils.GetNameTag(name)
		imgHash := imgIdx.GetImageHash(nameTag)

		if err := os.RemoveAll(rootdir.GetImageLayoutPath(imgHash)); err != nil {
			log.Errorln("failed to remove image layout, err:", err)
		}
		imgIdx.DeleteByNameTag(nameTag)

		log.Println("image", nameTag, "deleted")
	}
}
