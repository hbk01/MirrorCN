package main

import (
	"encoding/json"
	"io/ioutil"
)

func Parse(path string) (config Config) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}
	return config
}
