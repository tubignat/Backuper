package common

import (
	"io/ioutil"
	"log"
	"os"
)

func IsExist(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func ReadFile(filename string) *[]byte {
	content, error := ioutil.ReadFile(filename)
	if error != nil {
		log.Panic(error)
	}
	return &content
}

func WriteFile(filename string, content *[]byte) {
	if error := ioutil.WriteFile(filename, *content, 0644); error != nil {
		log.Panic(error)
	}
}
