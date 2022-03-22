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

	case "run": // create + start
		imageNameTag, containerName := args[1], args[2]
		c = container.CreateContainer(
			image.GetImage(utils.GetNameTag(imageNameTag)), containerName,
		)
		c.Start()
		time.Sleep(1 * time.Second)
		c.Exec(args[3:]...)

	case "create":
		imageNameTag, containerName := args[1], args[2]
		container.CreateContainer(
			image.GetImage(utils.GetNameTag(imageNameTag)), containerName,
		)

	case "start": // by name
		container.StartContainersByName(args[1:]...)

	case "exec": // by name, run entrypoint
		name := args[1]
		c = container.GetContainerByName(name)
		c.Exec(args[2:]...)

	case "stop": // by name
		container.StopContainersByName(args[1:]...)

	case "delete": // by name
		container.DeleteContainersByName(args[1:]...)

	case "ps":
		container.ListAllContainer()

	case "pull":
		image.Pull(utils.GetNameTag(args[1]))

	}

}
