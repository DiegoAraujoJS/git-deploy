package navigation

import (
	"log"

	"github.com/DiegoAraujoJS/webdev-git-server/pkg/utils"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func Checkout(repository string, hash string) (*plumbing.Reference, error) {
	repo := utils.GetRepository(repository)
	w, err := repo.Worktree()

	if err != nil {
		log.Println("first error. ", err.Error())
		return nil, err
	}

	err = w.Checkout(&git.CheckoutOptions{
		Hash: plumbing.NewHash(hash),
	})
	if err != nil {
		return nil, err
	}

	ref, err := repo.Head()
	if err != nil {
		log.Println("second error.", err.Error())
		return nil, err
	}

	return ref, nil
}
