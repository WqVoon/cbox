package utils

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"

	"github.com/wqvoon/cbox/pkg/log"
)

func PathIsExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

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
		defer file.Close()

		if len(content) == 0 {
			return
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

func CopyDirContent(from, to string) {
	if _, err := os.Stat(from); err != nil {
		log.Errorln("faild to stat", from, "err:", err)
	}

	CreateDirIfNotExist(to)

	cmd := exec.Command("cp", "-r", from, to)

	if output, err := cmd.CombinedOutput(); err != nil {
		log.Println(string(output))
		log.Errorf("faild to copy %s -> %s, err: %v\n", from, to, err)
	}
}

func CopyFile(fromPath, toPath string) {
	toFile, err := os.OpenFile(toPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		log.Errorln("faild to create file", toPath, "err:", err)
	}

	fromFile, err := os.OpenFile(fromPath, os.O_RDONLY, 0777)
	if err != nil {
		log.Errorln("faild to create file", fromPath, "err:", err)
	}

	if _, err := io.Copy(toFile, fromFile); err != nil {
		log.Errorf("failed to copy %s -> %s, err: %v\n", fromPath, toPath, err)
	}
}
