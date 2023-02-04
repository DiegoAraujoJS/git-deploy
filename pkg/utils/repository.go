package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/go-git/go-git/v5"
)

type Config struct {
	Directories []struct {
		Name      string
		Directory string
	}
}

var Repositories = map[string]*git.Repository{}
var ConfigValue Config

func Connect() {
	content, err := ioutil.ReadFile("./config.json")

	if err != nil {
		log.Fatal(err.Error())
	}

	err = json.Unmarshal(content, &ConfigValue)

	if err != nil {
		log.Fatal(err.Error())
	}

	for _, dir := range ConfigValue.Directories {
		r, err := git.PlainOpen(dir.Directory)
		if err != nil {
			log.Fatal(err.Error())
		}
		Repositories[dir.Name] = r
	}

}

func GetRepository(repo string) *git.Repository {
	return Repositories[repo]
}
