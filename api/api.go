package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/DiegoAraujoJS/webdev-git-server/api/routes"
)

const PORT = "3001"

func ListenAndServe() {
    http.HandleFunc("/getRepos", routes.GetRepos)
	http.HandleFunc("/getTags", routes.GetReleaseVersions)
	http.HandleFunc("/checkout", routes.CheckoutBranch)
    fmt.Println("Listening on port " + PORT)
	err := http.ListenAndServe(":"+PORT, nil)
	if err != nil {
		log.Fatal(err.Error())
	}
}
