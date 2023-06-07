package routes

import (
	"net/http"
	"sync"

	"github.com/DiegoAraujoJS/webdev-git-server/pkg/navigation"
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
        navigation.Rw_lock.Lock()
        defer navigation.Rw_lock.Unlock()
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
