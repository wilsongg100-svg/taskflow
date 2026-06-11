package cmd

import (
	"fmt"
	"os"
	"taskflow/internal/application"
	"taskflow/internal/domain"
	"taskflow/internal/infrastructure/inmemory"
	"taskflow/internal/infrastructure/worker"

	"github.com/spf13/cobra"
)

var (
	taskService *application.TaskService
	workerPool  *worker.Pool
)

var rootCmd = &cobra.Command{
	Use:   "taskflow",
	Short: "A concurrent task pipeline system",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// build the dependency chain once, at startup
	// events channel is the bridge between the service and the rest of the system
	events := make(chan domain.DomainEvent, 100)
	repo := inmemory.NewRepository()
	taskService = application.NewTaskService(repo, events)
	workerPool = worker.NewPool(3, taskService)

	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(startCmd)
}
