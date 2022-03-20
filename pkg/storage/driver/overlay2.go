package driver

import (
	"fmt"
	"path"
	"strings"

	"github.com/wqvoon/cbox/pkg/log"
	"github.com/wqvoon/cbox/pkg/utils"
	"golang.org/x/sys/unix"
)

var _ Interface = Register(&Overlay2{})

type Overlay2 struct{}

// layerPaths 中的 layer 来自于 manifest.json，也就是正序的
// 但是 overlay2 读 lower layer 时是倒序的，所以得反着来
func (o *Overlay2) Mount(dst string, layerPaths ...string) {
	// 做一个深拷贝，避免影响到其他功能
	layers := utils.ReverseStringSlice(
		utils.NewStringSlice(layerPaths...),
	)

	// upperdir 和 workdir 与 dst 在同一个目录中
	fsPath := path.Dir(dst)
	upperdirPath := path.Join(fsPath, "upperdir")
	workdirPath := path.Join(fsPath, "workdir")
	{
		utils.CreateDirIfNotExist(upperdirPath)
		utils.CreateDirIfNotExist(workdirPath)
	}

	mntOptions := fmt.Sprintf(
		"lowerdir=%s,upperdir=%s,workdir=%s",
		strings.Join(layers, ":"),
		upperdirPath,
		workdirPath,
	)

	if err := unix.Mount("cbox-overlay2-fs", dst, "overlay", 0, mntOptions); err != nil {
		log.Errorln("failed to mount overlay2 fs, err:", err)
	}
}

func (o *Overlay2) UnMount(dst string) {
	if err := unix.Unmount(dst, 0); err != nil {
		log.Errorln("failed to unmount overlay2 fs, err:", err)
	}
}

func (o *Overlay2) String() string { return "overlay2" }
