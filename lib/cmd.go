package lib

import "os/exec"

func ExecuteCmd(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	output, err := cmd.CombinedOutput()

	return string(output), err
}
