package cmd

import (
	"fmt"
	"strings"

	"github.com/Moorad/workforest/internal/tmux"
	"github.com/Moorad/workforest/internal/tui"
	"github.com/Moorad/workforest/internal/utils"
	"github.com/spf13/cobra"
)

var switchCmd = &cobra.Command{
	Use:   "switch",
	Short: "Switch globally between tmux session",
	Run: func(cmd *cobra.Command, args []string) {
		sessions, err := tmux.ListSessions()
		if err != nil {
			utils.Panic("Failed to list tmux sessions", err)
		}

		utils.Debug("Session info: %v", sessions)

		items := []tui.Item{}

		for _, sessionInfo := range sessions {
			sessionName := sessionInfo.Name
			if strings.Contains(sessionName, "_") {
				parts := strings.Split(sessionName, "_")
				sessionName = fmt.Sprintf("%s (%s)", parts[0], parts[1])
			}

			label := fmt.Sprintf("%s - %d windows", sessionName, sessionInfo.WindowCount)

			if sessionInfo.IsAttached {
				label += " [current]"
			}

			items = append(items, tui.NewItem(label, sessionInfo.Name))
		}

		choice, err := tui.PromptList("Choose a session to switch to", items)
		if err != nil {
			utils.Panic("Failed to prompt list", err)
		}

		if choice.IsEmpty() {
			utils.GracefullyExit("No session picked")
		}

		err = tmux.DirectSwitchOrAttach(choice.Value)
		if err != nil {
			utils.Panic("Failed to attach/switch session", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(switchCmd)
}
