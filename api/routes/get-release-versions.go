package routes

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/DiegoAraujoJS/webdev-git-server/pkg/navigation"
)

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Headers", "*")
}

type FullResponse struct {
	*navigation.BranchResponse
}

func GetReleaseVersions(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)


	repo := r.URL.Query().Get("repo")
	response, err := json.Marshal(&FullResponse{
		BranchResponse: navigation.GetReleaseBranchesWithTheirVersioning(repo),
	})
	if err != nil {
        log.Println("Error while getting release versions -> "+err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
