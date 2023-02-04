package main

import (
	// "github.com/DiegoAraujoJS/webdev-git-server/api"
	"github.com/DiegoAraujoJS/webdev-git-server/api"
	"github.com/DiegoAraujoJS/webdev-git-server/database"
	"github.com/DiegoAraujoJS/webdev-git-server/pkg/utils"
)

const A = 0

func init() {
	utils.Connect()
	database.CreateDatabase()
}

func main() {
	api.ListenAndServe()
}
