package navigation

import (
	"fmt"
	"log"

	"github.com/DiegoAraujoJS/webdev-git-server/pkg/utils"
	"github.com/go-git/go-git/v5"
)

func Pull() {
	repo := utils.GetRepository()

    w, err := repo.Worktree()

    if err != nil {
        log.Fatal(err.Error())
    }

    err = w.Pull(&git.PullOptions{
		RemoteName: "origin",
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
