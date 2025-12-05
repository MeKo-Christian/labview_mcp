package main

import (
	"context"
	"flag"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/MeKo-Christian/labview_mcp/internal/mcpserver"
	"github.com/MeKo-Christian/labview_mcp/internal/python"
)

var (
	transport = flag.String("transport", "stdio", "Transport to use: stdio or http")
	httpAddr  = flag.String("http-addr", ":8787", "HTTP listen address (when --transport=http)")
)

// PingInput defines the input parameters for the ping tool
type PingInput struct {
	Message string `json:"message" jsonschema:"Message to echo back"`
}

// PingOutput defines the output for the ping tool
type PingOutput struct {
	Reply string `json:"reply" jsonschema:"Echoed reply message"`
}

func main() {
	flag.Parse()

	// Setup logging
	// For stdio transport, logs must go to stderr to keep stdout clean
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	slog.Info("Starting LabVIEW MCP Server",
		slog.String("transport", *transport),
		slog.String("version", "0.1.0"),
	)

	// Setup graceful shutdown context
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// Create Python executor
	pythonExec, err := python.NewExecutor()
	if err != nil {
		log.Fatalf("Failed to create Python executor: %v", err)
	}

	slog.Info("Python executor ready")

	// Create server with Python executor
	srv := mcpserver.New(mcpserver.Options{
		Name:           "labview-mcp-go",
		Version:        "0.1.0",
		PythonExecutor: pythonExec,
	})

	// Run server with selected transport
	var runErr error
	switch *transport {
	case "stdio":
		slog.Info("Running with stdio transport")
		runErr = srv.RunStdio(ctx)
		if runErr != nil {
			slog.Error("stdio transport failed", slog.Any("error", runErr))
		} else {
			slog.Info("stdio transport stopped")
		}
	case "http":
		slog.Info("Running with HTTP/SSE transport", slog.String("address", *httpAddr))
		runErr = srv.RunHTTP(ctx, *httpAddr)
		if runErr != nil {
			slog.Error("HTTP transport failed", slog.Any("error", runErr))
		} else {
			slog.Info("HTTP transport stopped")
		}
	default:
		log.Fatalf("Unknown transport: %s (use 'stdio' or 'http')", *transport)
	}

	if runErr != nil {
		os.Exit(1)
	}
}
