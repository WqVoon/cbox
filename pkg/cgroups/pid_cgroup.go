package cgroups

import (
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
