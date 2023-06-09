package routes

import (
	"net/http"
	"sync"

	"github.com/DiegoAraujoJS/webdev-git-server/globals"
	"github.com/DiegoAraujoJS/webdev-git-server/pkg/utils"
	"github.com/go-git/go-git/v5"
)

func GetRepos(w http.ResponseWriter, r *http.Request) {
	var repos []string
	for _, v := range utils.ConfigValue.Directories {
		repos = append(repos, v.Name)
	}
    WriteResponseOk(&w, repos)

    go func() {
        wg := sync.WaitGroup{}
        globals.Get_commits_rw_mutex.Lock()
        defer globals.Get_commits_rw_mutex.Unlock()
        for _, repo := range utils.Repositories {
            wg.Add(1)
            go func(repo *git.Repository) {
                utils.ForceUpdateAllBranches(repo)
                wg.Done()
            }(repo)
        }
        wg.Wait()
    }()
}
