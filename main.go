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

type ImageIdx map[string]ImageEntity

// image name -> image hash
type ImageEntity map[string]string

type Manifest []struct {
	Config   string
	RepoTags []string
	Layers   []string
}

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

	manifest := getManifest(rootDir, idx, "hello-world", "latest")
	log.Println("get manifest:")
	for idx, oneManifest := range manifest {
		log.Println("- manifest", idx)

		log.Println(" - config:", oneManifest.Config)
		log.Println(" - layers:", oneManifest.Layers)
		log.Println(" - repoTags:", oneManifest.RepoTags)
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

func getManifest(rootDir *string, idx ImageIdx, name string, version string) Manifest {
	var ret Manifest

	entity, isIn := idx[name]
	if !isIn {
		log.Fatalf("no such image in imageIdx: %s:%s\n", name, version)
	}

	hash, isIn := entity[version]
	if !isIn {
		log.Fatalf("no such image in imageIdx: %s:%s\n", name, version)
	}

	manifestFilePath := path.Join(*rootDir, "images", hash, "manifest.json")

	if _, err := os.Stat(manifestFilePath); err != nil {
		log.Fatalf("faild to get manifest of image %s:%s, err: %v\n", name, version, err)
	}

	data, err := ioutil.ReadFile(manifestFilePath)
	if err != nil {
		log.Fatalln("can not read idx file, err:", err)
	}

	if err := json.Unmarshal(data, &ret); err != nil {
		log.Fatalln("can not unmarshal idx file, err:", err)
	}

	return ret
}
