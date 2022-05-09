package rootdir

import (
	"path"
)

func GetRootDirPath() string { return rootDirPath }

//--------------- Image Start ---------------

// cbox-dir 中存放所有 images 的位置
// 内部有一个 idx.json 文件作为索引，以及以 imageHash 命名的文件夹，各文件夹内部是对应镜像的 Layout
func GetImageRootPath() string { return path.Join(GetRootDirPath(), "images") }

//---------- Image Root Start ----------

// image 的 idx.json 文件的路径
func GetImageIdxPath() string { return path.Join(GetImageRootPath(), "idx.json") }

// image 的 Layout，内部有 manifest.json 文件、image config 文件、image fs 文件夹
func GetImageLayoutPath(imageHash string) string { return path.Join(GetImageRootPath(), imageHash) }

//----- Image Layout Start -----

// 获取 Image 的 info 文件路径，内部保存 Image 的运行时信息
func GetImageInfoPath(imageHash string) string {
	return path.Join(GetImageLayoutPath(imageHash), "info")
}

// 获取 Image 的 FS 路径，解包后的文件就放在这里
func GetImageFsPath(imageHash, layerPath string) string {
	return path.Join(GetImageLayoutPath(imageHash), layerPath, "fs")
}

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
func GetContainerRootPath() string { return path.Join(GetRootDirPath(), "containers") }

//---------- Container Root Start ----------

// container 的 idx.json 文件
func GetContainerIdxPath() string { return path.Join(GetContainerRootPath(), "idx.json") }

// container 的 Layout，内部有 fs、ns 文件夹
func GetContainerLayoutPath(containerID string) string {
	return path.Join(GetContainerRootPath(), containerID)
}

//----- Container Layout Start -----

// container 的 info 文件，内部保存运行时的相关信息
func GetContainerInfoPath(containerID string) string {
	return path.Join(GetContainerLayoutPath(containerID), "info")
}

// container 的 ns 文件夹，内部绑定挂载了容器对应的 ns
func GetContainerNSPath(containerID string) string {
	return path.Join(GetContainerLayoutPath(containerID), "ns")
}

// container 的 fs，内部最重要的是 mnt 文件夹，此外根据不同的 StorageDriver 可能会有其他文件夹
func GetContainerFSPath(containerID string) string {
	return path.Join(GetContainerLayoutPath(containerID), "fs")
}

// Container FS Start

// 经过 StorageDriver 处理后的 mnt 文件夹，内部是容器的最终文件系统
func GetContainerMountPath(containerID string) string {
	return path.Join(GetContainerFSPath(containerID), "mnt")
}

// 获取 mnt 文件夹中 dns 配置文件的路径，当前仅支持 /etc/resolv.conf
func GetContainerDNSConfigPath(containerID string) string {
	return path.Join(GetContainerMountPath(containerID), "etc", "resolv.conf")
}

// 获取健康检查文件的路径，这个文件中记录了实时的不健康原因，根据 inContainerView 来区分是宿主机视角还是容器视角
func GetContainerHealthCheckInfoPath(containerID string, inContainerView bool) string {
	fileName := "/.cbox-health-check"
	if inContainerView {
		return fileName
	}

	return path.Join(GetContainerMountPath(containerID), fileName)
}

// Container FS End

//----- Container Layout End -----

//---------- Container Root End ----------

//--------------- Container End ---------------

// 获取 cbox-dir 中 tarballs 文件夹的路径
func GetTarballRootPath() string {
	return path.Join(GetRootDirPath(), "tarballs")
}

// 获取 tarballs 文件夹中某个 Image 的路径
func GetImageTarballPath(imgHash string) string {
	return path.Join(GetTarballRootPath(), imgHash, "image.tar")
}

func GetConfigPath() string { return path.Join(GetRootDirPath(), "config.json") }
