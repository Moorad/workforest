package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Moorad/workforest/internal/config"
	"github.com/Moorad/workforest/internal/git"
	"github.com/Moorad/workforest/internal/tmux"
	"github.com/Moorad/workforest/internal/tui"
	"github.com/Moorad/workforest/internal/worktrees"
	"github.com/fatih/color"
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
				panic("Failed to resolve passed path or current path")
			}
			wts, err := git.GetWorktrees(path)

			isWorktreePicked := false
			var workingPath string
			var configPath string
			if len(wts) > 1 {
				choice := tui.PromptList("Pick a worktree to switch to", worktrees.ItemizeWorktrees(wts))

				if choice == (tui.Item{}) {
					fmt.Println("No worktree picked, exiting")
					os.Exit(0)
				}
				workingPath = choice.Value
				isWorktreePicked = true
			} else {
				defaultPath, err := config.GetDefaultConfigPath()
				if err != nil {
					panic(err)
				}
				workingPath = defaultPath
			}

			println("working path:", workingPath)

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

			configExists, err := config.CheckConfigExists(configPath)
			if err != nil {
				panic(err)
			}

			if !configExists {
				color.Red("No config file was found for the following directory: %s", configPath)
				os.Exit(1)
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

			cfg := config.Parse(configPath)

			fmt.Println(workingPath)

			sessionName := cfg.Name

			if isWorktreePicked {
				sessionName += "_" + filepath.Base(workingPath)
			}

			if tmux.DoesSessionExist(sessionName) {
				tmux.DirectSwitchOrAttach(sessionName)
			}

			session, err := tmux.NewSession(sessionName, workingPath)
			if err != nil {
				fmt.Println(err)
				panic(fmt.Errorf("failed to create session: %e", err))
			}

			for i, window := range cfg.Windows {

				if i == 0 {
					err = session.RenameWindowMain(window.Name)
				} else {
					err = session.NewWindow(window.Name)
				}

				if err != nil {
					panic("Failed to create/rename new window")
				}

				if window.Command != "" {
					err := session.SendKeys(window.Name, window.Command, "Enter")
					if err != nil {
						panic("Failed to send keys to window")
					}
				}

			}

			err = session.SwitchOrAttach()
			if err != nil {
				panic("Failed to attach")
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
