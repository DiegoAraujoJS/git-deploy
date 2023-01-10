package navigation

import (
	"errors"
	"log"

	"github.com/DiegoAraujoJS/webdev-git-server/pkg/utils"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func Checkout(tag string) (*plumbing.Reference, error) {
	repo := utils.GetRepository()
	w, err := repo.Worktree()

    versions := GetReleaseBranchesWithTheirVersioning()
    var commit *object.Commit
    for _, v := range versions {
        if v.NewReference == tag {
            commit = v.Commit
        }
    }
	if err != nil {
		log.Fatal(err.Error())
	}

	err = w.Checkout(&git.CheckoutOptions{
		Hash: commit.Hash,
	})
	if err != nil {
        return &plumbing.Reference{}, errors.New("could not checkout")
	}

	ref, err := repo.Head()
	if err != nil {
		log.Fatal(err.Error())
	}

    return ref, nil
}
