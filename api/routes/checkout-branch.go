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
        WriteError(&w, "Repository not found", http.StatusNotFound)
        return
    }

    action := &builddeploy.Action{
        ID: builddeploy.GenerateActionID(),
        Repo: r.URL.Query().Get("repo"),
        Hash: r.URL.Query().Get("commit"),
    }

    builddeploy.CheckoutBuildInsertChan <- action

    WriteResponseOk(&w, strconv.Itoa(action.ID))
}
