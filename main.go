package main

import (
	"flag"

	"github.com/wqvoon/cbox/pkg/container"
	"github.com/wqvoon/cbox/pkg/flags"
	"github.com/wqvoon/cbox/pkg/image"
	"github.com/wqvoon/cbox/pkg/log"
	"github.com/wqvoon/cbox/pkg/rootdir"
	"github.com/wqvoon/cbox/pkg/storage/driver"
	"github.com/wqvoon/cbox/pkg/utils"
)

func main() {
	flags.ParseAll()

	log.Println("Hello cbox!")
	log.Printf("use %q as storage driver", driver.D)

	rootdir.Init()
	log.Println("successfully create root dir:", rootdir.GetRootPath())

	args := flag.Args()
	cmd := args[0]

	var c *container.Container
	switch cmd {
	case "create":
		imageNameTag, containerName := args[1], args[2]
		c = container.CreateContainer(image.GetImage(utils.GetNameTag(imageNameTag)), containerName)
		c.Start()

	case "get":
		containerName := args[1]
		c = container.GetContainerByName(containerName)
		c.Start()
	}

}
