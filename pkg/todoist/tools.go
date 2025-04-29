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
		{
			Tool:    tp.GetTask(),
			Handler: tp.HandleGetTask,
		},
		{
			Tool:    tp.CreateTask(),
			Handler: tp.HandleCreateTask,
		},
		{
			Tool:    tp.UpdateTask(),
			Handler: tp.HandleUpdateTask,
		},
		{
			Tool:    tp.CloseTask(),
			Handler: tp.HandleCloseTask,
		},
		{
			Tool:    tp.DeleteTask(),
			Handler: tp.HandleDeleteTask,
		},
		{
			Tool:    tp.GetProjects(),
			Handler: tp.HandleGetProjects,
		},
		{
			Tool:    tp.GetProject(),
			Handler: tp.HandleGetProject,
		},
		// Add other tools here
	}
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}
