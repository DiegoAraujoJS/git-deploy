package utils

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type SortedCommitsOptions struct {
    All bool
    Number int
}

// This function is thought over the fact that the most recent commit of a git repository is one of its leafs. The most recent before that is either one of the parents of the most recent one or one of the leafs of the repository, and so on.
func GetSortedCommitsMap(repo *git.Repository, opts *SortedCommitsOptions) []*object.Commit {
    var sorted = []*object.Commit{}

    // 1. Define the set of leafs
    var set = map[plumbing.Hash]*object.Commit{}
    refs, _ := repo.References()
    for r, err := refs.Next(); err == nil; r, err = refs.Next() {
        if !r.Name().IsBranch() && !r.Name().IsRemote() {continue}

        c, c_err := repo.CommitObject(r.Hash())
        if c_err != nil {continue}
        set[c.Hash] = c
    }

    // 2. Find the most recent commit by iterating over the set.
    step_two:
    var most_recent *object.Commit
    for _, commit := range set {
        if most_recent == nil {
            most_recent = commit
            continue
        }
        if commit.Committer.When.After(most_recent.Committer.When) {
            most_recent = commit
        }
    }

    // 3. Add the commit found in 2. to a list.
    sorted = append(sorted, most_recent)

    // Special return case: list is full of commits according to opts.Number cap.
    if len(sorted) == opts.Number && !opts.All {return sorted}

    // 4. If the commit found in 2. has no parents, return the list of 3. Else redefine the set (remove the commit found in 2., add its parents), and go to step 2.
    if most_recent.NumParents() == 0 {return sorted}

    delete(set, most_recent.Hash)
    parents_iter := most_recent.Parents()
    for parent, err := parents_iter.Next(); err == nil; parent, err = parents_iter.Next() {
        set[parent.Hash] = parent
    }
    goto step_two
}
