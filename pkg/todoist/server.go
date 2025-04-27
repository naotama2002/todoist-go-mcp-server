package todoist

import (
	"net/http"

	"github.com/mark3labs/mcp-go/server"
	"github.com/sirupsen/logrus"
)

// Server represents a Todoist MCP server
type Server struct {
	mcpServer *server.MCPServer
	tools     *ToolProvider
	logger    *logrus.Logger
	httpServer *http.Server
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

	return &Server{
		mcpServer: mcpServer,
		tools:     tools,
		logger:    logger,
	}
}

// Start starts the Todoist MCP server
func (s *Server) Start(addr string) error {
	// Register tools
	for _, tool := range s.tools.GetTools() {
		s.mcpServer.AddTool(tool.Tool, tool.Handler)
		s.logger.WithField("tool", tool.Tool.Name).Info("Registered tool")
	}

	// Create HTTP server
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok","message":"Todoist MCP Server is running"}`))
	})

	s.httpServer = &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	// Start the server
	s.logger.WithField("addr", addr).Info("Starting Todoist MCP server")
	return s.httpServer.ListenAndServe()
}
