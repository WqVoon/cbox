package utils

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"path"

	"github.com/wqvoon/cbox/pkg/log"
)

func GetObjFromJsonFile(filePath string, obj interface{}) {
	if _, err := os.Stat(filePath); err != nil {
		log.Errorln("faild to stat json file, err:", err)
	}

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Errorln("faild to read json file, err:", err)
	}

	if err := json.Unmarshal(data, obj); err != nil {
		log.Errorln("faild to unmarshal json file, err:", err)
	}
}

func SaveObjToJsonFile(filePath string, obj interface{}) {
	data, err := json.Marshal(obj)
	if err != nil {
		log.Errorln("faild to marshal obj to json, err:", err)
	}

	if err := ioutil.WriteFile(filePath, data, 0644); err != nil {
		log.Errorln("faild to write json obj to file, err:", err)
	}
}

func CreateDirIfNotExist(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err = os.MkdirAll(path, 0755); err != nil {
			log.Errorln("faild to create directory, err:", err)
		}
	}
}

func WriteFileIfNotExist(filePath string, content []byte) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		dirPath := path.Dir(filePath)
		if dirPath != "." {
			CreateDirIfNotExist(dirPath)
		}

		file, err := os.Create(filePath)
		if err != nil {
			log.Errorln("faild to create file, err:", err)
		}

		if _, err = file.Write(content); err != nil {
			log.Errorln("faild to write file, err:", err)
		}
	}
}

func CreateDirWithExclusive(path string) {
	if _, err := os.Stat(path); err == nil {
		log.Errorln(path, "has existed, so can not re-create")
	}

	if err := os.MkdirAll(path, 0755); err != nil {
		log.Errorln("faild to create directory, err:", err)
	}
}

func CopyDir(from, to string) {
	if _, err := os.Stat(from); err != nil {
		log.Errorln("faild to stat", from, "err:", err)
	}

	CreateDirIfNotExist(to)

	cmd := exec.Command("cp", "-R", from, to)

	if err := cmd.Run(); err != nil {
		log.Errorf("faild to copy %s -> %s, err: %v\n", from, to, err)
	}
}
