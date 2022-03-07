package main

import (
	"log"

	"github.com/wqvoon/cbox/pkg/flags"
	"github.com/wqvoon/cbox/pkg/image"
	"github.com/wqvoon/cbox/pkg/rootdir"
)

func main() {
	log.SetFlags(0)
	flags.ParseAll()

	log.Println("Hello cbox!")

	rootdir.Init()
	log.Println("successfully create root dir:", rootdir.GetPath())

	idx := image.GetIdx()
	log.Println("get idx:")
	for name, entity := range idx {
		log.Println("-", name)

		for version, hash := range entity {
			log.Println(" -", version, ":", hash)
		}
	}

	manifest := idx.GetManifest("hello-world:latest")
	log.Println("get manifest:")
	for idx, oneManifest := range manifest {
		log.Println("- manifest", idx)

		log.Println(" - config:", oneManifest.Config)
		log.Println(" - layers:", oneManifest.Layers)
		log.Println(" - repoTags:", oneManifest.RepoTags)
	}
}
