package routes

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/DiegoAraujoJS/webdev-git-server/pkg/utils"
)

func GetRepos(w http.ResponseWriter, r *http.Request) {
	var repos []string
	for _, v := range utils.ConfigValue.Directories {
		repos = append(repos, v.Name)
	}
	response, err := json.Marshal(repos)
	if err != nil {
		log.Println(err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
