package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tasks",
	RunE: func(cmd *cobra.Command, args []string) error {
		// tasks, err := taskService.List()
		fmt.Println("ID        STATUS    TITLE")
		fmt.Println("────────  ────────  ──────────────────")
		// for _, t := range tasks { fmt.Printf(...) }
		return nil
	},
}
