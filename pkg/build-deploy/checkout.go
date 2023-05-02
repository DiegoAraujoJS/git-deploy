package builddeploy

import (
	"bytes"
	"log"

	"github.com/DiegoAraujoJS/webdev-git-server/pkg/utils"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func Checkout(repository string, hash string, stdout *bytes.Buffer) (*plumbing.Reference, error) {
	repo := utils.GetRepository(repository)
	w, err := repo.Worktree()

	if err != nil {
		log.Println("first error. ", err.Error())
		return nil, err
	}

    stdout.WriteString("Checkout commit " + hash + "\n")
	err = w.Checkout(&git.CheckoutOptions{
		Hash: plumbing.NewHash(hash),
        Force: true,
	})
	if err != nil {
		return nil, err
	}

	ref, err := repo.Head()
    stdout.WriteString("Successfully changed repository head to " + ref.Hash().String() + "\n")
	if err != nil {
		log.Println("second error.", err.Error())
		return nil, err
	}

	return ref, nil
}
