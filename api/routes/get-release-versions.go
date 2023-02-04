package routes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/DiegoAraujoJS/webdev-git-server/pkg/navigation"
	"github.com/DiegoAraujoJS/webdev-git-server/pkg/utils"
)

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Headers", "*")
}

type FullResponse struct {
	*navigation.BranchResponse
}

func GetReleaseVersions(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get release versions")
	enableCors(&w)
	content, _ := ioutil.ReadFile("./config.json")
	json.Unmarshal(content, &utils.ConfigValue)

	repo := r.URL.Query().Get("repo")
	response, err := json.Marshal(&FullResponse{
		BranchResponse: navigation.GetReleaseBranchesWithTheirVersioning(repo),
	})
	if err != nil {
		log.Println(err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
