package navigation

import (
	"fmt"
	"log"

	"github.com/DiegoAraujoJS/webdev-git-server/pkg/utils"
	"github.com/go-git/go-git/v5"
)

func Pull(repository string) error {
	repo := utils.GetRepository(repository)

    // return early if repo has no remote
    if _, err := repo.Remote("origin"); err != nil {
        fmt.Println("No remote found for this repository")
        return nil
    }

    w, err := repo.Worktree()

    if err != nil {
        log.Println(err.Error())
        return err
    }

    err = w.Pull(&git.PullOptions{
		RemoteName: "origin",
	})

	if err != nil {
		log.Println(err.Error())
        return err
	}

	if err != nil {
		log.Println(err.Error())
        return err
	}
    fmt.Println("Successfully pulled from origin!")
    return nil
}
