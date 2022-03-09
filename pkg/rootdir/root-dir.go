package rootdir

import (
	"path"

	"github.com/wqvoon/cbox/pkg/flags"
	"github.com/wqvoon/cbox/pkg/utils"
)

func Init() {
	rootPath := flags.GetRootDirPath()
	subPaths := []string{"containers", "images"}

	for _, subPath := range subPaths {
		path := path.Join(rootPath, subPath)
		utils.CreateDirIfNotExist(path)
	}
}

var GetRootPath = flags.GetRootDirPath

func GetContainerRootPath() string { return path.Join(GetRootPath(), "containers") }

func GetImageRootPath() string { return path.Join(GetRootPath(), "images") }

func GetImageIdxPath() string { return path.Join(GetImageRootPath(), "images.json") }

func GetImagePath(imageHash string) string { return path.Join(GetImageRootPath(), imageHash) }

func GetManifestPath(imageHash string) string {
	return path.Join(GetImagePath(imageHash), "manifest.json")
}

func GetImageConfigPath(imageHash, configFileName string) string {
	return path.Join(GetImagePath(imageHash), configFileName)
}

func GetContainerPath(containerID string) string {
	return path.Join(GetContainerRootPath(), containerID)
}

func GetContainerFSPath(containerID string) string {
	return path.Join(GetContainerRootPath(), containerID, "fs")
}

func GetContainerMountPath(containerID string) string {
	return path.Join(GetContainerFSPath(containerID), "mnt")
}
