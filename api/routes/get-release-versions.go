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
    LastBuild *struct{
        Version string
        Date string
    } `json:"last_build"`
}

func GetReleaseVersions(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get release versions")
	enableCors(&w)
	content, _ := ioutil.ReadFile("./config.json")
    json.Unmarshal(content, &utils.ConfigValue)

	response, err := json.Marshal(&FullResponse{
        BranchResponse: navigation.GetReleaseBranchesWithTheirVersioning(),
        LastBuild: &utils.ConfigValue.LastBuild,
    })
	if err != nil {
		log.Println(err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
