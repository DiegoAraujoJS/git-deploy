package main

import (
	// "github.com/DiegoAraujoJS/webdev-git-server/api"
	builddeploy "github.com/DiegoAraujoJS/webdev-git-server/pkg/build-deploy"
	"github.com/DiegoAraujoJS/webdev-git-server/pkg/utils"
)

const A = 0

func init() {
	utils.Connect()
}

func main() {
	builddeploy.DeployIIS()
	// api.ListenAndServe()
}
