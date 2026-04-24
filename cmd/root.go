package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "workforest",
	Short: "A tool to simplify git worktree + tmux workflows",
	Long:  `A tool that allows you to easily create, switch and manage tmux session mapped to git worktrees`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}
