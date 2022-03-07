package rootdir

import (
	"log"
	"os"
	"path"

	"github.com/wqvoon/cbox/pkg/flags"
)

func Init() {
	rootPath := flags.GetRootDirPath()
	subPaths := []string{"containers", "images"}

	for _, subPath := range subPaths {
		path := path.Join(rootPath, subPath)

		if _, err := os.Stat(path); os.IsNotExist(err) {
			if err = os.MkdirAll(path, 0755); err != nil {
				log.Fatalln("faild to create directory, err:", err)
			}
		}
	}
}

var GetPath = flags.GetRootDirPath
