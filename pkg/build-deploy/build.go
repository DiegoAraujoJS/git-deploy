package builddeploy

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func Build(directory string, script string) error {
	fmt.Println(script)
	cmd := exec.Command("cp", script, directory)
	stdout, err := cmd.Output()
	if err != nil {
		log.Fatal(err.Error())
		return err
	}
	fmt.Println(string(stdout))
    if _, err = os.Stat(directory + "/build.sh") ; os.IsNotExist(err) {
        fmt.Println(err.Error())
    }
    fmt.Println(directory + "/build.sh")
	cmd = exec.Command("/bin/bash", directory + "/build.sh")

	stdout, err = cmd.Output()
	if err != nil {
		log.Fatal(err.Error())
		return err
	}
	fmt.Println(string(stdout))

	return err
}
