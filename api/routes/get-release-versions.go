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
    slice_len := j - i
	// Filter by branch if branch is not empty
	branch := r.URL.Query().Get("branch")
    filtered_commits := make([]*navigation.Commit, 0, slice_len)
    count := 0
    for _, commit := range commits {
        if count == slice_len { break }
        if branch == "" || funk.Contains(commit.Branch, branch) {
            if count >= i {filtered_commits = append(filtered_commits, commit)}
            count++
        }
    }

    WriteResponseOk(&w, filtered_commits)
}
