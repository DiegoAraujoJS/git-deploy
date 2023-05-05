package routes

import (
	"net/http"
	"strconv"

	"github.com/DiegoAraujoJS/webdev-git-server/pkg/build-deploy"
	"github.com/DiegoAraujoJS/webdev-git-server/pkg/utils"
)

type CheckoutResponse struct {
    Errors  map[string]string
    Version string
}

func CheckoutBranch(w http.ResponseWriter, r *http.Request) {
    _, ok := utils.Repositories[r.URL.Query().Get("repo")]

    if !ok {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusNotAcceptable)
        w.Write([]byte(`{"error": "Repository not found"}`))
        return
    }

    action_id := builddeploy.GenerateActionID()
    builddeploy.CheckoutBuildInsertChan <- &builddeploy.Action{
        ID: action_id,
        Repo: r.URL.Query().Get("repo"),
        Hash: r.URL.Query().Get("commit"),
    }

    w.Header().Set("Content-Type", "text")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(strconv.Itoa(action_id)))
}
