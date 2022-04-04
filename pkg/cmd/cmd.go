package cmd

import (
	"github.com/wqvoon/cbox/pkg/log"
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

func RegisterCmd(name string, fn func(args []string), description string) {
	if cmdSet.name2Fn == nil {
		cmdSet.name2Fn = make(map[string]func(args []string))
	}

	cmdSet.name2Fn[name] = fn
	cmdSet.cmds = append(cmdSet.cmds, &Cmd{name, description})
}

func Run(cmdName string, args []string) {
	if fn, isIn := cmdSet.name2Fn[cmdName]; isIn {
		fn(args)
	} else {
		log.Println("unsupported cmd:", cmdName)
	}
}

func init() {
	RegisterCmd("help", func([]string) {
		for _, cmd := range cmdSet.cmds {
			log.Println(cmd.name, "\t", cmd.description)
		}
	}, "显示帮助信息")
}
