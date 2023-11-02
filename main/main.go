package main

import (
	"github.com/DiegoAraujoJS/webdev-git-server/api"
	"github.com/DiegoAraujoJS/webdev-git-server/database"
	"github.com/DiegoAraujoJS/webdev-git-server/pkg/utils"
)

func init() {

	utils.Connect()
    database.InitDatabase()
}

func main() {
	api.ListenAndServe()
}
