package rootdir

import (
	"path"

	"github.com/wqvoon/cbox/pkg/flags"
	"github.com/wqvoon/cbox/pkg/utils"
)

func Init() {
	rootPath := flags.GetRootDirPath()
	subPaths := []string{
		path.Join("containers", "idx.json"),
		path.Join("images", "idx.json"),
	}

	data := []byte("{}")

	for _, subPath := range subPaths {
		path := path.Join(rootPath, subPath)

		utils.WriteFileIfNotExist(path, data)
	}
}

var GetRootPath = flags.GetRootDirPath

//--------------- Image Start ---------------

// cbox-dir 中存放所有 images 的位置
// 内部有一个 idx.json 文件作为索引，以及以 imageHash 命名的文件夹，各文件夹内部是对应镜像的 Layout
func GetImageRootPath() string { return path.Join(GetRootPath(), "images") }

//---------- Image Root Start ----------

// image 的 idx.json 文件的路径
func GetImageIdxPath() string { return path.Join(GetImageRootPath(), "idx.json") }

// image 的 Layout，内部有 manifest.json 文件、image config 文件、image fs 文件夹
func GetImageLayoutPath(imageHash string) string { return path.Join(GetImageRootPath(), imageHash) }

//----- Image Layout Start -----

// image Layout 中的 manife.json 文件
func GetManifestPath(imageHash string) string {
	return path.Join(GetImageLayoutPath(imageHash), "manifest.json")
}

// image 的 config 文件，configFileName 可以从 manifest 中获取
func GetImageConfigPath(imageHash, configFileName string) string {
	return path.Join(GetImageLayoutPath(imageHash), configFileName)
}

//----- Image Layout End -----

//---------- Image Root End ----------

/*
//--------------- Image End ---------------






//--------------- Container Start ---------------
*/

// cbox-dir 中存放所有 containers 的位置
// 内部有一个 idx.json 文件作为索引，以及以 containerID 命名的文件夹，各文件夹内部是对应容器的 Layout
func GetContainerRootPath() string { return path.Join(GetRootPath(), "containers") }

//---------- Container Root Start ----------

// container 的 idx.json 文件
func GetContainerIdxPath() string { return path.Join(GetContainerRootPath(), "idx.json") }

// container 的 Layout，内部有 fs、ns 文件夹
func GetContainerLayoutPath(containerID string) string {
	return path.Join(GetContainerRootPath(), containerID)
}

//----- Container Layout Start -----

// container 的 fs，内部最重要的是 mnt 文件夹，此外根据不同的 StorageDriver 可能会有其他文件夹
func GetContainerFSPath(containerID string) string {
	return path.Join(GetContainerLayoutPath(containerID), "fs")
}

// Container FS Start

// 经过 StorageDriver 处理后的 mnt 文件夹，内部是容器的最终文件系统
func GetContainerMountPath(containerID string) string {
	return path.Join(GetContainerFSPath(containerID), "mnt")
}

// Container FS End

//----- Container Layout End -----

//---------- Container Root End ----------

//--------------- Container End ---------------
