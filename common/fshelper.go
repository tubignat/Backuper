package common

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

func ReadFile(filename string) *[]byte {
	content, error := ioutil.ReadFile(filename)
	if error != nil {
		log.Panic(error)
	}
	return &content
}

func FromJSON(content []byte) interface{} {
	var object interface{}
	error := json.Unmarshal(content, &object)
	if error != nil {
		log.Fatal(error)
	}
	return object
}
