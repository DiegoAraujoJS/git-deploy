package routes

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/DiegoAraujoJS/webdev-git-server/pkg/navigation"
)

type FullResponse struct {
	*navigation.BranchResponse
}

func GetAllCommits(w http.ResponseWriter, r *http.Request) {
	repo := r.URL.Query().Get("repo")

	response, err := json.Marshal(&FullResponse{
		BranchResponse: navigation.GetAllCommits(repo),
	})
	if err != nil {
        log.Println("Error while getting release versions -> "+err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
