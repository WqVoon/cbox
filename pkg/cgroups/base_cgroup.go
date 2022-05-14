package cgroups

import (
	"io/ioutil"
	"os"
	"path"
	"strconv"

	"github.com/wqvoon/cbox/pkg/log"
	"github.com/wqvoon/cbox/pkg/utils"
)

type BaseCGroup string

// 判断是否是一个有效的 cgroup 目录，判断依据为特征文件是否存在
func (c *BaseCGroup) IsValid() bool {
	dirPath := c.GetDirPath()
	checkedFiles := []string{
		"cgroup.procs",      // 属于这一 cgroups 的进程
		"notify_on_release", // 该 cgroups 退出时是否执行 release_agent
		"tasks",             // 属于这一 cgroups 的线程
		// 不检查 release_agent，因为只有 root cgroup 中才存在这一文件
	}

	for _, filename := range checkedFiles {
		if !utils.PathIsExist(path.Join(dirPath, filename)) {
			return false
		}
	}

	return true
}

func (c *BaseCGroup) ValidOrDie() {
	if !c.IsValid() {
		log.Errorf("%q is not a valid cgroup\n", *c)
	}
}

func (c *BaseCGroup) DeleteSelf() {
	dirPath := c.GetDirPath()

	c.ValidOrDie()

	if err := os.Remove(dirPath); err != nil {
		log.Errorf("failed to remove cgroup %q err: %v\n", dirPath, err)
	}
}

func (c *BaseCGroup) DeleteSubCGroup(name string) {
	dirPath := c.GetDirPath()
	subCGroupPath := path.Join(dirPath, name)

	c.ValidOrDie()

	if err := os.Remove(subCGroupPath); err != nil {
		log.Errorf("failed to remove cgroup %q err: %v\n", dirPath, err)
	}
}

func (c *BaseCGroup) JoinProcessToSelf(pid int) {
	dirPath := c.GetDirPath()

	c.ValidOrDie()

	cProcsPath := path.Join(dirPath, "cgroup.procs")

	if err := ioutil.WriteFile(cProcsPath, []byte(strconv.Itoa(pid)), 0644); err != nil {
		log.Errorf("failed to join process %d to cgroup %q, err: %v\n", pid, dirPath, err)
	}
}

func (c *BaseCGroup) SetNotifyOnRelease(setToTrue bool) {
	val := "0"
	if setToTrue {
		val = "1"
	}

	dirPath := c.GetDirPath()
	notifyConfigPath := path.Join(dirPath, "notify_on_release")

	if err := ioutil.WriteFile(notifyConfigPath, []byte(val), 0644); err != nil {
		log.Errorf("failed to set notify_no_release for cgroup %q, err: %v\n", dirPath, err)
	}
}

func (c *BaseCGroup) GetDirPath() string { return string(*c) }

func (c *BaseCGroup) GetType() string { return "base cgroup" }
