package main

import (
	"flag"
	"time"

	"github.com/wqvoon/cbox/pkg/cmd"
	"github.com/wqvoon/cbox/pkg/container"
	"github.com/wqvoon/cbox/pkg/flags"
	"github.com/wqvoon/cbox/pkg/image"
	"github.com/wqvoon/cbox/pkg/rootdir"
	"github.com/wqvoon/cbox/pkg/runtime"
	runtimeUtils "github.com/wqvoon/cbox/pkg/runtime/utils"
	"github.com/wqvoon/cbox/pkg/utils"
)

func main() {
	flags.ParseAll()
	rootdir.Init()

	if runtimeUtils.IsRuntimeMode() {
		runtime.Handle()
	}

	{
		cmd.RegisterCmd("test", func(args []string) {
			c := container.CreateContainer(
				image.GetImage(utils.GetNameTag("alpine")), "test",
			)
			c.Start()
			// TODO: 这里先等1秒，后面整个更优雅的做法
			time.Sleep(1 * time.Second)
			c.Exec(args[0:]...)
		}, "启动，运行一条龙服务（for test）")

		cmd.RegisterCmd("done", func(args []string) {
			c := container.GetContainerByName("test")
			c.Stop()
			c.Delete()
		}, "停止，删除一条龙服务（for test）")

		cmd.RegisterCmd("run", func(args []string) {
			imageNameTag, containerName := args[0], args[1]
			c := container.CreateContainer(
				image.GetImage(utils.GetNameTag(imageNameTag)), containerName,
			)
			c.Start()
			time.Sleep(1 * time.Second)
			c.Exec(args[2:]...)
		}, "create + start，默认运行 entrypoint，命令格式 `cbox run <IMAGE> <CONTAINER> [...args]`")

		cmd.RegisterCmd("create", func(args []string) {
			imageNameTag, containerName := args[0], args[1]
			container.CreateContainer(
				image.GetImage(utils.GetNameTag(imageNameTag)), containerName,
			)
		}, "创建容器，命令格式 `cbox create <IMAGE> <CONTAINER>`")

		cmd.RegisterCmd("start", func(args []string) {
			container.StartContainersByName(args...)
		}, "启动容器，命令格式 `cbox start <CONTAINER>`")

		cmd.RegisterCmd("exec", func(args []string) {
			name := args[0]
			c := container.GetContainerByName(name)
			c.Exec(args[1:]...)
		}, "在容器环境下执行命令, 默认运行 entrypoint，命令格式 `cbox exec <CONTAINER> [...args]`")

		cmd.RegisterCmd("stop", func(args []string) {
			container.StopContainersByName(args...)
		}, "停止容器，命令格式 `cbox stop <CONTAINER> [...<CONTAINER>]`")

		cmd.RegisterCmd("rm", func(args []string) {
			container.DeleteContainersByName(args...)
		}, "删除容器，命令格式 `cbox rm <CONTAINER> [...<CONTAINER>]`")

		cmd.RegisterCmd("ps", func(args []string) {
			container.ListAllContainer()
		}, "列出所有的容器，命令格式 `cbox ps`")

		cmd.RegisterCmd("pull", func(args []string) {
			image.Pull(utils.GetNameTag(args[0]))
		}, "拉取镜像到本地，命令格式 `cbox pull <CONTAINER>`")

		cmd.RegisterCmd("images", func(args []string) {
			image.ListAllImage()
		}, "列出本地所有镜像，命令格式 `cbox images`")

		cmd.RegisterCmd("rmi", func(args []string) {
			image.DeleteImagesNyName(args...)
		}, "删除本地的镜像，命令格式 `cbox rmi <IMAGE> [...<IMAGE>]`")
	}

	args := flag.Args()
	cmd.Run(args[0], args[1:])
}
