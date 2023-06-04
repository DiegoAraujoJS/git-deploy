package routes

import (
	"encoding/json"
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

	response, err := json.Marshal(navigation.GetRepoTags(r.URL.Query().Get("repo")))
	if err != nil {
		WriteError(&w, "Error while getting release versions", 403)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func GetCommits(w http.ResponseWriter, r *http.Request) {
	_, ok := utils.Repositories[r.URL.Query().Get("repo")]

	if !ok {
		WriteError(&w, "Repository not found", http.StatusNotFound)
		return
	}

	commits := navigation.GetAllCommits(r.URL.Query().Get("repo"))
	// Filter by branch if branch is not empty
	branch := r.URL.Query().Get("branch")
	if branch != "" {
		filtered_commits := make([]*navigation.Commit, 0, len(commits))
		for _, commit := range commits {
			if funk.Contains(commit.Branch, branch) {
				filtered_commits = append(filtered_commits, commit)
			}
		}
		commits = filtered_commits
	}
    i, j := NormalizeSliceIndexes(len(commits), r)

	response, err := json.Marshal(commits[i:j])
	if err != nil {
		WriteError(&w, "Error while getting release versions", 403)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
