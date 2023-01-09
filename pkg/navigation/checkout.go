package navigation

import (
	"fmt"
	"log"

	"github.com/DiegoAraujoJS/webdev-git-server/pkg/utils"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func TagNameToHash(tag string) plumbing.Hash {
	repo := utils.GetRepository()
    var tag_object plumbing.Hash
    for _, v := range GetTags() {
        fmt.Println(v.Name)
        if v.Name == tag {
            fmt.Println("match", v.Name)
            obj, err := repo.TagObject(v.Hash)
            if err != nil {
                log.Fatal(err.Error())
            }
            tag_object = obj.Target
            break
        }
    }
    return tag_object
}

func StringToHash(hash string) plumbing.Hash {
	id := plumbing.NewHash(hash)
	return id
}

func Checkout(hash plumbing.Hash) plumbing.Hash {
	repo := utils.GetRepository()

	w, err := repo.Worktree()

	if err != nil {
		log.Fatal(err.Error())
	}

    fmt.Println(hash)

	err = w.Checkout(&git.CheckoutOptions{
		Hash: hash,
	})

	if err != nil {
		log.Fatal("could not checkout", err.Error())
	}

	ref, err := repo.Head()

	if err != nil {
		log.Fatal(err.Error())
	}

    return ref.Hash()
}
