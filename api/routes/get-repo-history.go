package routes

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/DiegoAraujoJS/webdev-git-server/database"
)

// Uses the function database.SelectVersionChangeEvents to get all the version change events for a given repo. It builds a JSON that is a list of the same type as the return value of database.SelectVersionChangeEvents.
func GetRepoHistory(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
    name := r.URL.Query().Get("repo")
    versionChangeEvents, err := database.SelectVersionChangeEvents(name)
    if err != nil {
        log.Println(err.Error())
        return
    }
    var versionChangeEventsJSON []byte
    versionChangeEventsJSON, err = json.Marshal(versionChangeEvents)
    if err != nil {
        log.Println(err.Error())
        return
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(versionChangeEventsJSON)
}
