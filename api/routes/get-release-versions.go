package routes

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

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

    repo_tags.Branches = utils.MergeSort(repo_branches, sortByName)

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
    start := time.Now()
    repo, ok := utils.Repositories[r.URL.Query().Get("repo")]

    if !ok {
        WriteError(&w, "Repository not found", http.StatusNotFound)
        return
    }

    branch := r.URL.Query().Get("branch")

    globals.Get_commits_rw_mutex.RLock()
    defer globals.Get_commits_rw_mutex.RUnlock()

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

    //  i       j       branch      solution

    //  YES     YES     YES         Get all commits from branch up to j. Filter from i.
    //  YES     YES     NO          Get all commits up to j with map. Filter from i.
    //  YES     NO      YES         Get all commits from branch. Filter from i.
    //  YES     NO      NO          Get all commits. Sort. Filter from i.
    //  NO      YES     YES         Get all commits from branch up to j.
    //  NO      YES     NO          Get all commits up to j with map.
    //  NO      NO      YES         Get all commits from branch.
    //  NO      NO      NO          Get all commits. Sort.

    if branch != "" {
        var log_options = &git.LogOptions{
            Order: git.LogOrderCommitterTime,
        }

        ref, err := utils.GetBranch(repo, branch)
        log_options.From = ref.Hash()
        if err != nil {
            WriteError(&w, err.Error(), http.StatusNotFound)
            return
        }

        log, _ := repo.Log(log_options)
        c, err := log.Next()

        var continue_loop = func (count int, j int) bool {
            if j_err != nil {
                return true
            }
            return count < j
        }

        for count := 0; err == nil && continue_loop(count, j); count++ {
            if count >= i {
                filtered_commits = append(filtered_commits, c)
            }
            c, err = log.Next()
        }
    } else {
        opts := &sortedCommitsOptions{
            all: j_err != nil,
            number: j,
        }

        if i_err == nil {
            filtered_commits = getSortedCommitsMap(repo, opts)[i:]
        } else {
            filtered_commits = getSortedCommitsMap(repo, opts)
        }
    }

    WriteResponseOk(&w, filtered_commits)
    fmt.Println(time.Since(start))
}
