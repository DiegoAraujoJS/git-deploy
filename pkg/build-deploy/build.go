package builddeploy

import (
	"log"
	"os"
	"os/exec"
)

// The function Build takes the name of the repo as a parameter. It executes the python script located at the ./scripts folder, that has to be named as the repo with the "py" extension.
//
// Example: for repo named "test", it executes (if exists) "./scripts/test.py".
func Build(action *Action) error {
    script := "./scripts/" + action.Repo + ".py"
    if _, err := os.Stat(script); os.IsNotExist(err) {
        error := "No build script " + script + " found for repo " + action.Repo + "\n" + err.Error()
        log.Println(error)
        action.Status.Stderr.WriteString(error)
        return err
	}
	cmd := exec.Command("python", script)
    cmd.Stdout = action.Status.Stdout
    cmd.Stderr = action.Status.Stderr
    action.Status.Stdout.WriteString("Executing build script " + script + " for repo " + action.Repo + "\n")
    err := cmd.Run()
    action.Status.Stdout.WriteString("Finished executing build script " + script + " for repo " + action.Repo + "\n")
    if err != nil {
        error := "Error while executing build script " + script + " for repo " + action.Repo + "\n" + err.Error()
        log.Println(error)
        action.Status.Stderr.WriteString(error)
        return err
    }
	return nil
}
