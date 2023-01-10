package utils

import (
	"log"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func GetBranchHash(branch_name string) plumbing.Hash {
	repo := GetRepository()

	branches, err := repo.Branches()

	if err != nil {
		log.Fatal(err.Error())
	}

	for {
		branch, err := branches.Next()
		if err != nil {
			break
		}
		if strings.Contains(branch.Name().String(), branch_name) {
			return branch.Hash()
		}
	}
	panic(branch_name + " branch not found")
}

func GetCommitsFromBranchToMaster(b *plumbing.Reference) []*object.Commit {
	repo := GetRepository()

	var commits []*object.Commit
	commit, _ := repo.CommitObject(b.Hash())
	master_commit, _ := repo.CommitObject(GetBranchHash("master"))
	merge_base, err := commit.MergeBase(master_commit)

	if err != nil {
		log.Fatal(err.Error())
	}

	log, _ := repo.Log(&git.LogOptions{From: commit.Hash})

	for {
		next_commit, err := log.Next()
		if err != nil {
			break
		}
		if next_commit.Hash.String() == merge_base[0].Hash.String() {
			break
		}
		commits = append(commits, next_commit)
	}
	return commits
}
