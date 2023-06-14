package routes

import (
	"net/http"
	"strconv"

	"github.com/DiegoAraujoJS/webdev-git-server/pkg/build-deploy"
	"github.com/DiegoAraujoJS/webdev-git-server/pkg/utils"
	"github.com/go-git/go-git/v5/plumbing"
)

type CheckoutResponse struct {
    Errors  map[string]string
    Version string
}

func CheckoutBranch(w http.ResponseWriter, r *http.Request) {
    repo, ok := utils.Repositories[r.URL.Query().Get("repo")]

    if !ok {
        WriteError(&w, "Repository not found", http.StatusNotFound)
        return
    }

    // In the line below we check that the commit belongs to the repository.
    ref, err := repo.CommitObject(plumbing.NewHash(r.URL.Query().Get("commit")))
    if err != nil {
        WriteError(&w, "Commit not found", http.StatusNotFound)
        return
    }

    action := &builddeploy.Action{
        ID: builddeploy.GenerateActionID(),
        Repo: r.URL.Query().Get("repo"),
        Hash: ref.Hash,
    }

    builddeploy.CheckoutBuildInsertChan <- action

    WriteResponseOk(&w, strconv.Itoa(action.ID))
}
