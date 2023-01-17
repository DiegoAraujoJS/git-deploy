package builddeploy

import (
	"fmt"
	"os/exec"

	"github.com/DiegoAraujoJS/webdev-git-server/pkg/utils"
)

func DeployIIS() error {
    cmd := exec.Command("cp", "-a", utils.ConfigValue.BuildOutputFolder + "/.", utils.ConfigValue.IISDirectory)
    cmd.Dir = utils.ConfigValue.IISDirectory

    stdout, err := cmd.Output()
    if err != nil {
        fmt.Println(err.Error())
        return err
    }
    fmt.Println("output", string (stdout))
    return err
}
