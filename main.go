package main

import (
	"flag"
	"time"

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

	args := flag.Args()
	cmd := args[0]

	var c *container.Container
	switch cmd {
	case "test": // 启动，运行一条龙服务（for test）
		c = container.CreateContainer(
			image.GetImage(utils.GetNameTag("alpine")), "test",
		)
		c.Start()
		// TODO: 这里先等1秒，后面整个更优雅的做法
		time.Sleep(1 * time.Second)
		c.Exec(args[1:]...)

	case "done": // 停止，删除一条龙服务（for test）
		c = container.GetContainerByName("test")
		c.Stop()
		c.Delete()

	case "run": // create + start，默认运行 entrypoint，命令格式 `cbox run <IMAGE> <CONTAINER> [...args]`
		imageNameTag, containerName := args[1], args[2]
		c = container.CreateContainer(
			image.GetImage(utils.GetNameTag(imageNameTag)), containerName,
		)
		c.Start()
		time.Sleep(1 * time.Second)
		c.Exec(args[3:]...)

	case "create": // 创建容器，命令格式 `cbox create <IMAGE> <CONTAINER>`
		imageNameTag, containerName := args[1], args[2]
		container.CreateContainer(
			image.GetImage(utils.GetNameTag(imageNameTag)), containerName,
		)

	case "start": // 启动容器，命令格式 `cbox start <CONTAINER>`
		container.StartContainersByName(args[1:]...)

	case "exec": // 在容器环境下执行命令, 默认运行 entrypoint，命令格式 `cbox exec <CONTAINER> [...args]`
		name := args[1]
		c = container.GetContainerByName(name)
		c.Exec(args[2:]...)

	case "stop": // 停止容器，命令格式 `cbox stop <CONTAINER> [...<CONTAINER>]`
		container.StopContainersByName(args[1:]...)

	case "rm": // 删除容器，命令格式 `cbox rm <CONTAINER> [...<CONTAINER>]`
		container.DeleteContainersByName(args[1:]...)

	case "ps": // 列出所有的容器，命令格式 `cbox ps`
		container.ListAllContainer()

	case "pull": // 拉取镜像到本地，命令格式 `cbox pull <CONTAINER>`
		image.Pull(utils.GetNameTag(args[1]))

	case "images": // 列出本地所有镜像，命令格式 `cbox images`
		image.ListAllImage()

	case "rmi": // 删除本地的镜像，命令格式 `cbox rmi <IMAGE> [...<IMAGE>]`
		image.DeleteImagesNyName(args[1:]...)
	}

}
