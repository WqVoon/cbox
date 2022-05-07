package builder

import (
	"github.com/wqvoon/cbox/pkg/image"
	"github.com/wqvoon/cbox/pkg/utils"
)

type BuildTask struct {
	// 镜像哈希，以 `cbox` 做前缀，长度为 12，自动生成
	ImageHash string

	// 镜像名，值应该是一个 nameTag 的字符串形式
	ImageNameStr string `json:"name"`
	// ImageNameStr 解析后的 NameTag
	ImageNameTag *utils.NameTag

	// 基础镜像，值应该是一个 nameTag 的字符串形式
	BaseImageStr string `json:"from"`
	// BaseImageStr 解析后的 NameTag
	BaseImageNameTag *utils.NameTag
	// BaseImageNameTag 对应的镜像
	BaseImage *image.Image

	// 环境变量，key 是变量名，value 是对应的值，会覆盖掉 BaseImage 中环境变量的同名值
	Env map[string]string `json:"env"`

	// 从宿主机复制文件/文件夹到镜像内
	Copy []CopyType `json:"copy"`

	// 需要被执行的 Cmd，其执行在所有的 Copy 结束后
	Cmds []Cmd `json:"cmds"`

	// 入口命令
	Entrypint Cmd `json:"entrypoint"`
}

type Cmd []string

// Src 是宿主机目录，可以是相对路径
// Dst 是镜像内目录，必须是绝对路径
type CopyType struct {
	Src string `json:"src"`
	Dst string `json:"dst"`
}

func LoadFromJsonFile(filePath string) *BuildTask {
	bt := &BuildTask{}
	utils.GetObjFromJsonFile(filePath, bt)
	return bt
}
