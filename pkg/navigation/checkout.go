package navigation

import (
	"errors"
	"fmt"
	"log"

	"github.com/DiegoAraujoJS/webdev-git-server/pkg/utils"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func Checkout(hash string) (*plumbing.Reference, error) {
	repo := utils.GetRepository()
	w, err := repo.Worktree()

	if err != nil {
		log.Fatal(err.Error())
	}

	err = w.Checkout(&git.CheckoutOptions{
		Hash: plumbing.NewHash(hash),
	})
	if err != nil {
		return &plumbing.Reference{}, errors.New("could not checkout")
	}

	ref, err := repo.Head()
	if err != nil {
		log.Fatal(err.Error())
	}

    fmt.Println("head", ref.Hash())

	return ref, nil
}
