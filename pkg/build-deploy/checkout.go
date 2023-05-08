package builddeploy

import (
	"log"

	"github.com/DiegoAraujoJS/webdev-git-server/pkg/utils"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func Checkout(action *Action) (*plumbing.Reference, error) {
	repo := utils.Repositories[action.Repo]
	w, err := repo.Worktree()

	if err != nil {
        log.Println("Error while getting repository worktree -> "+err.Error())
        action.Status.Stderr.WriteString(err.Error())
		return nil, err
	}

    action.Status.Stdout.WriteString("Checkout commit " + action.Hash + "\n")
	err = w.Checkout(&git.CheckoutOptions{
		Hash: plumbing.NewHash(action.Hash),
        Force: true,
	})
	if err != nil {
        log.Println("Error while checking out commit -> "+err.Error())
        action.Status.Stderr.WriteString(err.Error())
		return nil, err
	}

	ref, err := repo.Head()
    action.Status.Stdout.WriteString("Successfully changed repository head to " + ref.Hash().String() + "\n")
	if err != nil {
        log.Println("Error while getting repository head -> "+err.Error())
        action.Status.Stderr.WriteString(err.Error())
		return nil, err
	}

	return ref, nil
}
