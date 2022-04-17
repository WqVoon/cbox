package flags

import (
	"flag"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
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
		prepareDriverName()
		prepareDNSFilePath()
		prepareVolume()

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

func prepareDriverName() {
	const driverEnvName = "CBOX_STORAGE_DRIVER"
	const defaultDriverName = "raw_copy"

	// 命令行的 flag 优先
	if driverName != nil && len(*driverName) != 0 {
		return
	}

	// 否则先看环境变量，这里不用 LookupEnv 是为了避免拿到空值
	*driverName = os.Getenv(driverEnvName)

	// 如果环境变量还没有，就设成默认值
	if driverName == nil || len(*driverName) == 0 {
		*driverName = defaultDriverName
	}
}

func prepareDNSFilePath() {
	var err error
	const defaultPath = "/etc/resolv.conf"

	// 说明没有设置过这个 flag，此时应该设置默认值或者退出，否则下面的 Abs 会出问题
	if *dnsFilePath == "" {
		if _, err := os.Stat(defaultPath); err == nil {
			*dnsFilePath = defaultPath
		}
		return
	}

	*dnsFilePath, err = filepath.Abs(*dnsFilePath)

	if err != nil {
		log.Fatalln("can not convert dnsFilePath to abs path")
	}
}

func prepareVolume() {
	if volume == nil || *volume == "" {
		return
	}

	volumes := strings.Split(*volume, ",")
	parsedVolumes = [][]string{}

	for _, v := range volumes {
		splitedV := strings.Split(v, ":")
		if len(splitedV) != 2 {
			log.Fatalln("error format of volume definition:", v)
		}

		hostPath, containerPath := splitedV[0], splitedV[1]

		hostPath, err := filepath.Abs(hostPath)
		if err != nil {
			log.Fatal("failed to convert hostPath to abs path, err:", err)
		}

		if !filepath.IsAbs(containerPath) {
			log.Fatal("containerPath must be abs path")
		}

		parsedVolumes = append(parsedVolumes, []string{hostPath, containerPath})
	}
}
