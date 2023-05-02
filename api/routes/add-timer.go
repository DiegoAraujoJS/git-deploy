package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	builddeploy "github.com/DiegoAraujoJS/webdev-git-server/pkg/build-deploy"
)

func AddTimer(w http.ResponseWriter, r *http.Request) {
    repo, branch, seconds := r.URL.Query().Get("repo"), r.URL.Query().Get("branch"), r.URL.Query().Get("seconds")

    fmt.Println(repo, branch, seconds)

    secs, err := strconv.Atoi(seconds)
    if err != nil || secs < 60 {
        w.Header().Set("Content-Type", "text")
        w.WriteHeader(http.StatusNotAcceptable)
        w.Write([]byte("Format of \"seconds\" not correct or either has to be > 60"))
        return
    }

    w.Header().Set("Content-Type", "text")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("ok"))

    builddeploy.AddTimer(&builddeploy.AutobuildConfig{
        Repo: repo,
        Seconds: secs,
        Branch: branch,
    })
}

func GetTimers(w http.ResponseWriter, r *http.Request) {
    response, _ := json.Marshal(builddeploy.ActiveTimers)

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(response)
}
