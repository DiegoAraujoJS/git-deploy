package routes

import (
	"net/http"
	"strconv"

	builddeploy "github.com/DiegoAraujoJS/webdev-git-server/pkg/build-deploy"
	"github.com/DiegoAraujoJS/webdev-git-server/pkg/utils"
)

func AddTimer(w http.ResponseWriter, r *http.Request) {
    repo, ok := utils.Repositories[r.URL.Query().Get("repo")]
    if !ok {
        WriteError(&w, "Repository " + r.URL.Query().Get("repo") + " not found", http.StatusNotAcceptable)
        return
    }

    if _, err := utils.GetBranch(repo, r.URL.Query().Get("branch")); err != nil {
        WriteError(&w, "Branch " + r.URL.Query().Get("branch") + " not found", http.StatusNotAcceptable)
        return
    }

    if secs, err := strconv.Atoi(r.URL.Query().Get("seconds")); err == nil && secs >= 60 {

        builddeploy.AddTimer(&builddeploy.AutobuildConfig{
            Repo: r.URL.Query().Get("repo"),
            Seconds: secs,
            Branch: r.URL.Query().Get("branch"),
        })
        WriteResponseOk(&w, "Timer added")
        return
    }

    WriteError(&w, "Format of \"seconds\" not correct or either has to be >= 60", http.StatusNotAcceptable)
}

func DeleteTimer(w http.ResponseWriter, r *http.Request) {
    repo := r.URL.Query().Get("repo")

    if _, ok := builddeploy.ActiveTimers[repo]; ok {
        builddeploy.DeleteTimer(repo)
        WriteResponseOk(&w, "Timer deleted")
        return
    }

    WriteError(&w, "Timer not found", http.StatusNotAcceptable)
}

func GetTimers(w http.ResponseWriter, r *http.Request) {

    var configs = []*builddeploy.AutobuildConfig{}
    for _, timer := range builddeploy.ActiveTimers {
        configs = append(configs, timer.Config)
    }
    WriteResponseOk(&w, configs)
}
