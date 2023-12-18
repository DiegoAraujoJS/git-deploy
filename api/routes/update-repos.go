package routes

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/DiegoAraujoJS/webdev-git-server/globals"
	"github.com/DiegoAraujoJS/webdev-git-server/pkg/utils"
	"github.com/go-git/go-git/v5"
)

func updateRepos () []error {
    globals.Get_commits_rw_mutex.Lock()
    defer globals.Get_commits_rw_mutex.Unlock()
    wg := sync.WaitGroup{}
	var errors []error
    for _, repo := range utils.Applications.GetAllRepos() {
        wg.Add(1)
        go func(repo *git.Repository) {
            defer wg.Done()
            defer func (){
                if r := recover(); r != nil {
                    errors = append(errors, fmt.Errorf(r.(string)))
                }
            }()
            if error := utils.ForceUpdateAllBranches(repo); error != nil && error != git.ErrRemoteNotFound && error != git.NoErrAlreadyUpToDate {
                errors = append(errors, error)
            }
        }(repo)
    }
    wg.Wait()
    return errors
}

func UpdateRepos(w http.ResponseWriter, r *http.Request) {
    updateRepos()
    WriteResponseOk(&w, "Repositories updated successfully 👌")
}
