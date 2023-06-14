package routes

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func sortByName(n string, m string) bool {
    return n < m
}

func committerWhenAfter(a *object.Commit, b *object.Commit) bool {
    return a.Committer.When.After(b.Committer.When)
}

type sortedCommitsOptions struct {
    all bool
    number int
}

// This function is thought over the fact that the most recent commit of a git repository is one of its leafs. The most recent before that is either one of the parents of the most recent one or one of the leafs of the repository, and so on.
func getSortedCommitsMap(repo *git.Repository, opts *sortedCommitsOptions) []*object.Commit {
    var sorted = []*object.Commit{}

    // 1. Define the set of leafs
    var set = map[plumbing.Hash]*object.Commit{}
    branches, _ := repo.Branches()
    for r, err := branches.Next(); err == nil; r, err = branches.Next() {
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
        if committerWhenAfter(commit, most_recent) {
            most_recent = commit
        }
    }

    // 3. Add the commit found in 2. to a list.
    sorted = append(sorted, most_recent)

    // 4. If the commit found in 2. has no parents, return the list of 3. Else redefine the set (remove the commit found in 2., add its parents), and go to step 2.
    if most_recent.NumParents() == 0 || (len(sorted) == opts.number && !opts.all) {
        return sorted
    } else {
        delete(set, most_recent.Hash)
        parents_iter := most_recent.Parents()
        for c, err := parents_iter.Next(); err == nil; c, err = parents_iter.Next() {
            set[c.Hash] = c
        }
        goto step_two
    }
}
