package routes

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/DiegoAraujoJS/webdev-git-server/pkg/navigation"
	"github.com/DiegoAraujoJS/webdev-git-server/pkg/utils"
)

type FullResponse struct {
	*navigation.BranchResponse
}

func GetAllCommits(w http.ResponseWriter, r *http.Request) {
	_, ok := utils.Repositories[r.URL.Query().Get("repo")]

    if !ok {
        log.Println("Error while getting release versions -> Repository not found" + r.URL.Query().Get("repo"))
        WriteError(&w, "Repository not found", 403)
        return
    }

	response, err := json.Marshal(&FullResponse{
		BranchResponse: navigation.GetAllCommits(r.URL.Query().Get("repo")),
	})
	if err != nil {
        log.Println("Error while getting release versions -> "+err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
