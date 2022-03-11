package main

import (
	"github.com/wqvoon/cbox/pkg/container"
	"github.com/wqvoon/cbox/pkg/flags"
	"github.com/wqvoon/cbox/pkg/image"
	"github.com/wqvoon/cbox/pkg/log"
	"github.com/wqvoon/cbox/pkg/rootdir"
	"github.com/wqvoon/cbox/pkg/utils"
)

func main() {
	flags.ParseAll()

	log.Println("Hello cbox!")

	rootdir.Init()
	log.Println("successfully create root dir:", rootdir.GetRootPath())

	img := image.GetImage(utils.GetNameTag("hello-world"))
	log.Println(img)

	ctr := container.CreateContainer(img, "test")
	log.Println(ctr)

	ctr = container.GetContainerByName(ctr.Name)
	log.Println("ByName:", ctr)

	ctr = container.GetContainerByID(ctr.ID)
	log.Println("ByID:", ctr)

	ctr.Start()
}
