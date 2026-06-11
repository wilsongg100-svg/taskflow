package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"taskflow/internal/domain"
	"taskflow/internal/infrastructure/worker"
	"time"

	"github.com/spf13/cobra"
)

var workers int

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the worker pool and process tasks",
	RunE: func(cmd *cobra.Command, args []string) error {
		// rebuild the pool with the requested worker count
		workerPool = worker.NewPool(workers, taskService)
		workerPool.Start()

		fmt.Printf("✓ Worker pool started with %d workers\n", workers)
		fmt.Println("  Waiting for tasks... (Ctrl+C to stop)")

		go func() {
			scanner := bufio.NewScanner(os.Stdin)
			fmt.Print("> ")
			for scanner.Scan() {
				line := strings.TrimSpace(scanner.Text())
				handleInput(cmd.Context(), line)
				fmt.Print("> ")
			}
		}()

		// still block on signal
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		fmt.Println("\nShutting down gracefully...")
		workerPool.Stop()
		fmt.Println("✓ All workers stopped")

		return nil
	},
}

func init() {
	startCmd.Flags().IntVarP(&workers, "workers", "w", 3, "Number of worker goroutines")
}

func handleInput(ctx context.Context, line string) {
	fmt.Printf("[DEBUG] parts: %v | taskService nil: %v\n", strings.Fields(line), taskService == nil)
	parts := strings.Fields(line)
	if len(parts) == 0 {
		return
	}

	switch parts[0] {
	case "create":
		// parts[1] is title, parts[2] is optional priority
		title := parts[1]
		priority := ""
		if parts[2] != "" {
			priority = parts[2]
		}
		task, err := taskService.Create(ctx, title, priority)
		if err != nil {
			fmt.Printf("Error creating task: %v\n", err)
		}
		if err := workerPool.Submit(task); err != nil {
			fmt.Printf("Error submitting task: %v\n", err)
			return
		}
		fmt.Printf("✓ Task created and submitted: %s\n", task.ID)
	case "test":
		//part[1] is number of test tasks to create
		startTime := time.Now()
		fmt.Printf("✓ Starting test: creating and submitting tasks...\n")
		taskCount := 5
		if len(parts) > 1 {
			fmt.Sscanf(parts[1], "%d", &taskCount)
		}
		tasks := make([]*domain.Task, taskCount)
		for i := 0; i < taskCount; i++ {
			task, err := taskService.Create(ctx, fmt.Sprintf("Test Task %d", i+1), "medium")
			if err != nil {
				fmt.Printf("Error creating task: %v\n", err)
				return
			}
			tasks[i] = task
		}
		for _, t := range tasks {
			if err := workerPool.Submit(t); err != nil {
				fmt.Printf("Error submitting task: %v\n", err)
				return
			}
			fmt.Printf("✓ Test task created and submitted: %s\n", t.ID)
		}
		endTime := time.Now()
		fmt.Printf("✓ All test tasks created and submitted in %v\n", endTime.Sub(startTime))
	case "list":
		tasks, err := taskService.List(ctx)
		if err != nil {
			fmt.Printf("Error listing tasks: %v\n", err)
			return
		}
		if len(tasks) == 0 {
			fmt.Println("No tasks found.")
			return
		}
		fmt.Printf("%-36s  %-8s  %-8s  %s\n", "ID", "STATUS", "PRIORITY", "TITLE")
		fmt.Printf("%-36s  %-8s  %-8s  %s\n", "----", "------", "--------", "-----")
		for _, t := range tasks {
			fmt.Printf("%-36s  %-8s  %-8s  %s\n", t.ID, t.Status, t.Priority, t.Title)
		}
	case "status":
		if len(parts) < 2 {
			fmt.Println("Usage: status <task-id>")
			return
		}
		task, err := taskService.Get(ctx, domain.TaskID(parts[1]))
		if err != nil {
			fmt.Printf("Error getting task: %v\n", err)
			return
		}
		fmt.Printf("ID:         %s\n", task.ID)
		fmt.Printf("Title:      %s\n", task.Title)
		fmt.Printf("Status:     %s\n", task.Status)
		fmt.Printf("Priority:   %s\n", task.Priority)
	case "help":
		fmt.Println("Commands: create <title> [priority], list, status <id>, exit")
	case "exit":
		fmt.Println("Bye!")
		os.Exit(0)
	default:
		fmt.Printf("Unknown command: %s\n", parts[0])
	}
}
