package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type ServerConfig struct {
	Root	string `json:"root"`
	MaxImages	int `json:"max_images"`
	Port	int	`json:"port"`
}

func loadImageCacheConfig() (*ServerConfig, error) {
	file, err := os.Open("config.json")
	if err != nil {
		return nil, err
	}

	jsonbytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	conf := new(ServerConfig)

	err = json.Unmarshal(jsonbytes, conf)
	if err != nil {
		return nil, err
	}

	return conf, nil
}