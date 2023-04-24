package navigation

import (
	"log"

	"github.com/DiegoAraujoJS/webdev-git-server/pkg/utils"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

var All_commits map[string]*BranchResponse = map[string]*BranchResponse{}

type Branch struct {
	Commit       *object.Commit `json:"commit"`
	NewReference string  `json:"new_reference"`
    Branch []string `json:"branches"`
}

type BranchResponse struct {
	Commits           []*Branch `json:"commits"`
	CurrentVersion string    `json:"current_version"`
    Head           *object.Commit    `json:"head"`
    Branches       []string    `json:"branches"`
}

func GetAllCommits(repository string) *BranchResponse {
    repo := utils.GetRepository(repository)

    if _, ok := All_commits[repository]; !ok {
        branches, err := repo.Branches()
        if err != nil {
            log.Fatal(err.Error())
        }

        var commits_map = map[string]*Branch{}
        var branches_names = map[string]struct{}{}

        for {
            branch, err := branches.Next()
            if branch == nil || err != nil { break }
            branches_names[branch.Name().Short()] = struct{}{}
            log, _ := repo.Log(&git.LogOptions{
                From: branch.Hash(),
            })

            for {
                commit, err := log.Next()
                if commit == nil || err != nil { break }
                if _, ok := commits_map[commit.Hash.String()] ; !ok {
                    commits_map[commit.Hash.String()] = &Branch{
                        Commit: commit,
                        Branch: []string{branch.Name().Short()},
                    }
                } else {
                    commits_map[commit.Hash.String()].Branch = append(commits_map[commit.Hash.String()].Branch, branch.Name().Short())
                }
            }
        }

        All_commits[repository] = &BranchResponse{}
        for _, v := range commits_map {
            All_commits[repository].Commits = append(All_commits[repository].Commits, v)
        }
        for k := range branches_names {
            All_commits[repository].Branches = append(All_commits[repository].Branches, k)
        }
        // Sort commits by date. The most recent is the first.
        All_commits[repository].Commits = utils.MergeSort(All_commits[repository].Commits, func(n *Branch, m *Branch) bool {
           return m.Commit.Committer.When.Before(n.Commit.Author.When)
        })
        // Sort branches by name.
        All_commits[repository].Branches = utils.MergeSort(All_commits[repository].Branches, func(n string, m string) bool {
            return n < m
        })

    }

    head, _ := repo.Head()
    head_commit, _ := repo.CommitObject(head.Hash())
    All_commits[repository].Head = head_commit
    return All_commits[repository]
}
