package cmd

import (
	"flag"
	"time"

	"github.com/wqvoon/cbox/pkg/container"
	"github.com/wqvoon/cbox/pkg/image"
	"github.com/wqvoon/cbox/pkg/log"
	"github.com/wqvoon/cbox/pkg/utils"
)

type Cmd struct {
	name        string // 命令名，用于匹配和 help
	description string // 命令描述，用于 help
}

type CmdSet struct {
	name2Fn map[string]func(args []string)
	cmds    []*Cmd
}

// 保存所有被注册的命令
var cmdSet CmdSet

func RegisterCmd(name, description string, fn func(args []string)) {
	if cmdSet.name2Fn == nil {
		cmdSet.name2Fn = make(map[string]func(args []string))
	}

	cmdSet.name2Fn[name] = fn
	cmdSet.cmds = append(cmdSet.cmds, &Cmd{name, description})
}

func Run() {
	cmdName := "help"

	args := flag.Args()
	if len(args) > 0 {
		cmdName, args = args[0], args[1:]
	}

	if fn, isIn := cmdSet.name2Fn[cmdName]; isIn {
		fn(args)
	} else {
		log.Errorln("unsupported cmd:", cmdName)
	}
}

func init() {
	RegisterCmd(
		"help",
		"显示帮助信息",
		func([]string) {
			for _, cmd := range cmdSet.cmds {
				log.Println(cmd.name, "\t", cmd.description)
			}
		})

	RegisterCmd(
		"test",
		"启动，运行一条龙服务（for test）",
		func(args []string) {
			c := container.CreateContainer(
				image.GetImage(utils.GetNameTag("alpine")), "test",
			)
			c.Start()
			// TODO: 这里先等1秒，后面整个更优雅的做法
			time.Sleep(1 * time.Second)
			c.Exec(args[0:]...)
		})

	RegisterCmd(
		"done",
		"停止，删除一条龙服务（for test）",
		func(args []string) {
			c := container.GetContainerByName("test")
			c.Stop()
			c.Delete()
		})

	RegisterCmd(
		"run",
		"create + start，默认运行 entrypoint，命令格式 `cbox run <IMAGE> <CONTAINER> [...args]`",
		func(args []string) {
			if len(args) < 2 {
				log.Errorln("malformed command, run `cbox help` for more info")
			}

			imageNameTag, containerName := args[0], args[1]
			c := container.CreateContainer(
				image.GetImage(utils.GetNameTag(imageNameTag)), containerName,
			)
			c.Start()
			time.Sleep(1 * time.Second)
			c.Exec(args[2:]...)
		})

	RegisterCmd(
		"create",
		"创建容器，命令格式 `cbox create <IMAGE> <CONTAINER>`",
		func(args []string) {
			if len(args) != 2 {
				log.Errorln("malformed command, run `cbox help` for more info")
			}

			imageNameTag, containerName := args[0], args[1]
			container.CreateContainer(
				image.GetImage(utils.GetNameTag(imageNameTag)), containerName,
			)
		})

	RegisterCmd(
		"start",
		"启动容器，命令格式 `cbox start <CONTAINER>`",
		func(args []string) {
			container.StartContainersByName(args...)
		})

	RegisterCmd(
		"exec",
		"在容器环境下执行命令, 默认运行 entrypoint，命令格式 `cbox exec <CONTAINER> [...args]`",
		func(args []string) {
			if len(args) == 0 {
				log.Errorln("malformed command, run `cbox help` for more info")
			}

			name := args[0]
			c := container.GetContainerByName(name)
			c.Exec(args[1:]...)
		})

	RegisterCmd(
		"stop",
		"停止容器，命令格式 `cbox stop <CONTAINER> [...<CONTAINER>]`",
		func(args []string) {
			container.StopContainersByName(args...)
		})

	RegisterCmd(
		"rm",
		"删除容器，命令格式 `cbox rm <CONTAINER> [...<CONTAINER>]`",
		func(args []string) {
			container.DeleteContainersByName(args...)
		})

	RegisterCmd(
		"ps",
		"列出所有的容器，命令格式 `cbox ps`",
		func(args []string) {
			container.ListAllContainer()
		})

	RegisterCmd(
		"pull",
		"拉取镜像到本地，命令格式 `cbox pull <CONTAINER>`",
		func(args []string) {
			if len(args) != 1 {
				log.Errorln("malformed command, run `cbox help` for more info")
			}

			image.Pull(utils.GetNameTag(args[0]))
		})

	RegisterCmd(
		"images",
		"列出本地所有镜像，命令格式 `cbox images`",
		func(args []string) {
			image.ListAllImage()
		})

	RegisterCmd(
		"rmi",
		"删除本地的镜像，命令格式 `cbox rmi <IMAGE> [...<IMAGE>]`",
		func(args []string) {
			image.DeleteImagesNyName(args...)
		})
}
