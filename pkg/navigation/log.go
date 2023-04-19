package navigation

import (
	"log"
	"strconv"
	"strings"

	"github.com/DiegoAraujoJS/webdev-git-server/pkg/utils"
	"github.com/go-git/go-git/v5/plumbing/object"
)

const prefix = "refs/heads/"

type Branch struct {
	Commit       *object.Commit `json:"commit"`
	NewReference string  `json:"new_reference"`
}

type BranchResponse struct {
	Commits           []*Branch `json:"commits"`
	CurrentVersion string    `json:"current_version"`
    Head           *object.Commit    `json:"head"`
}

func GetReleaseBranchesWithTheirVersioning(repository string) *BranchResponse {
	repo := utils.GetRepository(repository)

	var result []*Branch

	branches, err := repo.Branches()
	if err != nil {
		log.Fatal(err.Error())
	}
	var current_version string
    var head_commit *object.Commit
	for {
		branch, _ := branches.Next()
		if branch == nil { break }
		if !strings.Contains(branch.Name().String(), "RELEASE") { continue }
        commits_from_master := utils.GetCommitsFromBranchToMaster(repository, branch)
        for k, c := range commits_from_master {
            undercase_split := strings.Split(branch.Name().String(), "_")
            version_number_string := undercase_split[len(undercase_split)-1]
            version := version_number_string + "." + strconv.Itoa(len(commits_from_master)-k)

            result = append(result, &Branch{
                Commit: c,
                NewReference: version,
            })
            
            if head, err := repo.Head(); err == nil && head.Hash().String() == c.Hash.String() {
                current_version = version
                head_commit = c
            }
        }

	}
	return &BranchResponse{
		Commits:           result,
		CurrentVersion: current_version,
        Head: head_commit,
	}
}

func GetAllCommits(repository string) *BranchResponse {
    repo := utils.GetRepository(repository)

    var result []*Branch

    commits, err := repo.CommitObjects()

    if err != nil {
        log.Fatal(err.Error())
    }
    for {
        commit, _ := commits.Next()
        if commit == nil { break }
        result = append(result, &Branch{
            Commit: commit,
        })
    }
    head, _ := repo.Head()
    head_commit, _ := repo.CommitObject(head.Hash())

    // Sort commits by date. The most recent is the first.
    result = utils.BubbleSort(result, func(n *Branch, m *Branch) bool {
        return n.Commit.Committer.When.Before(m.Commit.Author.When)
    })

    return &BranchResponse{
        Commits: result,
        Head: head_commit,
    }
}
