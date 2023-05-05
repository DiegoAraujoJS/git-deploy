package routes

import (
	"encoding/json"
	"net/http"
	"strconv"

	builddeploy "github.com/DiegoAraujoJS/webdev-git-server/pkg/build-deploy"
)

type StatusResponse struct {
    Finished    bool
    Moment      string
    Stdout      string
    Stderr      string
}

func GetStatus(w http.ResponseWriter, r *http.Request) {
    ID := r.URL.Query().Get("ID")
    int_ID, err := strconv.Atoi(ID)

    if err != nil {
        if config, ok := builddeploy.ActiveTimers[ID]; ok {
            response, _ := json.Marshal(StatusResponse {
                Stdout: config.Config.Stdout.String(),
                Stderr: config.Config.Stderr.String(),
            })

            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusOK)
            w.Write(response)
            return
        }
    }

    action, ok := builddeploy.ActionStatus[int_ID]
    if action == nil || !ok {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"error": "Action not found"}`))
        return
    }

    response, err := json.Marshal(StatusResponse{
        Finished: action.Status.Finished,
        Moment: builddeploy.StatusDict[action.Status.Moment],
        Stdout: action.Status.Stdout.String(),
        Stderr: action.Status.Stderr.String(),
    })
    if err != nil {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("Error parsing response"))
        return
    }

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
