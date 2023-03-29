package navigation

import (
	"log"
	"strconv"
	"strings"

	"github.com/DiegoAraujoJS/webdev-git-server/pkg/utils"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

const prefix = "refs/heads/"


type Commit struct {
	Hash         string
	Author       object.Signature
	Committer    object.Signature
	PGPSignature string
	Message      string
	TreeHash     plumbing.Hash
	ParentHashes []plumbing.Hash
}

type Branch struct {
	Commit       *Commit `json:"commit"`
	NewReference string  `json:"new_reference"`
}

type BranchResponse struct {
	Commits           []*Branch `json:"commits"`
	CurrentVersion string    `json:"current_version"`
    Head           *Commit    `json:"head"`
}

func GetReleaseBranchesWithTheirVersioning(repository string) *BranchResponse {
	repo := utils.GetRepository(repository)

	var result []*Branch

	branches, err := repo.Branches()
	if err != nil {
		log.Fatal(err.Error())
	}
	var current_version string
    var head_commit *Commit
	for {
		branch, _ := branches.Next()
		if branch == nil { break }
		if !strings.Contains(branch.Name().String(), "RELEASE") { continue }
        commits_from_master := utils.GetCommitsFromBranchToMaster(repository, branch)
        for k, c := range commits_from_master {
            undercase_split := strings.Split(branch.Name().String(), "_")
            version_number_string := undercase_split[len(undercase_split)-1]
            version := version_number_string + "." + strconv.Itoa(len(commits_from_master)-k)

            commit := &Commit{
                Hash:      c.Hash.String(),
                Author:    c.Author,
                Committer: c.Committer,
                Message:   c.Message,
            }

            result = append(result, &Branch{
                Commit: commit,
                NewReference: version,
            })
            
            if head, err := repo.Head(); err == nil && head.Hash().String() == c.Hash.String() {
                current_version = version
                head_commit = commit
            }
        }

	}
	return &BranchResponse{
		Commits:           result,
		CurrentVersion: current_version,
        Head: head_commit,
	}
}
