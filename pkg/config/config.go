package config

import (
	"github.com/wqvoon/cbox/pkg/rootdir"
	"github.com/wqvoon/cbox/pkg/utils"
)

var defaultConfig = struct {
	DriverName  string `json:"storage_driver"`
	DNSFilePath string `json:"dns_file_path"`
}{
	DriverName:  "raw_copy",
	DNSFilePath: "/etc/resolv.conf",
}

func Init() {
	utils.GetObjFromJsonFile(rootdir.GetConfigPath(), &defaultConfig)
}

func GetDriverName() string { return defaultConfig.DriverName }

func GetDNSFilePath() string { return defaultConfig.DNSFilePath }
