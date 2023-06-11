package routes

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/DiegoAraujoJS/webdev-git-server/globals"
	"github.com/DiegoAraujoJS/webdev-git-server/pkg/utils"
	"github.com/go-git/go-git/v5"
)

func updateRepos () []error {
    wg := sync.WaitGroup{}
    globals.Get_commits_rw_mutex.Lock()
    defer globals.Get_commits_rw_mutex.Unlock()
	var errors []error
    for _, repo := range utils.Repositories {
        wg.Add(1)
        go func(repo *git.Repository) {
            if error := utils.ForceUpdateAllBranches(repo); error != nil && error != git.ErrRemoteNotFound && error != git.NoErrAlreadyUpToDate {
                errors = append(errors, error)
            }
            wg.Done()
        }(repo)
    }
    wg.Wait()
    return errors
}

func UpdateRepos(w http.ResponseWriter, r *http.Request) {
    errors := updateRepos()
	if len(errors) != 0 {
		json, _ := json.Marshal(errors)
		WriteError(&w, string(json), http.StatusInternalServerError)
		return
	}
    WriteResponseOk(&w, "Repositories updated successfully 👌")
}
