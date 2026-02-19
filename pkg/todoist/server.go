package todoist

import (
	"context"
	"net/http"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/naotama2002/todoist-go-mcp-server/pkg/toolsets"
	"github.com/sirupsen/logrus"
)

// Server represents a Todoist MCP server
type Server struct {
	mcpServer    *mcp.Server
	tools        *ToolProvider
	logger       *logrus.Logger
	httpServer   *http.Server
	toolsetGroup *toolsets.ToolsetGroup
}

// NewServer creates a new Todoist MCP server
func NewServer(token string, logger *logrus.Logger) *Server {
	if logger == nil {
		logger = logrus.New()
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	}

	tools := NewToolProvider(token, logger)

	// Create a new MCP server with default options
	mcpServer := mcp.NewServer(
		&mcp.Implementation{
			Name:    "todoist-mcp-server",
			Version: "v0.1.0",
		},
		nil,
	)

	// Create default toolset group
	toolsetGroup := createDefaultToolsetGroup(tools, false)

	return &Server{
		mcpServer:    mcpServer,
		tools:        tools,
		logger:       logger,
		toolsetGroup: toolsetGroup,
	}
}

// createDefaultToolsetGroup creates the default toolset group for Todoist
func createDefaultToolsetGroup(tp *ToolProvider, readOnly bool) *toolsets.ToolsetGroup {
	group := toolsets.NewToolsetGroup(readOnly)

	// Create task management toolset
	taskToolset := toolsets.NewToolset("tasks", "Todoist task management tools")
	taskToolset.AddReadTools(
		toolsets.NewServerTool(tp.GetTaskFilterRules(), tp.HandleGetTaskFilterRules),
		toolsets.NewServerTool(tp.GetTasks(), tp.HandleGetTasks),
		toolsets.NewServerTool(tp.GetTask(), tp.HandleGetTask),
	)

	if !readOnly {
		taskToolset.AddWriteTools(
			toolsets.NewServerTool(tp.CreateTask(), tp.HandleCreateTask),
			toolsets.NewServerTool(tp.UpdateTask(), tp.HandleUpdateTask),
			toolsets.NewServerTool(tp.CloseTask(), tp.HandleCloseTask),
			toolsets.NewServerTool(tp.DeleteTask(), tp.HandleDeleteTask),
		)
	}

	// Create project management toolset
	projectToolset := toolsets.NewToolset("projects", "Todoist project management tools")
	projectToolset.AddReadTools(
		toolsets.NewServerTool(tp.GetProjects(), tp.HandleGetProjects),
		toolsets.NewServerTool(tp.GetProject(), tp.HandleGetProject),
	)

	// Add toolsets to the group
	group.AddToolset(taskToolset)
	group.AddToolset(projectToolset)

	// Enable all toolsets by default
	if err := group.EnableToolsets([]string{"all"}); err != nil {
		return nil
	}

	return group
}

// Start starts the Todoist MCP server over HTTP
func (s *Server) Start(ctx context.Context, addr string) error {
	// Register tools using toolset group
	s.toolsetGroup.RegisterTools(s.mcpServer)
	s.logger.Info("Registered tools from toolset group")

	// Create HTTP server with StreamableHTTPHandler
	handler := mcp.NewStreamableHTTPHandler(
		func(r *http.Request) *mcp.Server { return s.mcpServer },
		&mcp.StreamableHTTPOptions{JSONResponse: true},
	)

	mux := http.NewServeMux()
	mux.Handle("/", handler)

	s.httpServer = &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	// Start the server in a goroutine
	go func() {
		s.logger.WithField("addr", addr).Info("Starting Todoist MCP server over HTTP")
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.WithError(err).Fatal("HTTP server ListenAndServe error")
		}
	}()

	// Listen for context cancellation to gracefully shut down
	<-ctx.Done()
	s.logger.Info("Shutting down HTTP server...")

	// Create a deadline to wait for.
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
		s.logger.WithError(err).Fatal("HTTP server Shutdown failed")
	}
	s.logger.Info("HTTP server gracefully stopped")
	return nil
}

// StartStdio starts the Todoist MCP server over stdio
func (s *Server) StartStdio(ctx context.Context) error {
	// Register tools using toolset group
	s.toolsetGroup.RegisterTools(s.mcpServer)
	s.logger.Info("Registered tools from toolset group")

	s.logger.Info("Starting Todoist MCP server over stdio")
	return s.mcpServer.Run(ctx, &mcp.StdioTransport{})
}
