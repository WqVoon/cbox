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
	getObjFromJsonFile(idxFilePath, &ret)

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
	getObjFromJsonFile(manifestFilePath, &ret)

	return ret
}

func getObjFromJsonFile(filePath string, obj interface{}) {
	if _, err := os.Stat(filePath); err != nil {
		log.Fatalln("faild to stat json file, err:", err)
	}

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalln("can not read json file, err:", err)
	}

	if err := json.Unmarshal(data, obj); err != nil {
		log.Fatalln("can not unmarshal json file, err:", err)
	}
}
