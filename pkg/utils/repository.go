package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/go-git/go-git/v5"
)

// Config is the struct that holds the configuration for the git repositories. The configuration is read from a json file that is located in the root of the project.
type Config struct {
    Port string
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
