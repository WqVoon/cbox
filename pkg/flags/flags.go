package flags

import (
	"flag"
)

/*
	!!! 这个文件不可以引入任何 cbox 内部的包，因为它可能被任何内部的包使用 !!!
*/

var (
	rootDirPath   = flag.String("root_dir", "", "cbox root directory path (default $HOME/cbox-dir)")
	driverName    = flag.String("storage_driver", "", "use which storage driver (default raw_copy)")
	debug         = flag.Bool("debug", false, "show call stack when run failed (default false)")
	dnsFilePath   = flag.String("dns_file_path", "", "dns configuration file path")
	volume        = flag.String("volume", "", "bind mount some volumes")
	useDockerfile = flag.Bool("use-dockerfile", false, "build image by dockerfile rather than json file")
)

func Init() {
	flag.Parse()
}

func GetRootDirPath() string { return *rootDirPath }

func IsDebugMode() bool { return *debug }

func GetStorageDriver() string { return *driverName }

func GetDNSFilePath() string { return *dnsFilePath }

func GetVolume() string { return *volume }

func UseDockerfile() bool { return *useDockerfile }
