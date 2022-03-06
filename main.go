package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path"
)

var (
	rootDir = flag.String("root_dir", "", "cbox root directory (default $HOME/cbox-dir)")
)

type ImageIdx map[ImageName]ImageEntity
type ImageEntity map[ImageVersion]ImageHash
type ImageName string
type ImageVersion string
type ImageHash string

func main() {
	log.SetFlags(0)
	flag.Parse()

	log.Println("Hello cbox!")

	initRootDir(rootDir)
	log.Println("successfully create root dir:", *rootDir)

	idx := getImageIdx(rootDir)
	log.Println("get idx:")
	for name, entity := range idx {
		log.Println("-", name)

		for version, hash := range entity {
			log.Println(" -", version, ":", hash)
		}
	}
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

func getImageIdx(rootDir *string) ImageIdx {
	var ret ImageIdx

	idxFilePath := path.Join(*rootDir, "images", "images.json")
	if _, err := os.Stat(idxFilePath); os.IsNotExist(err) {
		ioutil.WriteFile(idxFilePath, []byte("{}"), 0644)
		return ret
	}

	data, err := ioutil.ReadFile(idxFilePath)
	if err != nil {
		log.Fatalln("can not read idx file, err:", err)
	}

	if err := json.Unmarshal(data, &ret); err != nil {
		log.Fatalln("can not unmarshal idx file, err:", err)
	}

	return ret
}
