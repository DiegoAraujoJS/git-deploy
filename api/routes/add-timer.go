package routes

import (
	"encoding/json"
	"net/http"
	"strconv"

	builddeploy "github.com/DiegoAraujoJS/webdev-git-server/pkg/build-deploy"
	"github.com/DiegoAraujoJS/webdev-git-server/pkg/utils"
)

func AddTimer(w http.ResponseWriter, r *http.Request) {
    repo, ok := utils.Repositories[r.URL.Query().Get("repo")]
    if !ok {
        WriteError(&w, "Repository not found", http.StatusNotAcceptable)
        return
    }

    branch := r.URL.Query().Get("branch")
    _, err := utils.GetBranch(repo, branch)
    if err != nil {
        WriteError(&w, "Branch not found", http.StatusNotAcceptable)
        return
    }

    seconds := r.URL.Query().Get("seconds")
    secs, err := strconv.Atoi(seconds)
    if err != nil || secs < 60 {
        WriteError(&w, "Format of \"seconds\" not correct or either has to be greater than or equal to 60", http.StatusNotAcceptable)
        return
    }

    w.Header().Set("Content-Type", "text")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("ok"))

    builddeploy.AddTimer(&builddeploy.AutobuildConfig{
        Repo: r.URL.Query().Get("repo"),
        Seconds: secs,
        Branch: branch,
    })
}

func DeleteTimer(w http.ResponseWriter, r *http.Request) {
    repo := r.URL.Query().Get("repo")

    if _, ok := builddeploy.ActiveTimers[repo]; ok {
        builddeploy.DeleteTimer(repo)
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("ok"))
        return
    }
    WriteError(&w, "Timer not found", http.StatusNotAcceptable)
}

func GetTimers(w http.ResponseWriter, r *http.Request) {

    var configs = map[string]*builddeploy.AutobuildConfig{}
    for k, v := range builddeploy.ActiveTimers {
        configs[k] = v.Config
    }

    response, _ := json.Marshal(configs)

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(response)
}
