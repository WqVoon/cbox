package flags

import (
	"flag"
	"log"
	"os"
	"path"
	"path/filepath"
)

/*
	!!! 这个文件不可以引入任何 cbox 内部的包，因为它可能被任何内部的包使用 !!!
*/

var (
	parsed      = false
	rootDirPath = flag.String("root_dir", "", "cbox root directory path (default $HOME/cbox-dir)")
	driverName  = flag.String("storage_driver", "", "use which storage driver (default raw_copy)")
	debug       = flag.Bool("debug", false, "show call stack when run failed (default false)")
	dnsFilePath = flag.String("dns_file_path", "", "dns configuration file path")
	volume      = flag.String("volume", "", "bind mount some volumes")

	// 列表中的每一项是一个长度为 2 的列表，其中第一个字符串是 hostPath，第二个是 containerPath
	parsedVolumes [][]string
)

func Init() {
	if !parsed {
		flag.Parse()
		prepareRootDirPath()
		parsed = true
	}
}

func GetRootDirPath() string {
	if rootDirPath == nil || *rootDirPath == "" {
		prepareRootDirPath()
	}
	return *rootDirPath
}

func IsDebugMode() bool {
	return *debug
}

func GetStorageDriver() string {
	return *driverName
}

func GetDNSFilePath() string {
	return *dnsFilePath
}

func GetVolumes() [][]string {
	return parsedVolumes
}

func GetVolume() string {
	return *volume
}

func prepareRootDirPath() {
	const rootDirEnvName = "CBOX_ROOT_DIR"

	if rootDirPath == nil || *rootDirPath == "" {
		*rootDirPath = os.Getenv(rootDirEnvName)
	}

	if rootDirPath == nil || *rootDirPath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.Fatalln("faild to get user home dir, err:", err)
		}

		*rootDirPath = path.Join(homeDir, "cbox-dir")
	}

	// rootdir 必须是绝对路径
	absPath, err := filepath.Abs(*rootDirPath)
	if err != nil {
		log.Fatalln("failed to get absolute path from", *rootDirPath)
	}
	*rootDirPath = absPath
}
