package cmd

import (
	"fmt"
	"taskflow/internal/domain"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status [task-id]",
	Short: "Get the status of a specific task",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		task, err := taskService.Get(cmd.Context(), domain.TaskID(args[0]))
		if err != nil {
			return fmt.Errorf("failed to get task: %w", err)
		}

		fmt.Printf("ID:         %s\n", task.ID)
		fmt.Printf("Title:      %s\n", task.Title)
		fmt.Printf("Priority:   %s\n", task.Priority)
		fmt.Printf("Status:     %s\n", task.Status)
		fmt.Printf("Created:    %s\n", task.CreatedAt.Format("2006-01-02 15:04:05"))
		fmt.Printf("Updated:    %s\n", task.UpdatedAt.Format("2006-01-02 15:04:05"))

		return nil
	},
}
