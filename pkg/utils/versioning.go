package utils

import (
	"fmt"
	"log"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func GetParents() map[string][]plumbing.Hash {
	repo := GetRepository()

	commits, err := repo.Log(&git.LogOptions{})

	if err != nil {
		log.Fatal(err.Error())
	}

	childs := map[string][]plumbing.Hash{}

	for {
		commit, err := commits.Next()

		if err != nil {
			break
		}

		fmt.Println(commit.String())
		for _, parent := range commit.ParentHashes {

			_, ok := childs[parent.String()]
			if !ok {
				childs[parent.String()] = []plumbing.Hash{commit.Hash}
			} else {
				childs[parent.String()] = append(childs[parent.String()], commit.Hash)
			}

		}

	}

	commits.Close()

	return childs
}

func GetMasterBranchHash() plumbing.Hash {
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
		if strings.Contains(branch.Name().String(), "master") {
			return branch.Hash()
		}
	}
	panic("master branch not found")
}
