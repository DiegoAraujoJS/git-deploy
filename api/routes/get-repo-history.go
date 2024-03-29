package routes

import (
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
    app := utils.Applications[repo]
    commit, err := app.Repo.CommitObject(plumbing.NewHash(hash))
    if err != nil {
        return nil, err
    }
    return commit, nil
}

// Uses the function database.SelectVersionChangeEvents to get all the version change events for a given repo. It builds a JSON that is a list of the same type as the return value of database.SelectVersionChangeEvents.
func GetRepoHistory(w http.ResponseWriter, r *http.Request) {
    _, ok := utils.Applications[r.URL.Query().Get("repo")]
    if !ok {
        WriteError(&w, "Repository not found", http.StatusNotFound)
        WriteResponseOk(&w, []*VersionChangeEventWithCommit{})
        return
    }
    versionChangeEvents, err := database.SelectVersionChangeEvents(r.URL.Query().Get("repo"))
    if err != nil {
        log.Println(err.Error())
        WriteResponseOk(&w, []*VersionChangeEventWithCommit{})
        return
    }
    i, j := NormalizeSliceIndexes(len(versionChangeEvents), r)
    var response = []*VersionChangeEventWithCommit{}
    for k := i; k < j; k++ {
        versionChangeEvent := versionChangeEvents[k]
        commit, err := getCommit(r.URL.Query().Get("repo"), versionChangeEvent.Hash)
        if err != nil {
            log.Println(err.Error())
            // We continue instead of returning mainly because we don't want to stop the loop if there is an error getting a commit. We just want to ignore that commit.
            continue
        }
        response = append(response, &VersionChangeEventWithCommit{
            Hash: versionChangeEvent.Hash,
            CreatedAt: versionChangeEvent.CreatedAt,
            Commit: commit,
        })
    }
    WriteResponseOk(&w, response)
}
