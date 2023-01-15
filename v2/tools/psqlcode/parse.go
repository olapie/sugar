package main

import (
	"os"

	"code.olapie.com/log"
	"gopkg.in/yaml.v2"
)

func ParseYAML(filename string) []*RepoModel {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.S().Fatal(err)
	}

	var repos []*RepoModel
	err = yaml.Unmarshal(data, &repos)
	if err != nil {
		log.S().Fatal(err)
	}
	return repos
}
