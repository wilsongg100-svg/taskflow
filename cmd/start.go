package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

var workers int

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the worker pool and process tasks",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("Starting %d workers...\n", workers)

		// pool := worker.NewPool(workers, repo)
		// pool.Start(ctx)

		// block until Ctrl+C — this is the Go idiom for graceful shutdown
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		fmt.Println("\nShutting down gracefully...")
		// pool.Stop() triggers context cancel + waits for WaitGroup
		return nil
	},
}

func init() {
	startCmd.Flags().IntVarP(&workers, "workers", "w", 3, "Number of worker goroutines")
}
