package main

import (
	"fmt"

	"github.com/DiegoAraujoJS/webdev-git-server/pkg/utils"
)

const A = 0

func init() {
	utils.Connect()
}

func main() {
    fmt.Println(utils.GetMasterBranchHash())
}
