package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/naotama2002/todoist-go-mcp-server/pkg/log"
	"github.com/naotama2002/todoist-go-mcp-server/pkg/todoist"
)

func main() {
	// Parse command line flags
	addr := flag.String("addr", ":8080", "Address to listen on")
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

	// Create and start the server
	server := todoist.NewServer(*token, logger)

	// Handle graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		logger.Info("Shutting down...")
		os.Exit(0)
	}()

	// Start the server
	logger.WithField("addr", *addr).Info("Starting Todoist MCP server")
	if err := server.Start(*addr); err != nil {
		logger.WithError(err).Fatal("Failed to start server")
	}
}
