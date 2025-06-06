package todoist

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/mark3labs/mcp-go/server"
	"github.com/naotama2002/todoist-go-mcp-server/pkg/toolsets"
	"github.com/sirupsen/logrus"
)

// Server represents a Todoist MCP server
type Server struct {
	mcpServer    *server.MCPServer
	tools        *ToolProvider
	logger       *logrus.Logger
	httpServer   *http.Server
	stdioServer  *server.StdioServer
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
	mcpServer := server.NewMCPServer(
		"todoist-mcp-server",
		"v0.1.0",
		server.WithToolCapabilities(true),
		server.WithLogging(),
	)

	// Create a new stdio server
	stdioServer := server.NewStdioServer(mcpServer)

	// Create default toolset group
	toolsetGroup := createDefaultToolsetGroup(tools, false)

	return &Server{
		mcpServer:    mcpServer,
		tools:        tools,
		logger:       logger,
		stdioServer:  stdioServer,
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

	// Create HTTP server with MCP protocol handler
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("Content-Type", "application/json")
			if _, err := w.Write([]byte(`{"status":"ok","message":"Todoist MCP Server is running"}`)); err != nil {
				s.logger.Errorf("Error writing response: %v", err)
			}
			return
		}

		// Read request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer func() {
			if err := r.Body.Close(); err != nil {
				s.logger.Errorf("Error closing request body: %v", err)
			}
		}()

		// Handle MCP message
		response := s.mcpServer.HandleMessage(r.Context(), body)

		// Convert response to JSON
		responseJSON, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Write response
		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write(responseJSON); err != nil {
			s.logger.Errorf("Error writing response: %v", err)
		}
	})

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
func (s *Server) StartStdio(ctx context.Context, in io.Reader, out io.Writer) error {
	// Register tools using toolset group
	s.toolsetGroup.RegisterTools(s.mcpServer)
	s.logger.Info("Registered tools from toolset group")

	s.logger.Info("Starting Todoist MCP server over stdio")
	return s.stdioServer.Listen(ctx, in, out)
}
