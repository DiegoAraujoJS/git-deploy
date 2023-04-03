package routes

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/DiegoAraujoJS/webdev-git-server/pkg/utils"
	"github.com/go-git/go-git/v5"
)

func UpdateRepos(w http.ResponseWriter, r *http.Request) {
    wg := sync.WaitGroup{}
    var errors []error
    for name, repo := range utils.Repositories {
        wg.Add(1)
        go func(repo *git.Repository, name string) {
            log.Println("Updating repo", name)
            if error := utils.ForceUpdateAllBranches(repo, &name); error != nil && error != git.ErrRemoteNotFound {
                log.Println("Error updating repo", name, error)
                errors = append(errors, error)
            }
            wg.Done()
        } (repo, name)
    }
    wg.Wait()
    if len(errors) != 0 {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusInternalServerError)
        json, _ := json.Marshal(errors)
        w.Write(json)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Repositories updated successfully"))
}
