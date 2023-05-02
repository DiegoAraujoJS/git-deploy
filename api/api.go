package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/DiegoAraujoJS/webdev-git-server/api/routes"
	"github.com/DiegoAraujoJS/webdev-git-server/pkg/utils"
)

const PORT string = "3001"

func ListenAndServe() {
    router := http.NewServeMux()

    router.HandleFunc("/getRepos", routes.GetRepos)
	router.HandleFunc("/getTags", routes.GetAllCommits)
	router.HandleFunc("/checkout", routes.CheckoutBranch)
    router.HandleFunc("/repoHistory", routes.GetRepoHistory)
    router.HandleFunc("/updateRepos", routes.UpdateRepos)
    router.HandleFunc("/getStatus", routes.GetStatus)
    router.HandleFunc("/addTimer", routes.AddTimer)
    router.HandleFunc("/getTimers", routes.GetTimers)

    handler := routes.EnableCorsMiddleware(
        routes.VerifyPasswordMiddleware(router),
    )

    if utils.ConfigValue.Port == "" {
        utils.ConfigValue.Port = PORT
    }

    fmt.Println("Listening on port " + utils.ConfigValue.Port)
	err := http.ListenAndServe(":"+utils.ConfigValue.Port, handler)
	if err != nil {
		log.Fatal(err.Error())
	}
}
