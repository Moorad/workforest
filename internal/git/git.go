package git

import (
	"strings"

	"github.com/Moorad/workforest/internal/exec"
)

func Health() (bool, string) {
	out, err := exec.CmdOutput("git", "version")
	if err != nil {
		return false, ""
	}

	return true, strings.Trim(strings.Split(out, " ")[2], "\n")
}
