package main

import (
	"log"

	"github.com/wqvoon/cbox/pkg/container"
	"github.com/wqvoon/cbox/pkg/flags"
	"github.com/wqvoon/cbox/pkg/image"
	"github.com/wqvoon/cbox/pkg/rootdir"
	"github.com/wqvoon/cbox/pkg/utils"
)

func main() {
	log.SetFlags(0)
	flags.ParseAll()

	log.Println("Hello cbox!")

	rootdir.Init()
	log.Println("successfully create root dir:", rootdir.GetRootPath())

	img := image.GetImageFromLocal(utils.GetNameTag("hello-world"))
	log.Println(img)

	ctr := container.CreateContainer(img, "test")
	log.Println(ctr)

	ctr.Start()
}
