package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var priority string

var createCmd = &cobra.Command{
	Use:   "create [title]",
	Short: "Create a new task and submit it to the pool",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		task, err := taskService.Create(cmd.Context(), args[0], priority)
		if err != nil {
			return fmt.Errorf("failed to create task: %w", err)
		}

		fmt.Printf("✓ Task created\n")
		fmt.Printf("  ID:       %s\n", task.ID)
		fmt.Printf("  Title:    %s\n", task.Title)
		fmt.Printf("  Priority: %s\n", task.Priority)
		fmt.Printf("  Status:   %s\n", task.Status)

		// submit to the worker pool so it gets processed
		if err := workerPool.Submit(task); err != nil {
			fmt.Printf("  [WARN] Could not submit to pool: %v\n", err)
		} else {
			fmt.Printf("  ✓ Submitted to worker pool\n")
		}

		return nil
	},
}

func init() {
	createCmd.Flags().StringVarP(&priority, "priority", "p", "medium", "Task priority (low|medium|high)")
}
