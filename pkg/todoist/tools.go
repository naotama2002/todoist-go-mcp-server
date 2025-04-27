package todoist

import (
	"github.com/mark3labs/mcp-go/server"
	"github.com/sirupsen/logrus"
)

// ToolProvider provides MCP tools for Todoist
type ToolProvider struct {
	client *Client
	logger *logrus.Logger
}

// NewToolProvider creates a new ToolProvider
func NewToolProvider(token string, logger *logrus.Logger) *ToolProvider {
	return &ToolProvider{
		client: NewClient(token, logger),
		logger: logger,
	}
}

// GetTools returns all Todoist tools
func (tp *ToolProvider) GetTools() []server.ServerTool {
	// Return all tools
	return []server.ServerTool{
		{
			Tool:    tp.GetTasks(),
			Handler: tp.HandleGetTasks,
		},
		// Add other tools here
	}
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}
