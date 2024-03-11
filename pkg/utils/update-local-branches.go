package utils

import (
	"fmt"
	"io/ioutil"

	"github.com/DiegoAraujoJS/webdev-git-server/globals"
	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

var public_key *ssh.PublicKeys

func getSshKey () *ssh.PublicKeys {
    if public_key != nil {
        return public_key
    }
    ssh_key, err := ioutil.ReadFile(ConfigValue.Credentials.Ssh)
    if err != nil {
        fmt.Println("Error reading private key from " + ConfigValue.Credentials.Ssh, err)
        return nil
    }
    public_key, err = ssh.NewPublicKeys("git", ssh_key, "")
    if err != nil {
        fmt.Println("Error creating public key", err)
        return nil
    }
    return public_key
}


func pruneLocalBranches(repo *git.Repository) error {
    remote, err := repo.Remote("origin")
    if err != nil {
        fmt.Println(err)
        return err
    }
    refs, err := remote.List(&git.ListOptions{Auth: getSshKey()})
    if err != nil {
        fmt.Println(err)
        return err
    }
    globals.Get_commits_rw_mutex.Lock()
    defer globals.Get_commits_rw_mutex.Unlock()
    local_branches, err := repo.Branches()
    if err != nil {
        fmt.Println(err)
        return err
    }
    local_branches_loop:
    for branch, branch_err := local_branches.Next(); branch_err == nil; branch, branch_err = local_branches.Next() {
        for _, ref := range refs {
            if ref.Name().IsBranch() && ref.Name().Short() == branch.Name().Short() {
                continue local_branches_loop
            }
        }
        fmt.Println("Deleting local branch that no longer exists on remote", branch.Name().String())
        repo.Storer.RemoveReference(branch.Name())
    }
    return nil
}

// This function fetches origin with a Force flag set to true, which causes all local branches to be updated to match their remote counterparts. The function then iterates over the remote branches and force-updates the local branches accordingly.
func ForceUpdateAllBranches(repo *git.Repository) error {
    sshKey := getSshKey()
    if sshKey == nil {return nil}
	// Fetch the remote
	remote, err := repo.Remote("origin")

	if err != nil {
        fmt.Println("Error getting remote", err)
		return err
	}

	err = remote.Fetch(&git.FetchOptions{
        RemoteName: "origin",
		RefSpecs:   []config.RefSpec{"+refs/heads/*:refs/heads/*"},
		Force:      true,
        Auth:       sshKey,
	})

	if err != nil && err != git.NoErrAlreadyUpToDate {
        fmt.Println("Error fetching remote", err)
		return err
	}

    go func (repo *git.Repository) {
        defer globals.GenericRecover()
        pruneLocalBranches(repo)
    }(repo)

	return err
}

func GetBranch(repo *git.Repository, branch string) (*plumbing.Reference, error) {
    refs, err := repo.References()

    for ref, err := refs.Next(); err == nil; ref, err = refs.Next() {
        if ref.Name().Short() == branch {
            return ref, nil
        }
    }
    return nil, err
}

