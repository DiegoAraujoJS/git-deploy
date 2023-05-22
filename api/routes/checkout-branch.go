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
        WriteError(&w, "Repository not found", http.StatusNotAcceptable)
        return
    }

    action := &builddeploy.Action{
        Repo: r.URL.Query().Get("repo"),
        Hash: r.URL.Query().Get("commit"),
    }

    builddeploy.CheckoutBuildInsertChan <- action

    w.Header().Set("Content-Type", "text")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(strconv.Itoa(action.ID)))
}
