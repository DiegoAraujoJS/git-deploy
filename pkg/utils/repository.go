package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/go-git/go-git/v5"
)

type Config struct {
	Directory string
}

var (
	Repository *git.Repository
    config Config
)


func Connect() {
    content, err := ioutil.ReadFile("./config.json")

    if err != nil {
        log.Fatal(err.Error())
    }

    err = json.Unmarshal(content, &config)

    if err != nil {
        log.Fatal(err.Error())
    }

	r, err := git.PlainOpen(config.Directory)
	if err != nil {
		log.Fatal(err.Error())
	}
	Repository = r
}

func GetRepository() *git.Repository {
	return Repository
}
