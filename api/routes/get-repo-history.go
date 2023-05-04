package routes

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/DiegoAraujoJS/webdev-git-server/database"
	"github.com/DiegoAraujoJS/webdev-git-server/pkg/utils"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type VersionChangeEventWithCommit struct {
    Hash        string
    CreatedAt   string
    Commit      *object.Commit
}

func getCommit(repo, hash string) (*object.Commit, error) {
    repoObj := utils.GetRepository(repo)
    commit, err := repoObj.CommitObject(plumbing.NewHash(hash))
    if err != nil {
        return nil, err
    }
    return commit, nil
}

// Uses the function database.SelectVersionChangeEvents to get all the version change events for a given repo. It builds a JSON that is a list of the same type as the return value of database.SelectVersionChangeEvents.
func GetRepoHistory(w http.ResponseWriter, r *http.Request) {
    name := r.URL.Query().Get("repo")
    versionChangeEvents, err := database.SelectVersionChangeEvents(name)
    if err != nil {
        log.Println(err.Error())
        return
    }
    var versionChangeEventsWithCommit []*VersionChangeEventWithCommit
    for _, versionChangeEvent := range versionChangeEvents {
        commit, err := getCommit(name, versionChangeEvent.Hash)
        if err != nil {
            log.Println(err.Error())
            return
        }
        versionChangeEventsWithCommit = append(versionChangeEventsWithCommit, &VersionChangeEventWithCommit{
            Hash: versionChangeEvent.Hash,
            CreatedAt: versionChangeEvent.CreatedAt,
            Commit: commit,
        })
    }
    var versionChangeEventsJSON []byte
    versionChangeEventsJSON, err = json.Marshal(versionChangeEventsWithCommit)
    if err != nil {
        log.Println(err.Error())
        return
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(versionChangeEventsJSON)
}
