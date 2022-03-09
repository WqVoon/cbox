package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

func GetObjFromJsonFile(filePath string, obj interface{}) {
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

func CreateDirIfNotExist(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err = os.MkdirAll(path, 0755); err != nil {
			log.Fatalln("faild to create directory, err:", err)
		}
	}
}

func CreateDirWithExclusive(path string) {
	if _, err := os.Stat(path); err == nil {
		log.Fatalln(path, "has existed, so can not re-create")
	}

	if err := os.MkdirAll(path, 0755); err != nil {
		log.Fatalln("faild to create directory, err:", err)
	}
}

func CopyDir(from, to string) {
	if _, err := os.Stat(from); err != nil {
		log.Fatalln("faild to stat", from, "err:", err)
	}

	CreateDirIfNotExist(to)

	cmd := exec.Command("cp", "-R", from, to)

	if err := cmd.Run(); err != nil {
		log.Fatalf("faild to copy %s -> %s, err: %v\n", from, to, err)
	}
}
