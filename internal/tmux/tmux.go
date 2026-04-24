package tmux

import (
	"os/exec"
	"strings"
)

func Health() (bool, string) {
	cmd := exec.Command("tmux", "-V")

	out, err := cmd.Output()
	if err != nil {
		return false, ""
	}

	return true, strings.Trim(strings.Split(string(out), " ")[1], "\n")
}
