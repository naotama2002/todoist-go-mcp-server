package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/naotama2002/todoist-go-mcp-server/pkg/log"
	"github.com/naotama2002/todoist-go-mcp-server/pkg/todoist"
)

func main() {
	// Parse command line flags
	mode := flag.String("mode", "http", "Server mode: 'http' or 'stdio'")
	addr := flag.String("addr", ":8080", "Address to listen on (HTTP mode only)")
	token := flag.String("token", "", "Todoist API token")
	flag.Parse()

	// Create logger
	logger := log.NewLogger()

	// Get token from environment variable if not provided
	if *token == "" {
		*token = os.Getenv("TODOIST_API_TOKEN")
		if *token == "" {
			logger.Fatal("Todoist API token is required. Set it with -token flag or TODOIST_API_TOKEN environment variable.")
		}
	}

	// Create the server
	server := todoist.NewServer(*token, logger)

	// Handle graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Start the server based on the mode
	switch *mode {
	case "http":
		// Start the server in HTTP mode
		logger.WithField("addr", *addr).Info("Starting Todoist MCP server in HTTP mode")
		if err := server.Start(*addr); err != nil {
			logger.WithError(err).Fatal("Failed to start server")
		}
	case "stdio":
		// Start the server in stdio mode
		logger.Info("Starting Todoist MCP server in stdio mode")

		// Create channels for errors
		errCh := make(chan error, 1)

		// Start the server in a goroutine
		go func() {
			errCh <- server.StartStdio(ctx, os.Stdin, os.Stdout)
		}()

		// Wait for shutdown signal or error
		select {
		case <-ctx.Done():
			logger.Info("Shutting down...")
		case err := <-errCh:
			if err != nil {
				logger.WithError(err).Fatal("Failed to run server")
			}
		}
	default:
		logger.Fatalf("Invalid mode: %s. Must be 'http' or 'stdio'", *mode)
	}
}
