package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/DiegoAraujoJS/webdev-git-server/api/routes"
	"github.com/DiegoAraujoJS/webdev-git-server/pkg/utils"
)

const PORT = "3001"

func ListenAndServe() {
    http.HandleFunc("/getRepos", routes.GetRepos)
	http.HandleFunc("/getTags", routes.GetReleaseVersions)
	http.HandleFunc("/checkout", routes.CheckoutBranch)
    if utils.ConfigValue.Port == "" {
        utils.ConfigValue.Port = PORT
    }
    fmt.Println("Listening on port " + utils.ConfigValue.Port)
	err := http.ListenAndServe(":"+utils.ConfigValue.Port, nil)
	if err != nil {
		log.Fatal(err.Error())
	}
}
