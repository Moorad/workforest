package git

import (
	"os/exec"
	"strings"
)

func Health() (bool, string) {
	cmd := exec.Command("git", "version")

	out, err := cmd.Output()
	if err != nil {
		return false, ""
	}

	return true, strings.Trim(strings.Split(string(out), " ")[2], "\n")
}
