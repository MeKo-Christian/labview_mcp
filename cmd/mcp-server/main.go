package main

import (
	"context"
	"flag"
	"log"
	"log/slog"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/MeKo-Christian/labview_mcp/internal/ipc"
	"github.com/MeKo-Christian/labview_mcp/internal/mcpserver"
)

var (
	transport   = flag.String("transport", "stdio", "Transport to use: stdio or http")
	httpAddr    = flag.String("http-addr", ":8787", "HTTP listen address (when --transport=http)")
	useBridge   = flag.Bool("use-bridge", true, "Spawn and use the LabVIEW bridge process")
	bridgePath  = flag.String("bridge-path", "", "Path to labview-bridge binary (auto-detected if empty)")
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
		slog.Bool("use_bridge", *useBridge),
	)

	// Setup graceful shutdown context
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// Setup bridge process if enabled
	var bridgeCmd *exec.Cmd
	var ipcClient *ipc.Client

	if *useBridge {
		// Determine bridge binary path
		bridgeBinPath := *bridgePath
		if bridgeBinPath == "" {
			// Auto-detect: look for bin/labview-bridge relative to this binary
			exePath, err := os.Executable()
			if err != nil {
				log.Fatalf("Failed to get executable path: %v", err)
			}
			exeDir := filepath.Dir(exePath)
			bridgeBinPath = filepath.Join(exeDir, "labview-bridge")
		}

		slog.Info("Starting LabVIEW bridge process", slog.String("path", bridgeBinPath))

		// Create bridge command
		bridgeCmd = exec.Command(bridgeBinPath)

		// Setup pipes for stdin/stdout
		stdinPipe, err := bridgeCmd.StdinPipe()
		if err != nil {
			log.Fatalf("Failed to create stdin pipe: %v", err)
		}

		stdoutPipe, err := bridgeCmd.StdoutPipe()
		if err != nil {
			log.Fatalf("Failed to create stdout pipe: %v", err)
		}

		// Bridge stderr goes to our stderr for logging
		bridgeCmd.Stderr = os.Stderr

		// Start bridge process
		if err := bridgeCmd.Start(); err != nil {
			log.Fatalf("Failed to start bridge: %v", err)
		}

		slog.Info("Bridge process started", slog.Int("pid", bridgeCmd.Process.Pid))

		// Create IPC client
		ipcClient = ipc.NewClient(ctx, stdoutPipe, stdinPipe)
	}

	// Create server with optional IPC client
	srv := mcpserver.New(mcpserver.Options{
		Name:      "labview-mcp-go",
		Version:   "0.1.0",
		IPCClient: ipcClient,
	})

	// Run server with selected transport
	err := error(nil)
	switch *transport {
	case "stdio":
		slog.Info("Running with stdio transport")
		err = srv.RunStdio(ctx)
		if err != nil {
			slog.Error("stdio transport failed", slog.Any("error", err))
		} else {
			slog.Info("stdio transport stopped")
		}
	case "http":
		slog.Info("Running with HTTP/SSE transport", slog.String("address", *httpAddr))
		err = srv.RunHTTP(ctx, *httpAddr)
		if err != nil {
			slog.Error("HTTP transport failed", slog.Any("error", err))
		} else {
			slog.Info("HTTP transport stopped")
		}
	default:
		log.Fatalf("Unknown transport: %s (use 'stdio' or 'http')", *transport)
	}

	// Cleanup
	if ipcClient != nil {
		slog.Info("Closing IPC client")
		ipcClient.Close()
	}

	if bridgeCmd != nil && bridgeCmd.Process != nil {
		slog.Info("Stopping bridge process")
		bridgeCmd.Process.Signal(syscall.SIGTERM)
		bridgeCmd.Wait()
	}

	if err != nil {
		os.Exit(1)
	}
}
