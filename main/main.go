package main

import (
	"fmt"
	"log"

	"github.com/go-git/go-git/v5"
)

const A = 0

func main() {
	fmt.Println("Hello from main")
	repo, err := git.PlainOpen(".git")
	if err != nil {
		log.Fatal(err.Error())
	}
	logResult, err := repo.Log(&git.LogOptions{})
	if err != nil {
		log.Fatal(err.Error())
	}
	for {
		l, err := logResult.Next()
		if err != nil {
			break
		}
		fmt.Println(l.String())
	}
}
