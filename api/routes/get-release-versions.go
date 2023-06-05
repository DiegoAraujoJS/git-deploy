package routes

import (
	"net/http"

	"github.com/DiegoAraujoJS/webdev-git-server/pkg/navigation"
	"github.com/DiegoAraujoJS/webdev-git-server/pkg/utils"
	"github.com/thoas/go-funk"
)

func GetRepoTags(w http.ResponseWriter, r *http.Request) {
	_, ok := utils.Repositories[r.URL.Query().Get("repo")]

	if !ok {
		WriteError(&w, "Repository not found", http.StatusNotFound)
		return
	}

    WriteResponseOk(&w, navigation.GetRepoTags(r.URL.Query().Get("repo")))
}

func GetCommits(w http.ResponseWriter, r *http.Request) {
	_, ok := utils.Repositories[r.URL.Query().Get("repo")]

	if !ok {
		WriteError(&w, "Repository not found", http.StatusNotFound)
		return
	}

	commits := navigation.GetAllCommits(r.URL.Query().Get("repo"))
    i, j := NormalizeSliceIndexes(len(commits), r)
	// Filter by branch if branch is not empty
	branch := r.URL.Query().Get("branch")
    filtered_commits := make([]*navigation.Commit, 0, j - i)
    for ; i < j; i++ {
        commit := commits[i]
        if branch == "" || funk.Contains(commit.Branch, branch) {
            filtered_commits = append(filtered_commits, commit)
        }
    }

    WriteResponseOk(&w, filtered_commits)
}
