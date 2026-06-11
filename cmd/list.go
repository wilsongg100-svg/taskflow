package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tasks",
	RunE: func(cmd *cobra.Command, args []string) error {
		tasks, err := taskService.List(cmd.Context())
		if err != nil {
			return fmt.Errorf("failed to list tasks: %w", err)
		}

		if len(tasks) == 0 {
			fmt.Println("No tasks found.")
			return nil
		}

		fmt.Printf("%-36s  %-8s  %-8s  %s\n", "ID", "STATUS", "PRIORITY", "TITLE")
		fmt.Printf("%-36s  %-8s  %-8s  %s\n", "----", "------", "--------", "-----")
		for _, t := range tasks {
			fmt.Printf("%-36s  %-8s  %-8s  %s\n", t.ID, t.Status, t.Priority, t.Title)
		}

		return nil
	},
}
