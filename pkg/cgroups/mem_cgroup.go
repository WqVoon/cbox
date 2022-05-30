package cgroups

import (
	"bytes"
	"io/ioutil"
	"path"
	"strconv"

	"github.com/wqvoon/cbox/pkg/log"
	"github.com/wqvoon/cbox/pkg/utils"
)

const (
	MemLimitInBytes = "memory.limit_in_bytes"
	MemUsageInBytes = "memory.usage_in_bytes"
)

func GetMemCGroupByPath(pathname string) *MemCGroup {
	baseCGroup := BaseCGroup(pathname)

	cg := &MemCGroup{BaseCGroup: baseCGroup}
	cg.ValidOrDie()

	return cg
}

// 限制能够使用的内存大小，暂不考虑 swap
type MemCGroup struct {
	BaseCGroup
}

// 限制当前 CGroup 中的进程能够使用多少内存，以 MiB 为单位
func (c *MemCGroup) SetMemLimit(mem int) {
	const MiB = 1024 * 1024

	limitFilePath := path.Join(c.GetDirPath(), MemLimitInBytes)
	memSize := mem * MiB

	err := ioutil.WriteFile(limitFilePath, []byte(strconv.Itoa(memSize)), 0644)
	if err != nil {
		log.Errorln("failed to set mem limit, err:", err)
	}
}

// 获取当前 CGroup 对内存的限制，以 Byte 为单位
func (c *MemCGroup) GetMemLimit() int {
	limitFilePath := path.Join(c.GetDirPath(), MemLimitInBytes)
	limitValBytes, err := ioutil.ReadFile(limitFilePath)
	if err != nil {
		log.Errorln("failed to read", MemLimitInBytes)
	}

	limitValBytes = bytes.TrimSpace(limitValBytes)

	limitValNum, err := strconv.Atoi(string(limitValBytes))
	if err != nil {
		log.Errorln("failed to parse", MemLimitInBytes)
	}

	return limitValNum
}

// 获取当前 CGroup 的内存用量，以 Byte 为单位
func (c *MemCGroup) GetMemUsage() int {
	usageFilePath := path.Join(c.GetDirPath(), MemUsageInBytes)
	usageValBytes, err := ioutil.ReadFile(usageFilePath)
	if err != nil {
		log.Errorln("failed to read", MemUsageInBytes)
	}

	usageValBytes = bytes.TrimSpace(usageValBytes)

	usageValNum, err := strconv.Atoi(string(usageValBytes))
	if err != nil {
		log.Errorln("failed to parse", MemUsageInBytes)
	}

	return usageValNum
}

func (c *MemCGroup) IsValid() bool {
	checkedFile := path.Join(c.GetDirPath(), MemLimitInBytes)

	if !utils.PathIsExist(checkedFile) {
		return false
	}

	return c.BaseCGroup.IsValid()
}

func (c *MemCGroup) GetOrCreateSubCGroup(name string) *MemCGroup {
	c.ValidOrDie()

	subCGroupPath := path.Join(c.GetDirPath(), name)

	utils.CreateDirIfNotExist(subCGroupPath)

	subCGroup := GetMemCGroupByPath(subCGroupPath)
	subCGroup.SetNotifyOnRelease(true)

	return subCGroup
}

func (c *MemCGroup) GetDirPath() string { return string(c.BaseCGroup) }

func (c *MemCGroup) GetType() string { return "mem cgroup" }
