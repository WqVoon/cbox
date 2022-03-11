package flags

import (
	"flag"
	"log"
	"os"
	"path"
)

/*
	!!! 这个文件不可以引入任何 cbox 内部的包，因为它可能被任何内部的包使用 !!!
*/

var (
	parsed      = false
	rootDirPath = flag.String("root_dir", "", "cbox root directory path (default $HOME/cbox-dir)")
	driverName  = flag.String("storage_driver", "raw_copy", "use which storage driver")
	debug       = flag.Bool("debug", false, "show call stack when run failed")
)

func ParseAll() {
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

func prepareRootDirPath() {
	if rootDirPath == nil || *rootDirPath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.Fatalln("faild to get user home dir, err:", err)
		}

		*rootDirPath = path.Join(homeDir, "cbox-dir")
	}
}
