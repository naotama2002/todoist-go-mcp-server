package todoist

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/mark3labs/mcp-go/server"
	"github.com/sirupsen/logrus"
)

// Server represents a Todoist MCP server
type Server struct {
	mcpServer  *server.MCPServer
	tools      *ToolProvider
	logger     *logrus.Logger
	httpServer *http.Server
	stdioServer *server.StdioServer
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

	return &Server{
		mcpServer:  mcpServer,
		tools:      tools,
		logger:     logger,
		stdioServer: stdioServer,
	}
}

// Start starts the Todoist MCP server over HTTP
func (s *Server) Start(addr string) error {
	// Register tools
	for _, tool := range s.tools.GetTools() {
		s.mcpServer.AddTool(tool.Tool, tool.Handler)
		s.logger.WithField("tool", tool.Tool.Name).Info("Registered tool")
	}

	// Create HTTP server with MCP protocol handler
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"status":"ok","message":"Todoist MCP Server is running"}`))
			return
		}

		// Read request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// Handle MCP message
		ctx := context.Background()
		response := s.mcpServer.HandleMessage(ctx, body)

		// Convert response to JSON
		responseJSON, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Write response
		w.Header().Set("Content-Type", "application/json")
		w.Write(responseJSON)
	})

	s.httpServer = &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	// Start the server
	s.logger.WithField("addr", addr).Info("Starting Todoist MCP server over HTTP")
	return s.httpServer.ListenAndServe()
}

// StartStdio starts the Todoist MCP server over stdio
func (s *Server) StartStdio(ctx context.Context, in io.Reader, out io.Writer) error {
	// Register tools
	for _, tool := range s.tools.GetTools() {
		s.mcpServer.AddTool(tool.Tool, tool.Handler)
		s.logger.WithField("tool", tool.Tool.Name).Info("Registered tool")
	}

	s.logger.Info("Starting Todoist MCP server over stdio")
	return s.stdioServer.Listen(ctx, in, out)
}
