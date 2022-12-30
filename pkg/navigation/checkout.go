package navigation

import (
	"fmt"
	"log"

	"github.com/DiegoAraujoJS/webdev-git-server/pkg/utils"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func StringToHash(hash string) plumbing.Hash {
	id := plumbing.NewHash(hash)
	return id
}

func Checkout(hash plumbing.Hash) {
	repo := utils.GetRepository()

	w, err := repo.Worktree()

	if err != nil {
		log.Fatal(err.Error())
	}

	err = w.Checkout(&git.CheckoutOptions{
		Hash: hash,
	})

	if err != nil {
		log.Fatal(err.Error())
	}

	ref, err := repo.Head()

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(ref.Hash())
}
