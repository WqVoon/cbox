package config

import (
	"github.com/wqvoon/cbox/pkg/rootdir"
	"github.com/wqvoon/cbox/pkg/utils"
)

var defaultConfig = struct {
	DriverName  string        `json:"storage_driver"`
	DNSFilePath string        `json:"dns_file_path"`
	CGroup      cgroupConfig  `json:"cgroup"`
	Network     networkConfig `json:"network"`
}{
	DriverName:  "raw_copy",
	DNSFilePath: "/etc/resolv.conf",
	CGroup: cgroupConfig{
		Enable: false,
		Name:   "cbox",

		CPUCgroupPath: "/sys/fs/cgroup/cpu",
		CPULimit:      -1, // 取负数就会设置为系统 cpu 数量

		MemCgroupPath: "/sys/fs/cgroup/memory",
		MemLimit:      4 * 1024, // 默认限额设置为 4GiB

		PIDCgroupPath: "/sys/fs/cgroup/pids",
		TaskLimit:     -1, // 设置 -1 就是使用原本的值
	},
	Network: networkConfig{
		Enable:     true,
		Name:       "cbox0",
		BridgeCIDR: "172.29.0.1/16",
		IPRange: ipRange{
			Start: "172.29.0.150/16",
			End:   "172.29.0.200/16",
		},
	},
}

type cgroupConfig struct {
	Enable bool   `json:"enable"` // 是否启用 cgroup 隔离
	Name   string `json:"name"`   // 会作为 cgroup 文件夹的名字

	CPUCgroupPath string `json:"cpu_cgroup_path"` // cpu cgroup 的绝对路径
	CPULimit      int    `json:"cpu_limit"`       // 限制使用多少个 cpu 核心

	MemCgroupPath string `json:"mem_cgroup_path"` // mem cgroup 的绝对路径
	MemLimit      int    `json:"mem_limit"`       // 限制使用多少 MiB 的内存

	PIDCgroupPath string `json:"pid_cgroup_path"` // pid cgroup 的绝对路径
	TaskLimit     int    `json:"task_limit"`      // 限制最多创建多少个 Task
}

type networkConfig struct {
	Enable     bool    `json:"enable"`      // 是否启用网络隔离
	Name       string  `json:"name"`        // 网桥的名字
	BridgeCIDR string  `json:"bridge_cidr"` // 网桥的 CIDR
	IPRange    ipRange `json:"ip_range"`
}

type ipRange struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

func Init() {
	configPath := rootdir.GetConfigPath()

	if utils.PathIsExist(configPath) {
		utils.GetObjFromJsonFile(configPath, &defaultConfig)
	} else {
		// 如果配置文件不存在，那么将 defaultConfig 的内容写入
		utils.PrettySaveObjToJsonFile(configPath, &defaultConfig)
	}
}

func GetDriverName() string { return defaultConfig.DriverName }

func GetDNSFilePath() string { return defaultConfig.DNSFilePath }

func GetCgroupConfig() cgroupConfig { return defaultConfig.CGroup }

func GetNetworkConfig() networkConfig { return defaultConfig.Network }
