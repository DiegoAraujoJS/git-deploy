package routes

import (
	"net/http"
	"strconv"

	"github.com/DiegoAraujoJS/webdev-git-server/pkg/build-deploy"
)

type CheckoutResponse struct {
    Errors  map[string]string
    Version string
}

func CheckoutBranch(w http.ResponseWriter, r *http.Request) {
    repo := r.URL.Query().Get("repo")

    action_id := builddeploy.GenerateActionID()
    builddeploy.CheckoutBuildInsertChan <- &builddeploy.Action{
        ID: action_id,
        Repo: repo,
        Hash: r.URL.Query().Get("commit"),
    }

    w.Header().Set("Content-Type", "text")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(strconv.Itoa(action_id)))
}
