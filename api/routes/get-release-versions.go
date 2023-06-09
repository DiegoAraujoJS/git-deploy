package routes

import (
	"log"
	"net/http"
	"strconv"

	"github.com/DiegoAraujoJS/webdev-git-server/globals"
	"github.com/DiegoAraujoJS/webdev-git-server/pkg/utils"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type RepoTags struct {
    CurrentVersion  string          `json:"current_version"`
    Head            *object.Commit  `json:"head"`
    Branches        []string        `json:"branches"`
}

func GetRepoTags(w http.ResponseWriter, r *http.Request) {
	repo, ok := utils.Repositories[r.URL.Query().Get("repo")]

	if !ok {
		WriteError(&w, "Repository not found", http.StatusNotFound)
		return
	}

    branches, err := repo.Branches()
    if err != nil {
        log.Println(err.Error())
        return
    }

    repo_tags := &RepoTags{}
    repo_branches := []string{}
    for {
        branch, err := branches.Next()
        if branch == nil || err != nil { break }
        repo_branches = append(repo_branches, branch.Name().Short())
    }

    repo_tags.Branches = utils.MergeSort(repo_branches, func(n string, m string) bool {
        return n < m
    })

    head, err := repo.Head()
    if err != nil {
        log.Println(err.Error())
        return
    }
    head_commit, err := repo.CommitObject(head.Hash())
    if err != nil {
        log.Println(err.Error())
        return
    }
    repo_tags.Head = head_commit

    WriteResponseOk(&w, repo_tags)
}

func GetCommits(w http.ResponseWriter, r *http.Request) {
	repo, ok := utils.Repositories[r.URL.Query().Get("repo")]
    globals.Get_commits_rw_mutex.RLock()
    defer globals.Get_commits_rw_mutex.RUnlock()

	if !ok {
		WriteError(&w, "Repository not found", http.StatusNotFound)
		return
	}

	branch := r.URL.Query().Get("branch")

    log_options := &git.LogOptions{
        All: true,
        Order: git.LogOrderCommitterTime,
    }
    if branch != "" {
        ref, err := utils.GetBranch(repo, branch)
        if err != nil {
            WriteError(&w, "Branch not found", http.StatusNotFound)
            return
        }
        log_options.From = ref.Hash()
        log_options.All = false
    }

    log, err := repo.Log(log_options)
    if err != nil {
        WriteError(&w, "Error getting commits", http.StatusInternalServerError)
        return
    }

    i, i_err := strconv.Atoi(r.URL.Query().Get("i"))
    j, j_err := strconv.Atoi(r.URL.Query().Get("j"))
    if i_err != nil {i = 0}
    if j < i {j = i}

    var filtered_commits []*object.Commit
    if j_err == nil {
        filtered_commits = make([]*object.Commit, 0, j - i)
    } else {
        filtered_commits = make([]*object.Commit, 0)
    }

    var continue_loop = func (count int, j int) bool {
        if j_err != nil {
            return true
        }
        return count < j
    }

    c, err := log.Next()
    count := 0
    for err == nil && log != nil && continue_loop(count, j) {
        if count >= i {
            filtered_commits = append(filtered_commits, c)
        }
        count++
        c, err = log.Next()
    }

    WriteResponseOk(&w, filtered_commits)
}
