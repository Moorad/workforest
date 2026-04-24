package fzf

import (
	"os/exec"
	"strings"
)

func Health() (bool, string) {
	cmd := exec.Command("fzf", "--version")

	out, err := cmd.Output()
	if err != nil {
		return false, ""
	}

	return true, strings.Split(string(out), " ")[0]
}
