package cgroups

import (
	"io/ioutil"
	"path"
	"runtime"
	"strconv"

	"github.com/wqvoon/cbox/pkg/log"
	"github.com/wqvoon/cbox/pkg/utils"
)

func GetCPUCGroupByPath(pathname string) *CPUCGroup {
	baseCGroup := BaseCGroup(pathname)

	cg := &CPUCGroup{BaseCGroup: baseCGroup}
	cg.ValidOrDie()

	return cg
}

// 用来控制 cpu 核数
type CPUCGroup struct {
	BaseCGroup
}

// 限制当前 CGroup 中的进程能够使用多少个 cpu 核心
func (c *CPUCGroup) SetCPULimit(cpuNum int) {
	sysCpuNum := runtime.NumCPU()
	if cpuNum > sysCpuNum || cpuNum < 0 {
		cpuNum = sysCpuNum
	}

	dirPath := c.GetDirPath()
	cfsPeriodPath := path.Join(dirPath, "cpu.cfs_period_us")
	cfsQuotaPath := path.Join(dirPath, "cpu.cfs_quota_us")

	// 基数单位，由于 cpu cgroup 的单位是微妙，所以这个值表示 1 秒
	baseNum := 1000000

	if err := ioutil.WriteFile(cfsPeriodPath, []byte(strconv.Itoa(baseNum)), 0644); err != nil {
		log.Errorln("failed to set cfsPeriod, err:", err)
	}

	if err := ioutil.WriteFile(cfsQuotaPath, []byte(strconv.Itoa(cpuNum*baseNum)), 0644); err != nil {
		log.Errorln("failed to set cfsPeriod, err:", err)
	}
}

func (c *CPUCGroup) IsValid() bool {
	dirPath := c.GetDirPath()
	checkedFiles := []string{
		"cpu.cfs_period_us",
		"cpu.cfs_quota_us",
	}

	for _, filename := range checkedFiles {
		if !utils.PathIsExist(path.Join(dirPath, filename)) {
			return false
		}
	}

	return c.BaseCGroup.IsValid()
}

func (c *CPUCGroup) GetOrCreateSubCGroup(name string) *CPUCGroup {
	c.ValidOrDie()

	subCGroupPath := path.Join(c.GetDirPath(), name)

	utils.CreateDirIfNotExist(subCGroupPath)

	subCGroup := GetCPUCGroupByPath(subCGroupPath)
	subCGroup.SetNotifyOnRelease(true)

	return subCGroup
}

func (c *CPUCGroup) GetDirPath() string { return string(c.BaseCGroup) }

func (c *CPUCGroup) GetType() string { return "cpu cgroup" }
