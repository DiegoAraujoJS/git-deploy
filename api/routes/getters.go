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

func getSortedCommitsMap(repo *git.Repository, n int) []*object.Commit {
    // 1. Define the set of leafs
    var set = map[string]*object.Commit{}
    branches, _ := repo.Branches()
    branches.ForEach(func(r *plumbing.Reference) error {
        c, _ := repo.CommitObject(r.Hash())
        set[c.Hash.String()] = c
        return nil
    })

    var sorted = []*object.Commit{}
    for {
        // 2. Find the most recent commit by iterating over the set.
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

        // 3. Remove the element found in 2 from the set.
        delete(set, most_recent.Hash.String())

        // 4. Add the element found in 1 to a list.
        sorted = append(sorted, most_recent)

        // 5. If the element found in 2. has no parents, return the list of 4. Else add the parents of the element found in 2 to the set.
        if most_recent.NumParents() == 0 || len(sorted) == n {return sorted}
        most_recent.Parents().ForEach(func(c *object.Commit) error {
            set[c.Hash.String()] = c
            return nil
        })

        // 6. Go to 2.
    }
}
