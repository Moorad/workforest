package worktrees

import (
	"path/filepath"

	"github.com/Moorad/workforest/internal/tui"
)

type Worktree struct {
	path string
}

func NewWorktree(path string) Worktree {
	return Worktree{path}
}

func (w *Worktree) Name() string {
	return filepath.Base(w.path)
}

func ItemizeWorktrees(wts []Worktree) []tui.Item {
	options := []tui.Item{}

	for _, worktree := range wts {
		options = append(options, tui.NewItem(worktree.Name(), worktree.path))
	}

	return options
}
