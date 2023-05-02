package utils

import (
	"fmt"
	"io/ioutil"
	"os"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

func fetchFromRemote(repo *git.Repository) error {
	remote, err := repo.Remote("origin")
	if err != nil {
		return err
	}

    var public_key *ssh.PublicKeys
    ssh_path := os.Getenv("HOME") + "/.ssh/id_ed25519"
    ssh_key, _ := ioutil.ReadFile(ssh_path)
    if err != nil {
        fmt.Println("Error reading private key from " + ssh_path, err)
        return err
    }
    public_key, err = ssh.NewPublicKeys("git", []byte(ssh_key), "")
    if err != nil {
        fmt.Println("Error creating public key", err)
        return err
    }

	err = remote.Fetch(&git.FetchOptions{
		RefSpecs:   []config.RefSpec{"refs/heads/*:refs/heads/*"},
		Force:      true,
        Auth:       public_key,
	})
    return err
}

// This function fetches origin with a Force flag set to true, which causes all local branches to be updated to match their remote counterparts. The function then iterates over the remote branches and force-updates the local branches accordingly.
func ForceUpdateAllBranches(repo *git.Repository) error {
	// Fetch the remote
    err := fetchFromRemote(repo)

	if err != nil {
		return err
	}

	// Get the remote branches
	remoteBranches, err := repo.Branches()
	if err != nil {
		return err
	}

	// Iterate over the remote branches
	err = remoteBranches.ForEach(func(ref *plumbing.Reference) error {
		// Resolve the commit for the remote branch
		commit, err := repo.CommitObject(ref.Hash())
		if err != nil {
			return err
		}

		// Force update the local branch to match the remote branch
		err = repo.Storer.SetReference(plumbing.NewHashReference(ref.Name(), commit.Hash))

        return err
	})
	return err
}

func GetBranch(repo *git.Repository, branch string) (*plumbing.Reference, error) {
    branches, _ := repo.Branches()
    for {
        b, _ := branches.Next()
        if b == nil {return b, nil}
        if b.Name().Short() == branch {return b, nil}
    }
}

func UpdateBranch(repo *git.Repository, branch *plumbing.Reference) (bool, error) {
    fmt.Println("update branch")
    last_commmit_hash := branch.Hash()
    fmt.Println("last commit hash", last_commmit_hash)
    err := fetchFromRemote(repo)
    if err != nil {return false, err}

    new_commit_hash := branch.Hash()
    fmt.Println("new commit hash", new_commit_hash)

    if last_commmit_hash != new_commit_hash {
        repo.Storer.SetReference(plumbing.NewHashReference(branch.Name(), new_commit_hash))
        return true, nil
    }
    return false, nil
}
