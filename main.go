package main

import (
	"os"

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

	cmd := os.Args[1]
	var c *container.Container
	switch cmd {
	case "create":
		imageNameTag, containerName := os.Args[2], os.Args[3]
		c = container.CreateContainer(image.GetImage(utils.GetNameTag(imageNameTag)), containerName)
		c.Start()

	case "get":
		containerName := os.Args[2]
		c = container.GetContainerByName(containerName)
		c.Start()
	}

}
