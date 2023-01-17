package builddeploy

import (
	"fmt"
	"os/exec"

	"github.com/DiegoAraujoJS/webdev-git-server/pkg/utils"
)

func Build() error {
    cmd := exec.Command("npx", "vite", "build")
    cmd.Dir = utils.ConfigValue.ClientDirectory

    stdout, err := cmd.Output()
    if err != nil {
        fmt.Println(err.Error())
        return err
    }
    fmt.Println(string (stdout))

    return err
}
