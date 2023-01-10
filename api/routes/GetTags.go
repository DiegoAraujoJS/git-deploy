package routes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/DiegoAraujoJS/webdev-git-server/pkg/navigation"
)

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func GetReleaseVersions(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get release versions")
	enableCors(&w)
	response, err := json.Marshal(navigation.GetReleaseBranchesWithTheirVersioning())
	if err != nil {
		log.Println(err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
