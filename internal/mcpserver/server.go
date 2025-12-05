package mcpserver

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/MeKo-Christian/labview_mcp/internal/python"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var (
	ErrStdioNotInitialized  = errors.New("stdio not initialized")
	ErrServerNotInitialized = errors.New("server not initialized")
)

type Options struct {
	Name           string
	Version        string
	PythonExecutor *python.Executor
}

type Server struct {
	s              *mcp.Server
	stdio          *mcp.StdioTransport
	http           *mcp.StreamableHTTPHandler
	sse            *mcp.SSEHandler
	pythonExecutor *python.Executor
}

func New(opts Options) *Server {
	if opts.Name == "" {
		opts.Name = "labview-mcp-go"
	}

	if opts.Version == "" {
		opts.Version = "0.1.0"
	}

	instructions := `# LabVIEW MCP Server

This server provides programmatic control over LabVIEW applications through the Model Context Protocol.

## Available Tools:

### Module Management
- **start_module**: Initialize the LabVIEW DQMH server module (must be called first)
- **stop_module**: Stop the LabVIEW server and clean up resources

### VI Creation & Manipulation
- **new_vi**: Create a new VI in the LabVIEW IDE
- **add_object**: Add objects to block diagrams or front panels (500+ object types)
- **connect_objects**: Wire two terminals together on the block diagram
- **save_vi**: Save a VI to disk

### Debugging & Inspection
- **get_object_terminals**: Get terminal information for wiring
- **get_vi_error_list**: Get compile errors for the current VI

### Utility
- **echo**: Simple connectivity test tool
- **ping**: Echo test with formatted response

## Workflow:
1. Start with start_module to initialize LabVIEW
2. Create VIs with new_vi
3. Add controls, indicators, and functions with add_object
4. Wire components with connect_objects
5. Check for errors with get_vi_error_list
6. Save your work with save_vi
7. Clean up with stop_module when done`

	s := mcp.NewServer(&mcp.Implementation{
		Name:    opts.Name,
		Version: opts.Version,
	}, &mcp.ServerOptions{
		Instructions: instructions,
	})

	// Register all tools
	registerTools(s)

	// Set global Python executor for tool handlers
	globalPythonExecutor = opts.PythonExecutor

	return &Server{
		s:              s,
		stdio:          &mcp.StdioTransport{},
		pythonExecutor: opts.PythonExecutor,
	}
}

func registerTools(s *mcp.Server) {
	// Ping tool for basic connectivity
	type PingInput struct {
		Message string `json:"message" jsonschema:"Message to echo back"`
	}
	type PingOutput struct {
		Reply string `json:"reply" jsonschema:"Echoed reply message"`
	}

	mcp.AddTool(s, &mcp.Tool{
		Name:        "ping",
		Description: "Echo test tool - responds with 'Pong: <message>'",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input PingInput) (*mcp.CallToolResult, PingOutput, error) {
		reply := "Pong: " + input.Message
		return nil, PingOutput{Reply: reply}, nil
	})

	// Module management
	mcp.AddTool(s, &mcp.Tool{
		Name:        "start_module",
		Description: "Initialize the LabVIEW DQMH server module. Must be called before any other LabVIEW-specific functions.",
	}, StartModuleTool)

	mcp.AddTool(s, &mcp.Tool{
		Name:        "stop_module",
		Description: "Stop the LabVIEW server and clean up all temporary data. Objects cannot be accessed after stopping.",
	}, StopModuleTool)

	// VI creation & manipulation
	mcp.AddTool(s, &mcp.Tool{
		Name:        "new_vi",
		Description: "Create a new VI in the LabVIEW IDE and return a reference for further modifications.",
	}, NewVITool)

	mcp.AddTool(s, &mcp.Tool{
		Name:        "add_object",
		Description: "Add an object to the block diagram, structure sub-diagram, or front panel. Supports over 500 object types including structures, controls/indicators, and functions.",
	}, AddObjectTool)

	mcp.AddTool(s, &mcp.Tool{
		Name:        "connect_objects",
		Description: "Connect two terminals of two objects with a wire on the block diagram. Use get_object_terminals to find terminal indices before wiring.",
	}, ConnectObjectsTool)

	mcp.AddTool(s, &mcp.Tool{
		Name:        "save_vi",
		Description: "Save a VI to the specified path.",
	}, SaveVITool)

	// Debugging & inspection
	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_object_terminals",
		Description: "Return terminal names/descriptions and their indices for use in connect_objects.",
	}, GetObjectTerminalsTool)

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_vi_error_list",
		Description: "Return the current VI error list (compile errors) in text format. Use this to verify if operations succeeded and identify issues.",
	}, GetVIErrorListTool)

	// Utility
	mcp.AddTool(s, &mcp.Tool{
		Name:        "echo",
		Description: "Simple echo tool for testing MCP connectivity. Does not require LabVIEW integration.",
	}, EchoTool)
}

func (srv *Server) RunStdio(ctx context.Context) error {
	if srv.s == nil || srv.stdio == nil {
		return ErrStdioNotInitialized
	}

	return srv.s.Run(ctx, srv.stdio)
}

// RunHTTP starts the HTTP/SSE server
func (srv *Server) RunHTTP(ctx context.Context, addr string) error {
	if srv.s == nil {
		return ErrServerNotInitialized
	}

	// Create HTTP handler that returns our MCP server for each session
	srv.http = mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server {
		return srv.s
	}, nil)

	// Create legacy SSE handler for backward compatibility
	srv.sse = mcp.NewSSEHandler(func(*http.Request) *mcp.Server {
		return srv.s
	}, &mcp.SSEOptions{})

	// Create a handler that routes requests appropriately
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Handle health endpoint
		if r.URL.Path == "/health" {
			if r.Method != http.MethodGet {
				http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"status":"healthy","service":"labview-mcp"}`))
			return
		}

		// Handle legacy SSE endpoint for backward compatibility
		if r.URL.Path == "/sse" {
			srv.sse.ServeHTTP(w, r)
			return
		}

		// All other routes go through the MCP handler
		srv.http.ServeHTTP(w, r)
	})

	// Apply security headers
	var finalHandler http.Handler = secureHeaders(handler)

	// Create HTTP server
	httpServer := &http.Server{
		Addr:              addr,
		Handler:           finalHandler,
		ReadHeaderTimeout: 10 * time.Second,
	}

	// Channel to signal server shutdown
	done := make(chan error, 1)

	// Start server in a goroutine
	go func() {
		slog.Info("Starting LabVIEW MCP HTTP/SSE server", slog.String("addr", addr))

		err := httpServer.ListenAndServe()

		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			done <- err
		}

		done <- nil
	}()

	// Wait for context cancellation or server error
	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		slog.Info("Shutting down LabVIEW MCP HTTP/SSE server...")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := httpServer.Shutdown(shutdownCtx)
		if err != nil {
			return err
		}

		return nil
	}
}

// secureHeaders adds basic security headers and MCP protocol version to every response
func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("Cache-Control", "no-store")
		// MCP Protocol Version header as recommended by specification
		w.Header().Set("MCP-Protocol-Version", "2024-11-05")
		next.ServeHTTP(w, r)
	})
}
