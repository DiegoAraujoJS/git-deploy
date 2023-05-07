package utils

import (
	"fmt"
	"io/ioutil"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

var public_key *ssh.PublicKeys

// This function fetches origin with a Force flag set to true, which causes all local branches to be updated to match their remote counterparts. The function then iterates over the remote branches and force-updates the local branches accordingly.
func ForceUpdateAllBranches(repo *git.Repository) error {
	// Fetch the remote
	remote, err := repo.Remote("origin")

	if err != nil {
        fmt.Println("Error getting remote", err)
		return err
	}

    if public_key == nil {
        ssh_key, _ := ioutil.ReadFile(ConfigValue.Credentials.Ssh)
        if err != nil {
            fmt.Println("Error reading private key from " + ConfigValue.Credentials.Ssh, err)
            return err
        }
        public_key, err = ssh.NewPublicKeys("git", []byte(ssh_key), "")
        if err != nil {
            fmt.Println("Error creating public key", err)
            return err
        }
    }

	err = remote.Fetch(&git.FetchOptions{
		RefSpecs:   []config.RefSpec{"refs/heads/*:refs/heads/*"},
		Force:      true,
        Auth:       public_key,
	})

	if err != nil {
        fmt.Println("Error fetching remote", err)
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
            fmt.Println("Error getting commit object", err)
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
        if b == nil {
            return nil, fmt.Errorf("branch not found: %s", branch)
        }
        if b.Name().Short() == branch {
            return b, nil
        }
    }
}
