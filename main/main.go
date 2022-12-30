package main

import (
	"github.com/DiegoAraujoJS/webdev-git-server/pkg/navigation"
	"github.com/DiegoAraujoJS/webdev-git-server/pkg/utils"
)

const A = 0

func init() {
	utils.Connect()
}

func main() {
	// navigation.Checkout(navigation.StringToHash("2f31df2771abed6ffbabe281b90738ed54815f5d"))

    navigation.ShowTags()

	// logResult, err := repo.Log(&git.LogOptions{})
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }
	// for {
	// 	l, err := logResult.Next()
	// 	if err != nil {
	// 		break
	// 	}
	// 	fmt.Println(l.Hash)
	// }
}
