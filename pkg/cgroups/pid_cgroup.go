package cgroups

import (
	"bytes"
	"io/ioutil"
	"path"
	"strconv"

	"github.com/wqvoon/cbox/pkg/log"
	"github.com/wqvoon/cbox/pkg/utils"
)

func GetPIDCGroupByPath(pathname string) *PIDCGroup {
	baseCGroup := BaseCGroup(pathname)

	cg := &PIDCGroup{BaseCGroup: baseCGroup}
	cg.ValidOrDie()

	return cg
}

// 用来控制 Task 数量，Task 包括进程与线程
type PIDCGroup struct {
	BaseCGroup
}

// 限制当前 CGroup 中最多能够创建多少个 Task
func (c *PIDCGroup) SetTaskLimit(num int) {
	limitFilePath := path.Join(c.GetDirPath(), "pids.max")

	limitVal := "max"
	if num >= 0 {
		limitVal = strconv.Itoa(num)
	}

	err := ioutil.WriteFile(limitFilePath, []byte(limitVal), 0644)
	if err != nil {
		log.Errorln("failed to set task limit, err:", err)
	}
}

// 获取当前 CGroups 最大允许创建多少个 Task，如果为 max 那么返回 -1
func (c *PIDCGroup) GetTaskLimit() int {
	limitFilePath := path.Join(c.GetDirPath(), "pids.max")
	limitValBytes, err := ioutil.ReadFile(limitFilePath)
	if err != nil {
		log.Errorln("failed to read pids.max, err:", err)
	}

	limitValBytes = bytes.TrimSpace(limitValBytes)

	limitValStr := string(limitValBytes)
	if limitValStr == "max" {
		return -1
	}

	limitValNum, err := strconv.Atoi(limitValStr)
	if err != nil {
		log.Errorln("failed to parse pids.max, err:", err)
	}

	return limitValNum
}

// 获取当前 CGroups 中有多少个 Task
func (c *PIDCGroup) GetCurrentTaskNum() int {
	currentFilePath := path.Join(c.GetDirPath(), "pids.current")
	currentValBytes, err := ioutil.ReadFile(currentFilePath)
	if err != nil {
		log.Errorln("failed to read pids.current, err:", err)
	}

	currentValBytes = bytes.TrimSpace(currentValBytes)

	currentValNum, err := strconv.Atoi(string(currentValBytes))
	if err != nil {
		log.Errorln("failed to parse pids.max, err:", err)
	}

	return currentValNum
}

// 当前 CGroups 是否还可以新增 Task
func (c *PIDCGroup) CanJoinTask() bool {
	taskLimit := c.GetTaskLimit()
	if taskLimit == -1 {
		return true
	}

	return c.GetCurrentTaskNum() < taskLimit
}

func (c *PIDCGroup) IsValid() bool {
	// 只检查 pids.current 而不检查 pids.max
	// 因为 pids 的 root cgroup 中没有这个文件
	checkedFile := path.Join(c.GetDirPath(), "pids.current")

	if !utils.PathIsExist(checkedFile) {
		return false
	}

	return c.BaseCGroup.IsValid()
}

func (c *PIDCGroup) GetOrCreateSubCGroup(name string) *PIDCGroup {
	c.ValidOrDie()

	subCGroupPath := path.Join(c.GetDirPath(), name)

	utils.CreateDirIfNotExist(subCGroupPath)

	subCGroup := GetPIDCGroupByPath(subCGroupPath)
	subCGroup.SetNotifyOnRelease(true)

	return subCGroup
}

func (c *PIDCGroup) GetDirPath() string { return string(c.BaseCGroup) }

func (c *PIDCGroup) GetType() string { return "pid cgroup" }
