package cmd

import (
	"fmt"

	"github.com/Moorad/workforest/internal/fzf"
	"github.com/Moorad/workforest/internal/git"
	"github.com/Moorad/workforest/internal/tmux"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	VERSION string
	COMMIT  string
	DATE    string
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version of workforest and its dependencies",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version: %s (%s - %s)\n", VERSION, COMMIT, DATE)

		isGitInstalled, gitVersion := git.Health()
		isFzfInstalled, fzfVersion := fzf.Health()
		isTmuxInstalled, tmuxVersion := tmux.Health()

		fmt.Printf("\nDependencies:\n\tgit: %s\n\tfzf: %s\n\ttmux: %s\n", formatInstallationMsg(isGitInstalled, gitVersion), formatInstallationMsg(isFzfInstalled, fzfVersion), formatInstallationMsg(isTmuxInstalled, tmuxVersion))
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func formatInstallationMsg(isInstalled bool, version string) string {
	if isInstalled {
		return fmt.Sprintf("%s (%s)", color.GreenString("Installed"), version)
	}

	return color.RedString("Missing")
}
