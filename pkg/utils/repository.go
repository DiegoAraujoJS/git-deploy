package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/go-git/go-git/v5"
)

type Config struct {
	Directory        string
	ClientDirectory  string
    BuildOutputFolder string
	BackendDirectory string
	LastBuild        struct{
        Version string
        Date string
    }
    IISDirectory string
}

var (
	Repository  *git.Repository
	ConfigValue Config
)

func Connect() {
	content, err := ioutil.ReadFile("./config.json")

	if err != nil {
		log.Fatal(err.Error())
	}

	err = json.Unmarshal(content, &ConfigValue)

	if err != nil {
		log.Fatal(err.Error())
	}

	if len(ConfigValue.Directory) == 0 || len(ConfigValue.ClientDirectory) == 0 || len(ConfigValue.BackendDirectory) == 0 {
		log.Fatal(err.Error())
	}

	r, err := git.PlainOpen(ConfigValue.Directory)
	if err != nil {
		log.Fatal(err.Error())
	}
	Repository = r
}

func GetRepository() *git.Repository {
	return Repository
}
