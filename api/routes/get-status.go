package routes

import (
	"encoding/json"
	"net/http"
	"strconv"

	builddeploy "github.com/DiegoAraujoJS/webdev-git-server/pkg/build-deploy"
)

type StatusResponse struct {
    Finished bool
    Moment string
    Stdout string
    Stderr string
}

func GetStatus(w http.ResponseWriter, r *http.Request) {
    ID := r.URL.Query().Get("ID")
    int_ID, _ := strconv.Atoi(ID)
    action := builddeploy.ActionStatus[int_ID]

    json, _ := json.Marshal(StatusResponse {
        Stdout: action.Status.Stdout.String(),
        Stderr: action.Status.Stderr.String(),
        Moment: builddeploy.StatusDict[action.Status.Moment],
        Finished: action.Status.Finished,
    })

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(json)

}
