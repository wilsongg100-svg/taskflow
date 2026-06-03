package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var priority string

var createCmd = &cobra.Command{
	Use:   "create [title]",
	Short: "Create a new task",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		title := args[0]
		// this is where you call your application layer:
		// task, err := taskService.Create(title, priority)
		fmt.Printf("Task created: %q (priority: %s)\n", title, priority)
		return nil
	},
}

func init() {
	createCmd.Flags().StringVarP(&priority, "priority", "p", "medium", "Task priority (low|medium|high)")
}
