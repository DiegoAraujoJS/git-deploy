package navigation

import (
	"log"

	"github.com/DiegoAraujoJS/webdev-git-server/pkg/utils"
	"github.com/go-git/go-git/v5"
)

func Pull(repository string) {
	repo := utils.GetRepository(repository)

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

	_, err = repo.Head()

	if err != nil {
		log.Fatal(err.Error())
	}
}
