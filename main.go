package main

import (
	"flag"
	"log"
	"os"
	"path"
)

var (
	rootDir = flag.String("root_dir", "", "cbox root directory (default $HOME/cbox-dir)")
)

func main() {
	flag.Parse()

	log.Println("Hello cbox!")

	initRootDir(rootDir)
	log.Println("successfully create root dir:", *rootDir)
}

func initRootDir(rootDir *string) {
	if rootDir == nil || *rootDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.Fatalln("faild to get user home dir, err:", err)
		}

		*rootDir = path.Join(homeDir, "cbox-dir")
	}

	rootPath := *rootDir
	subPaths := []string{"containers", "images"}

	for _, subPath := range subPaths {
		path := path.Join(rootPath, subPath)

		if _, err := os.Stat(path); os.IsNotExist(err) {
			if err = os.MkdirAll(path, 0755); err != nil {
				log.Fatalln("faild to create directory, err:", err)
			}
		}
	}
}
