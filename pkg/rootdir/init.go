package rootdir

import (
	"os"
	"path"
	"path/filepath"

	"github.com/wqvoon/cbox/pkg/flags"
	"github.com/wqvoon/cbox/pkg/log"
	"github.com/wqvoon/cbox/pkg/utils"
)

var rootDirPath string

func Init() {
	initRootDirPath()
	initRootDirLayout()
}

// 初始化 rootDirPath 变量，该函数保证在正常退出时 rootDirPath 一定不为空
func initRootDirPath() {
	const rootDirEnvName = "CBOX_ROOT_DIR"

	rootDirPath = flags.GetRootDirPath()

	if rootDirPath == "" {
		rootDirPath = os.Getenv(rootDirEnvName)
	}

	if rootDirPath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.Errorln("faild to get user home dir, err:", err)
		}

		rootDirPath = path.Join(homeDir, "cbox-dir")
	}

	// rootdir 必须是绝对路径
	absPath, err := filepath.Abs(rootDirPath)
	if err != nil {
		log.Errorln("failed to get absolute path from", rootDirPath)
	}

	rootDirPath = absPath
}

// 创建 root_dir 中的基本文件
func initRootDirLayout() {
	if !filepath.IsAbs(rootDirPath) {
		log.Errorln("root_dir must be a absolute path")
	}

	subPaths := []string{
		path.Join("containers", "idx.json"),
		path.Join("images", "idx.json"),
	}

	data := []byte("{}")

	for _, subPath := range subPaths {
		path := path.Join(rootDirPath, subPath)

		utils.WriteFileIfNotExist(path, data)
	}

	utils.CreateDirIfNotExist(path.Join(rootDirPath, "tarballs"))
}
