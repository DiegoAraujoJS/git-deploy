package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	builddeploy "github.com/DiegoAraujoJS/webdev-git-server/pkg/build-deploy"
)

type StatusResponse struct {
    Errors []error
    Finished bool
    Moment string
    Stdout string
}

func GetStatus(w http.ResponseWriter, r *http.Request) {
    ID := r.URL.Query().Get("ID")
    int_ID, _ := strconv.Atoi(ID)
    action := builddeploy.ActionStatus[int_ID]

    stdout := action.Status.Stdout.String()
    fmt.Println(stdout)

    json, _ := json.Marshal(StatusResponse {
        Stdout: stdout,
        Errors: action.Status.Errors,
        Moment: builddeploy.StatusDict[action.Status.Moment],
        Finished: action.Status.Finished,
    })

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(json)

}
