package navigation

import (
	"log"

	"github.com/DiegoAraujoJS/webdev-git-server/pkg/utils"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

var All_commits map[string][]*Commit = map[string][]*Commit{}

type Commit struct {
    *object.Commit
    NewReference    string          `json:"new_reference"`
    Branch          []string        `json:"branches"`
}

type RepoTags struct {
    CurrentVersion  string          `json:"current_version"`
    Head            *object.Commit  `json:"head"`
    Branches        []string        `json:"branches"`
}

func GetAllCommits(repository string) []*Commit {
    repo := utils.Repositories[repository]

    if _, ok := All_commits[repository]; !ok {
        branches, err := repo.Branches()
        if err != nil {
            log.Println("Error while fetching branches", err.Error())
        }

        var commits_map = map[string]*Commit{}

        for {
            branch, err := branches.Next()
            if branch == nil || err != nil { break }
            log, _ := repo.Log(&git.LogOptions{
                From: branch.Hash(),
            })

            for {
                commit, err := log.Next()
                if commit == nil || err != nil { break }
                if payload, ok := commits_map[commit.Hash.String()] ; ok {
                    payload.Branch = append(payload.Branch, branch.Name().Short())
                    continue
                }
                commits_map[commit.Hash.String()] = &Commit{
                    Commit: commit,
                    Branch: []string{branch.Name().Short()},
                }
            }
        }

        All_commits[repository] = []*Commit{}
        for _, v := range commits_map {
            All_commits[repository] = append(All_commits[repository], v)
        }
        // Sort commits by date. The most recent is the first.
        All_commits[repository] = utils.MergeSort(All_commits[repository], func(n *Commit, m *Commit) bool {
           return m.Committer.When.Before(n.Committer.When)
        })
        // Sort branches by name.
    }

    return All_commits[repository]
}

func GetRepoTags(repository string) *RepoTags {
    repo := utils.Repositories[repository]
    branches, err := repo.Branches()
    if err != nil {
        log.Println(err.Error())
    }

    repo_tags := &RepoTags{}
    repo_branches := []string{}
    for {
        branch, err := branches.Next()
        if branch == nil || err != nil { break }
        repo_branches = append(repo_branches, branch.Name().Short())
    }
    repo_tags.Branches = repo_branches

    repo_tags.Branches = utils.MergeSort(repo_tags.Branches, func(n string, m string) bool {
        return n < m
    })

    head, err := repo.Head()
    if err != nil {
        log.Println(err.Error())
    }
    head_commit, err := repo.CommitObject(head.Hash())
    if err != nil {
        log.Println(err.Error())
    }
    repo_tags.Head = head_commit

    return repo_tags
}
