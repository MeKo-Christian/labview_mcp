// Package main implements a simple LabVIEW bridge process that communicates
// with the MCP server via line-delimited JSON over stdin/stdout.
//
// This is a PoC implementation that returns canned responses to validate
// the IPC protocol before connecting to real LabVIEW.
package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/MeKo-Christian/labview_mcp/internal/ipc"
)

func main() {
	// Setup logging to stderr (stdout is reserved for IPC)
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	// Detect environment
	mode := "stub"
	platform := runtime.GOOS
	if platform == "windows" {
		mode = "windows-ready"
		slog.Info("LabVIEW Bridge starting",
			slog.String("mode", mode),
			slog.String("platform", platform),
			slog.String("note", "Python wrappers required for real LabVIEW integration"),
		)
	} else {
		slog.Info("LabVIEW Bridge starting",
			slog.String("mode", mode),
			slog.String("platform", platform),
			slog.String("note", "LabVIEW integration requires Windows - using canned responses"),
		)
	}


	// Create scanner for reading line-delimited JSON from stdin
	scanner := bufio.NewScanner(os.Stdin)

	// Process requests in a loop
	for scanner.Scan() {
		line := scanner.Bytes()

		// Parse request
		var req ipc.Request
		if err := json.Unmarshal(line, &req); err != nil {
			slog.Error("Failed to parse request", slog.Any("error", err))
			continue
		}

		slog.Info("Received request",
			slog.String("id", req.ID),
			slog.String("tool", req.Tool),
		)

		// Generate canned response based on tool name
		resp := handleRequest(&req)

		// Send response as line-delimited JSON
		respBytes, err := json.Marshal(resp)
		if err != nil {
			slog.Error("Failed to marshal response", slog.Any("error", err))
			continue
		}

		fmt.Fprintf(os.Stdout, "%s\n", respBytes)

		slog.Info("Sent response", slog.String("id", resp.ID))
	}

	if err := scanner.Err(); err != nil {
		slog.Error("Scanner error", slog.Any("error", err))
		os.Exit(1)
	}

	slog.Info("LabVIEW Bridge stopped")
}

// handleRequest processes a request by calling Python wrapper on Windows or returning canned response
func handleRequest(req *ipc.Request) *ipc.Response {
	var payload interface{}
	var errorMsg string

	// On Windows, try to call Python wrapper scripts
	if runtime.GOOS == "windows" {
		payload, errorMsg = callPythonWrapper(req)
		if errorMsg == "" {
			// Successfully called Python wrapper
			resp, err := ipc.NewResponse(req.ID, payload, "")
			if err != nil {
				slog.Error("Failed to create response", slog.Any("error", err))
				resp = &ipc.Response{
					Type:  ipc.MessageTypeResponse,
					ID:    req.ID,
					Error: fmt.Sprintf("failed to create response: %v", err),
				}
			}
			return resp
		}
		// If Python wrapper failed, fall through to canned responses
		slog.Warn("Python wrapper failed, using canned response", slog.String("error", errorMsg))
	}

	// Generate tool-specific canned responses (for non-Windows or fallback)
	switch req.Tool {
	case "start_module":
		payload = map[string]interface{}{
			"status":                 "bridge: module started",
			"module_already_running": false,
			"error_out":              "",
		}

	case "stop_module":
		payload = map[string]interface{}{
			"status":    "bridge: module stopped",
			"error_out": "",
		}

	case "new_vi":
		payload = map[string]interface{}{
			"vi_id":     42,
			"result":    "bridge: VI created with ID 42",
			"timed_out": false,
			"error_out": "",
		}

	case "add_object":
		// Parse input to get object details
		var input struct {
			ObjectName string `json:"object_name"`
			PositionX  int    `json:"position_x"`
			PositionY  int    `json:"position_y"`
		}
		if err := json.Unmarshal(req.Payload, &input); err == nil {
			payload = map[string]interface{}{
				"object_id": 123,
				"result":    fmt.Sprintf("bridge: added '%s' at (%d,%d)", input.ObjectName, input.PositionX, input.PositionY),
				"timed_out": false,
				"error_out": "",
			}
		} else {
			errorMsg = "failed to parse add_object input"
		}

	case "connect_objects":
		payload = map[string]interface{}{
			"result":    "bridge: objects connected",
			"timed_out": false,
			"error_out": "",
		}

	case "save_vi":
		// Parse input to get path
		var input struct {
			Path string `json:"path"`
		}
		if err := json.Unmarshal(req.Payload, &input); err == nil {
			payload = map[string]interface{}{
				"path_out":  fmt.Sprintf("bridge: saved to %s", input.Path),
				"timed_out": false,
				"error_out": "",
			}
		} else {
			errorMsg = "failed to parse save_vi input"
		}

	case "get_object_terminals":
		payload = map[string]interface{}{
			"result":    "bridge: terminals [0: input, 1: output, 2: error]",
			"timed_out": false,
			"error_out": "",
		}

	case "get_vi_error_list":
		payload = map[string]interface{}{
			"result":    "bridge: no compile errors",
			"timed_out": false,
			"error_out": "",
		}

	default:
		errorMsg = fmt.Sprintf("unknown tool: %s", req.Tool)
	}

	// Create response
	resp, err := ipc.NewResponse(req.ID, payload, errorMsg)
	if err != nil {
		slog.Error("Failed to create response", slog.Any("error", err))
		// Return error response
		resp = &ipc.Response{
			Type:  ipc.MessageTypeResponse,
			ID:    req.ID,
			Error: fmt.Sprintf("failed to create response: %v", err),
		}
	}

	return resp
}

// callPythonWrapper executes the Python wrapper script for the given tool
func callPythonWrapper(req *ipc.Request) (interface{}, string) {
	// Get the executable directory to find wrapper scripts
	exePath, err := os.Executable()
	if err != nil {
		return nil, fmt.Sprintf("failed to get executable path: %v", err)
	}
	exeDir := filepath.Dir(exePath)

	// Wrapper scripts are in ../wrapper relative to the binary
	wrapperDir := filepath.Join(exeDir, "..", "wrapper")
	scriptPath := filepath.Join(wrapperDir, req.Tool+".py")

	// Check if script exists
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		return nil, fmt.Sprintf("wrapper script not found: %s", scriptPath)
	}

	slog.Info("Calling Python wrapper",
		slog.String("tool", req.Tool),
		slog.String("script", scriptPath),
	)

	// Prepare Python command
	cmd := exec.Command("python", scriptPath)

	// Convert request payload to JSON for stdin
	var inputJSON []byte
	if req.Payload != nil {
		inputJSON, err = json.Marshal(req.Payload)
		if err != nil {
			return nil, fmt.Sprintf("failed to marshal payload: %v", err)
		}
	} else {
		inputJSON = []byte("{}")
	}

	// Set up stdin/stdout/stderr
	cmd.Stdin = bytes.NewReader(inputJSON)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Execute the command
	err = cmd.Run()

	// Log stderr if present
	if stderr.Len() > 0 {
		slog.Warn("Python wrapper stderr", slog.String("stderr", stderr.String()))
	}

	if err != nil {
		return nil, fmt.Sprintf("python execution failed: %v (stderr: %s)", err, stderr.String())
	}

	// Parse the JSON output
	var result interface{}
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		return nil, fmt.Sprintf("failed to parse python output: %v (output: %s)", err, stdout.String())
	}

	slog.Info("Python wrapper succeeded", slog.String("tool", req.Tool))
	return result, ""
}
