package main

import (
	"bitbucket.org/uts-group/micro-tumbler/common"
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	configFileName = "config.json"
)

type mainSettings struct {
	Port            string `json:"port"`
	TokenStorage    string `json:"tokenStorage"`
	DatabaseConnect string `json:"databaseConnect"`
	KeyPath         string `json:"keyPath"`
}

func loadSettings() (*mainSettings, error) {
	configContent := []byte(os.Getenv(common.EnvConfig))
	if len(configContent) == 0 {
		configFlag := flag.String("c", configFileName, "Path to config file")
		flag.Parse()

		configFile, err := filepath.Abs(*configFlag)
		if err != nil {
			return nil, err
		}
		configContent, err = ioutil.ReadFile(configFile)
		if err != nil {
			return nil, err
		}
	}

	var Settings mainSettings
	err := json.Unmarshal(configContent, &Settings)
	if nil != err {
		return nil, err
	}

	return &Settings, nil
}
