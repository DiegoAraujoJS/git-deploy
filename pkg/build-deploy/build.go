package builddeploy

import (
	"bytes"
	"os"
	"os/exec"
)

// The function Build takes the name of the repo as a parameter. It executes the python script located at the ./scripts folder, that has to be named as the repo with the "py" extension.
//
// Example: for repo named "test", it executes (if exists) "./scripts/test.py".
func Build(repo string, stdout *bytes.Buffer, stderr *bytes.Buffer) error {
    script := "./scripts/" + repo + ".py"
    if _, err := os.Stat(script); os.IsNotExist(err) {
        stderr.WriteString(err.Error())
        return err
	}
	cmd := exec.Command("python", script)
    cmd.Stdout = stdout
    cmd.Stderr = stderr
    err := cmd.Run()
    if err != nil {
        stderr.WriteString(err.Error())
        return err
    }
    stdout.WriteString("Successfully finished executing " + script)
	return nil
}
