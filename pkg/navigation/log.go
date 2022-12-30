package navigation

import (
	"log"

	"github.com/DiegoAraujoJS/webdev-git-server/pkg/utils"
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
