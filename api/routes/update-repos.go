package routes

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/DiegoAraujoJS/webdev-git-server/pkg/navigation"
	"github.com/DiegoAraujoJS/webdev-git-server/pkg/utils"
	"github.com/go-git/go-git/v5"
)

func UpdateRepos(w http.ResponseWriter, r *http.Request) {
	wg := sync.WaitGroup{}
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
	if len(errors) != 0 {
		json, _ := json.Marshal(errors)
		WriteError(&w, string(json), http.StatusInternalServerError)
		return
	}
    navigation.Rw_lock.Lock()
	navigation.All_commits = map[string][]*navigation.Commit{}
    navigation.Rw_lock.Unlock()
    WriteResponseOk(&w, "Repositories updated successfully 👌")
}
