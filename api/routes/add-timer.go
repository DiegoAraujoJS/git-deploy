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
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusNotAcceptable)
        w.Write([]byte(`{"error": "Repository not found"}`))
        return
    }

    branch := r.URL.Query().Get("branch")
    _, err := utils.GetBranch(repo, branch)
    if err != nil {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusNotAcceptable)
        w.Write([]byte(`{"error": "Branch not found"}`))
        return
    }

    seconds := r.URL.Query().Get("seconds")
    secs, err := strconv.Atoi(seconds)
    if err != nil || secs < 60 {
        w.Header().Set("Content-Type", "text")
        w.WriteHeader(http.StatusNotAcceptable)
        w.Write([]byte(`{"error": "Format of \"seconds\" not correct or either has to be > 60"}`))
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

    if timer, ok := builddeploy.ActiveTimers[repo]; ok {

        timer.Timer.Stop()
        delete(builddeploy.ActiveTimers, repo)
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("ok"))

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
