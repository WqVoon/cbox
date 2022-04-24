package dns

import (
	"path/filepath"

	"github.com/wqvoon/cbox/pkg/flags"
	"github.com/wqvoon/cbox/pkg/log"
	"github.com/wqvoon/cbox/pkg/utils"
)

var dnsFilePath string

// 目前 Init 方法强制保证 dnsFilePath 变量不会为空
func Init() {
	candidateList := []string{
		flags.GetDNSFilePath(),
		"/etc/resolv.conf",
		"/var/run/systemd/resolve/resolv.conf",
	}

	for _, filePath := range candidateList {
		if filePath == "" || !utils.PathIsExist(filePath) {
			continue
		}

		var err error

		if dnsFilePath, err = filepath.Abs(filePath); err != nil {
			log.Errorln("can not convert dnsFilePath to abs path")
		} else {
			return
		}
	}

	log.Errorln("can not init dnsFilePath")
}

func GetDNSFilePath() string { return dnsFilePath }
