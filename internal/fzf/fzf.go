package fzf

import (
	"strings"

	"github.com/Moorad/workforest/internal/exec"
)

func Health() (bool, string) {
	out, err := exec.CmdOutput("fzf", "--version")
	if err != nil {
		return false, ""
	}

	return true, strings.Split(out, " ")[0]
}
