package utils

import (
	"log"

	"github.com/go-git/go-git/v5"
)

var (
    Repository *git.Repository
)

func Connect() {
	r, err := git.PlainOpen("/Users/diegoaraujo/repos/utils_lenox")
	if err != nil {
		log.Fatal(err.Error())
	}
    Repository = r
}

func GetRepository() *git.Repository {
    return Repository
}
