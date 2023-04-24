package navigation

import (
	"log"

	"github.com/DiegoAraujoJS/webdev-git-server/pkg/utils"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

const prefix = "refs/heads/"

var All_commits map[string][]*Branch = map[string][]*Branch{}

type Branch struct {
	Commit       *object.Commit `json:"commit"`
	NewReference string  `json:"new_reference"`
    Branch []string `json:"branches"`
}

type BranchResponse struct {
	Commits           []*Branch `json:"commits"`
	CurrentVersion string    `json:"current_version"`
    Head           *object.Commit    `json:"head"`
}

func GetAllCommits(repository string) *BranchResponse {
    repo := utils.GetRepository(repository)

    if _, ok := All_commits[repository]; !ok {

        branches, err := repo.Branches()
        if err != nil {
            log.Fatal(err.Error())
        }

        var commits_map = map[string]*Branch{}

        for {
            branch, _ := branches.Next()
            if branch == nil { break }
            log, _ := repo.Log(&git.LogOptions{
                From: branch.Hash(),
            })

            for {
                commit, _ := log.Next()
                if (commit == nil) {break}
                if commits_map[commit.Hash.String()] == nil {
                    commits_map[commit.Hash.String()] = &Branch{
                        Commit: commit,
                        Branch: []string{branch.Name().Short()},
                    }
                } else {
                    commits_map[commit.Hash.String()].Branch = append(commits_map[commit.Hash.String()].Branch, branch.Name().Short())
                }
            }
            
        }

        All_commits[repository] = []*Branch{}
        for _, v := range commits_map {
            All_commits[repository] = append(All_commits[repository], v)
        }
        // Sort commits by date. The most recent is the first.
        All_commits[repository] = utils.MergeSort(All_commits[repository], func(n *Branch, m *Branch) bool {
           return m.Commit.Committer.When.Before(n.Commit.Author.When)
        })
    }

    head, _ := repo.Head()
    head_commit, _ := repo.CommitObject(head.Hash())

    return &BranchResponse{
        Commits: All_commits[repository],
        Head: head_commit,
    }
}
