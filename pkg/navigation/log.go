package navigation

import (
	"log"
	"strings"

	"github.com/DiegoAraujoJS/webdev-git-server/pkg/utils"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func GetTags() []*object.Tag {
	repo := utils.GetRepository()

	tags, err := repo.Tags()

	if err != nil {
		log.Fatal(err.Error())
	}

	var list_tags []*object.Tag

	tags.ForEach(func(r *plumbing.Reference) error {
		obj, err := repo.TagObject(r.Hash())
		if err != nil {
			log.Fatal(err.Error())
			return err
		}
		list_tags = append(list_tags, obj)
		return err
	})

	return list_tags
}

func GetRemoteBranches() []*plumbing.Reference {
	repo := utils.GetRepository()

	remote, err := repo.Remote("origin")
	if err != nil {
		log.Fatal(err.Error())
	}

	ref_list, err := remote.List(&git.ListOptions{})

	if err != nil {
		log.Fatal(err.Error())
	}

	prefix := "refs/heads/"

	var remote_branches []*plumbing.Reference

	for _, ref := range ref_list {
		if !strings.HasPrefix(ref.Name().String(), prefix) {
			continue
		}
		remote_branches = append(remote_branches, ref)
	}
	return remote_branches
}
