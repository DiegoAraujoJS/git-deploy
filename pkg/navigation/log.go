package navigation

import (
	"log"

	"github.com/DiegoAraujoJS/webdev-git-server/pkg/utils"
	"github.com/go-git/go-git/v5/plumbing/object"
)

const prefix = "refs/heads/"

var All_commits map[string][]*Branch = map[string][]*Branch{}

type Branch struct {
	Commit       *object.Commit `json:"commit"`
	NewReference string  `json:"new_reference"`
}

type BranchResponse struct {
	Commits           []*Branch `json:"commits"`
	CurrentVersion string    `json:"current_version"`
    Head           *object.Commit    `json:"head"`
}

func GetAllCommits(repository string) *BranchResponse {
    repo := utils.GetRepository(repository)

    if _, ok := All_commits[repository]; !ok {

        commits, err := repo.CommitObjects()
        if err != nil {
            log.Fatal(err.Error())
        }

        All_commits[repository] = []*Branch{}

        for {
            commit, _ := commits.Next()
            if commit == nil { break }
            All_commits[repository] = append(All_commits[repository], &Branch{
                Commit: commit,
            })
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
