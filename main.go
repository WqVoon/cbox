package main

import (
	"flag"

	"github.com/wqvoon/cbox/pkg/container"
	"github.com/wqvoon/cbox/pkg/flags"
	"github.com/wqvoon/cbox/pkg/image"
	"github.com/wqvoon/cbox/pkg/log"
	"github.com/wqvoon/cbox/pkg/rootdir"
	"github.com/wqvoon/cbox/pkg/runtime"
	runtimeUtils "github.com/wqvoon/cbox/pkg/runtime/utils"
	"github.com/wqvoon/cbox/pkg/storage/driver"
	"github.com/wqvoon/cbox/pkg/utils"
)

func main() {
	if runtimeUtils.IsRuntimeMode() {
		runtime.Handle()
	}

	flags.ParseAll()

	log.Println("Hello cbox!")
	log.Printf("use %q as storage driver", driver.D)

	rootdir.Init()
	log.Println("successfully create root dir:", rootdir.GetRootPath())

	args := flag.Args()
	cmd := args[0]

	var c *container.Container
	switch cmd {
	case "test": // 一条龙服务
		c = container.CreateContainer(image.GetImage(utils.GetNameTag("alpine")), "test")
		c.Start()
		c.Stop()
		c.Delete()

	case "run": // create + start
		imageNameTag, containerName := args[1], args[2]
		c = container.CreateContainer(image.GetImage(utils.GetNameTag(imageNameTag)), containerName)
		c.Start()

	case "stop": // by name
		name := args[1]
		c = container.GetContainerByName(name)
		c.Stop()

	case "delete": // by name
		name := args[1]
		c = container.GetContainerByName(name)
		c.Delete()
	}

}
