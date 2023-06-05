package routes

import (
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
            WriteResponseOk(&w, StatusResponse {
                Stdout: config.Config.Stdout.String(),
                Stderr: config.Config.Stderr.String(),
            })
            return
        }
        WriteError(&w, "No timer active for " + ID, http.StatusNotAcceptable)
        return
    }

    action, ok := builddeploy.ActionStatus[int_ID]
    if action == nil || !ok {
        WriteError(&w, "Action not found", http.StatusNotAcceptable)
        return
    }

    WriteResponseOk(&w, StatusResponse{
        Finished: action.Status.Finished,
        Moment: builddeploy.StatusDict[action.Status.Moment],
        Stdout: action.Status.Stdout.String(),
        Stderr: action.Status.Stderr.String(),
    })
}

func ClearStatus(w http.ResponseWriter, r *http.Request) {
    ID := r.URL.Query().Get("ID")
    _, err := strconv.Atoi(ID)

    if err != nil {
        if config, ok := builddeploy.ActiveTimers[ID]; ok {
            config.Config.Stdout.Reset()
            config.Config.Stderr.Reset()

            WriteResponseOk(&w, "Status cleared")
            return
        }
        WriteError(&w, "No timer active for " + ID, http.StatusNotAcceptable)
        return
    }
    WriteError(&w, "Not implemented", http.StatusNotImplemented)
}
