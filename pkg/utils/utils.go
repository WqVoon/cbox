package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
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
