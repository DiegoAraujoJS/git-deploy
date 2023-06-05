package routes

import (
	"net/http"

	"github.com/DiegoAraujoJS/webdev-git-server/pkg/utils"
)

func GetRepos(w http.ResponseWriter, r *http.Request) {
	var repos []string
	for _, v := range utils.ConfigValue.Directories {
		repos = append(repos, v.Name)
	}
    WriteResponseOk(&w, repos)
}
