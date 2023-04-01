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
    router := http.NewServeMux()

    router.HandleFunc("/getRepos", routes.GetRepos)
	router.HandleFunc("/getTags", routes.GetReleaseVersions)
	router.HandleFunc("/checkout", routes.CheckoutBranch)
    router.HandleFunc("/repoHistory", routes.GetRepoHistory)

    handler := routes.EnableCorsMiddleware(router)

    if utils.ConfigValue.Port == "" {
        utils.ConfigValue.Port = PORT
    }

    fmt.Println("Listening on port " + utils.ConfigValue.Port)
	err := http.ListenAndServe(":"+utils.ConfigValue.Port, handler)
	if err != nil {
		log.Fatal(err.Error())
	}
}
