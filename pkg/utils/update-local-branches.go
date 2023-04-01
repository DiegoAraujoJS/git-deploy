package utils

import (
	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
)
// This function fetches the specified remote with a Force flag set to true, which causes all local branches to be updated to match their remote counterparts. The function then iterates over the remote branches and updates the local branches accordingly.
func forceUpdateAllBranches(repo *git.Repository, remoteName string) error {
	// Fetch the remote
	remote, err := repo.Remote(remoteName)
	if err != nil {
		return err
	}

	err = remote.Fetch(&git.FetchOptions{
		RefSpecs: []config.RefSpec{"refs/heads/*:refs/heads/*"},
		Force:    true,
	})
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
		if err != nil {
			return err
		}

		return nil
	})

	return err
}
