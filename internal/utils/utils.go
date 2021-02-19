package utils

import (
	"bufio"
	"os/exec"
)

func GitRoot() (string, error) {
	cEx := exec.Command("git", "rev-parse", "--show-toplevel")
	stdout, _ := cEx.StdoutPipe()
	cEx.Stderr = cEx.Stdout
	err := cEx.Start()

	var m string
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		m += scanner.Text()
	}

	cEx.Wait()
	return m, err
}
