package cmd

import (
	"os"
	"path/filepath"

	"github.com/Moorad/workforest/internal/config"
	"github.com/Moorad/workforest/internal/git"
	"github.com/Moorad/workforest/internal/tmux"
	"github.com/Moorad/workforest/internal/tui"
	"github.com/Moorad/workforest/internal/utils"
	"github.com/Moorad/workforest/internal/worktrees"
	"github.com/spf13/cobra"
)

var (
	cfgPath string
	rootCmd = &cobra.Command{
		Use:   "workforest",
		Short: "A tool to simplify git worktree + tmux workflows",
		Long:  `A tool that allows you to easily create, switch and manage tmux session mapped to git worktrees`,
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			path, err := config.GetPath(args)
			if err != nil {
				utils.Panic("Failed to resolve config path", err)
			}

			wts, err := git.GetWorktrees(path)
			if err != nil {
				utils.Panic("Failed to get worktrees", err)
			}

			utils.Debug("wt", wts)

			var workingPath string
			var configPath string
			isWorktreePicked := false

			if len(wts) > 1 {
				choice, err := tui.PromptList("Pick a worktree to switch to", worktrees.ItemizeWorktrees(wts))
				if err != nil {
					utils.Panic("An error occured while prompting list", err)
				}

				if choice == (tui.Item{}) {
					utils.GracefullyExit("No worktree picked")
				}

				workingPath = choice.Value
				isWorktreePicked = true
			} else {
				defaultPath, err := config.GetDefaultConfigPath()
				if err != nil {
					utils.Panic("Failed to resolve default config path", err)
				}
				workingPath = defaultPath
			}

			utils.Debug("Working path: %s", workingPath)

			if cfgPath != "" {
				configPath = cfgPath
			} else {
				if isWorktreePicked {
					configPath = filepath.Dir(workingPath)
				} else {
					configPath = workingPath
				}

				configPath = filepath.Join(configPath, "/.workforest.yml")
			}

			utils.Debug("Config path: %s", configPath)

			configExists, err := config.CheckConfigExists(configPath)
			if err != nil {
				utils.Panic("Failed to check config existence", err)
			}

			if !configExists {
				utils.PanicMsg("No config file was found for the following directory: %s", configPath)
			}

			// 		if resolvedCfgPath == "" {
			// 			color.Red("No config file was found in current or parent directory and no --config path was provided.")
			// 			var confirm bool
			// 			err := huh.NewConfirm().
			// 				Title("Would you like to create an empty tmux session?").
			// 				Affirmative("Yes").
			// 				Negative("No").
			// 				Value(&confirm).Run()
			// 			if err != nil {
			// 				panic("Failed to prompt user")
			// 			}
			//
			// 			if !confirm {
			// 				return
			// 			}
			//

			cfg, err := config.Parse(configPath)
			if err != nil {
				utils.Panic("Failed to parse config file", err)
			}

			sessionName := cfg.Name

			if isWorktreePicked {
				sessionName += "_" + filepath.Base(workingPath)
			}

			utils.Debug("Session name: %s", sessionName)

			if tmux.DoesSessionExist(sessionName) {
				err := tmux.DirectSwitchOrAttach(sessionName)
				if err != nil {
					utils.Panic("Failed to direct switch/attach to tmux session", err)
				}
			}

			session, err := tmux.NewSession(sessionName, workingPath)
			if err != nil {
				utils.Panic("Failed to create tmux session", err)
			}

			for i, window := range cfg.Windows {

				if i == 0 {
					err = session.RenameWindowMain(window.Name)
				} else {
					err = session.NewWindow(window.Name)
				}

				if err != nil {
					utils.Panic("Failed to create/rename new tmux window", err)
				}

				if window.Command != "" {
					err := session.SendKeys(window.Name, window.Command, "Enter")
					if err != nil {
						utils.Panic("Failed to send keys to tmux window", err)
					}
				}

			}

			err = session.SwitchOrAttach()
			if err != nil {
				utils.Panic("Failed to switch/attack to tmux session", err)
			}
		},
	}
)

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgPath, "config", "", "config file (default <current directory>/.workforest.yml)")
}
