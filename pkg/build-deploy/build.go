package builddeploy

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

// The function Build takes the name of the repo as a parameter. It executes the bash script located at the ./scripts folder, that has to be named as the repo with the "sh" extension. Example: for repo named "test", "test.sh".
func Build(repo string) error {
    script := "./scripts/" + repo + ".py"
    if _, err := os.Stat(script); os.IsNotExist(err) {
        payload := fmt.Errorf("file does not exist", fmt.Sprint(err))
        return payload
	}
	cmd := exec.Command("python", script)
    var out bytes.Buffer
    var stderr bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &stderr
    err := cmd.Run()
    if err != nil {
        return fmt.Errorf(fmt.Sprint(err) + ": " + stderr.String())
    }
	return err
}
