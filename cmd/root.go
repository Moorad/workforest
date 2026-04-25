package cmd

import (
	"os"

	"charm.land/huh/v2"
	"github.com/Moorad/workforest/internal/config"
	"github.com/Moorad/workforest/internal/tmux"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	cfgPath string
	rootCmd = &cobra.Command{
		Use:   "workforest",
		Short: "A tool to simplify git worktree + tmux workflows",
		Long:  `A tool that allows you to easily create, switch and manage tmux session mapped to git worktrees`,
		Run: func(cmd *cobra.Command, args []string) {
			resolvedCfgPath, err := config.ResolveConfigPath(cfgPath)
			if err != nil {
				panic(err)
			}

			if resolvedCfgPath == "" {
				color.Red("No config file was found in current or parent directory and no --config path was provided.")
				var confirm bool
				err := huh.NewConfirm().
					Title("Would you like to create an empty tmux session?").
					Affirmative("Yes").
					Negative("No").
					Value(&confirm).Run()
				if err != nil {
					panic("Failed to prompt user")
				}

				if !confirm {
					return
				}

				session, err := tmux.NewSession("default-workforest")
				if err != nil {
					panic("Failed to create session")
				}

				err = session.SwitchOrAttach()
				if err != nil {
					panic("Failed to create session")
				}

			}

			cfg := config.Parse(resolvedCfgPath)

			session, err := tmux.NewSession(cfg.Name)
			if err != nil {
				panic("Failed to create session")
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
