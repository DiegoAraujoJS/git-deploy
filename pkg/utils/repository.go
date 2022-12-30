package utils

import (
	"log"

	"github.com/go-git/go-git/v5"
)

var Repo *git.Repository

func Connect() {
	r, err := git.PlainOpen(".git")
	if err != nil {
		log.Fatal(err.Error())
	}
	Repo = r
}

func GetRepo() *git.Repository {
	return Repo
}
