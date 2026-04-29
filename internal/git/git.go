package git

import (
	"strings"

	"github.com/Moorad/workforest/internal/exec"
	"github.com/Moorad/workforest/internal/worktrees"
)

func Health() (bool, string) {
	out, err := exec.NewCommand("git", "version").Output()
	if err != nil {
		return false, ""
	}

	return true, strings.Trim(strings.Split(out, " ")[2], "\n")
}

func GetWorktrees(path string) ([]worktrees.Worktree, error) {
	cmd := exec.NewCommand("git", "worktree", "list", "--porcelain")

	cmd.Cmd.Dir = path

	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	splitOut := strings.Split(strings.Trim(out, "\n"), "\n\n")

	wts := []worktrees.Worktree{}

	for _, worktreeStr := range splitOut {
		splitText := strings.Split(strings.Trim(worktreeStr, "\n"), "\n")
		path := strings.Split(splitText[0], " ")[1]
		wts = append(wts, worktrees.NewWorktree(path))
	}

	return wts, nil
}
