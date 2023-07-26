package builddeploy

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/DiegoAraujoJS/webdev-git-server/pkg/utils"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

// This function will call "git checkout <hash>" on the repository. It will not use the Checkout function from the go-git library as it has big performance problems in windows. It uses the git CLI tool.
func checkoutFromCli (action *Action) error {
    if utils.ConfigValue.CliBinaryForCheckout == "" { return fmt.Errorf("no cli binary for checkout") }
    // Save the current path
    original_wd_path, err := os.Getwd()
    if err != nil {
        return err
    }
    // Find the repository path by iterating over utils.ConfigValue.Directories
    var repoPath string
    for _, v := range utils.ConfigValue.Directories {
        if v.Name == action.Repo {
            repoPath = v.Directory
        }
    }
    if repoPath == "" {
        return fmt.Errorf("repository not found")
    }
    // Change dir to repository path and run the command
    if err = os.Chdir(repoPath); err != nil {
        return err
    }
    cmd := exec.Command(utils.ConfigValue.CliBinaryForCheckout, "checkout", "--force", "--quiet", action.Hash.String())
    cmd.Stdout = action.Status.Stdout
    cmd.Stderr = action.Status.Stderr
    if err = cmd.Run(); err != nil {
        return err
    }
    // Change again to the original directory, so not to make any change to the working directory state.
    return os.Chdir(original_wd_path)
}

func Checkout(action *Action) (*plumbing.Reference, error) {
	repo := utils.Repositories[action.Repo]

    action.Status.Stdout.WriteString("Checkout commit " + action.Hash.String() + "\n")
    var err error
    if utils.ConfigValue.CliBinaryForCheckout != "" {
        err = checkoutFromCli(action)
    } else {
        w, w_err := repo.Worktree()
        if w_err != nil {
            log.Println("Error while getting repository worktree -> "+err.Error())
            action.Status.Stderr.WriteString(err.Error())
            return nil, err
        }
        err = w.Checkout(&git.CheckoutOptions{
            Hash: action.Hash,
            Force: true,
        })
    }
	if err != nil {
        log.Println("Error while checking out commit -> "+err.Error())
        action.Status.Stderr.WriteString(err.Error())
		return nil, err
	}

	ref, err := repo.Head()
	if err != nil {
        log.Println("Error while getting repository head -> "+err.Error())
        action.Status.Stderr.WriteString(err.Error())
		return nil, err
	}
    action.Status.Stdout.WriteString("Successfully changed repository head to " + ref.Hash().String() + "\n")

	return ref, nil
}
