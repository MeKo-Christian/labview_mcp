// Package python provides direct Python script execution for LabVIEW wrapper scripts.
package python

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"log/slog"
)

// Executor handles Python script execution with platform detection.
type Executor struct {
	wrapperDir string
	isWindows  bool
}

// NewExecutor creates a new Python executor.
// It auto-detects the wrapper directory relative to the executable.
func NewExecutor() (*Executor, error) {
	// Get executable path to find wrapper directory
	exePath, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("failed to get executable path: %w", err)
	}
	exeDir := filepath.Dir(exePath)

	// Wrapper scripts are in ../wrapper relative to the binary
	wrapperDir := filepath.Join(exeDir, "..", "wrapper")

	return &Executor{
		wrapperDir: wrapperDir,
		isWindows:  runtime.GOOS == "windows",
	}, nil
}

// Call executes a Python wrapper script on Windows or returns canned response on other platforms.
func (e *Executor) Call(ctx context.Context, tool string, input interface{}) (json.RawMessage, error) {
	if e.isWindows {
		return e.callPython(ctx, tool, input)
	}
	return e.cannedResponse(tool, input)
}

// callPython executes the Python wrapper script for the given tool.
func (e *Executor) callPython(ctx context.Context, tool string, input interface{}) (json.RawMessage, error) {
	scriptPath := filepath.Join(e.wrapperDir, tool+".py")

	// Check if script exists
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("wrapper script not found: %s", scriptPath)
	}

	slog.Info("Calling Python wrapper",
		slog.String("tool", tool),
		slog.String("script", scriptPath),
	)

	// Prepare Python command
	cmd := exec.CommandContext(ctx, "python", scriptPath)

	// Convert input to JSON for stdin
	var inputJSON []byte
	var err error
	if input != nil {
		inputJSON, err = json.Marshal(input)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal input: %w", err)
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
		return nil, fmt.Errorf("python execution failed: %w (stderr: %s)", err, stderr.String())
	}

	// Parse the JSON output
	var result json.RawMessage
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		return nil, fmt.Errorf("failed to parse python output: %w (output: %s)", err, stdout.String())
	}

	slog.Info("Python wrapper succeeded", slog.String("tool", tool))
	return result, nil
}

// cannedResponse returns tool-specific stub responses for non-Windows platforms.
func (e *Executor) cannedResponse(tool string, input interface{}) (json.RawMessage, error) {
	var response interface{}

	switch tool {
	case "start_module":
		response = map[string]interface{}{
			"status":                 "stub: module started",
			"module_already_running": false,
			"error_out":              "",
		}

	case "stop_module":
		response = map[string]interface{}{
			"status":    "stub: module stopped",
			"error_out": "",
		}

	case "new_vi":
		response = map[string]interface{}{
			"vi_id":     42,
			"result":    "stub: VI created",
			"timed_out": false,
			"error_out": "",
		}

	case "add_object":
		// Parse input to get object details
		var addInput struct {
			ObjectName string `json:"object_name"`
			PositionX  int    `json:"position_x"`
			PositionY  int    `json:"position_y"`
		}
		inputBytes, _ := json.Marshal(input)
		if err := json.Unmarshal(inputBytes, &addInput); err == nil {
			response = map[string]interface{}{
				"object_id": 123,
				"result":    fmt.Sprintf("stub: added '%s' at (%d,%d)", addInput.ObjectName, addInput.PositionX, addInput.PositionY),
				"timed_out": false,
				"error_out": "",
			}
		} else {
			return nil, fmt.Errorf("failed to parse add_object input")
		}

	case "connect_objects":
		response = map[string]interface{}{
			"result":    "stub: connected objects",
			"timed_out": false,
			"error_out": "",
		}

	case "save_vi":
		// Parse input to get path
		var saveInput struct {
			Path string `json:"path"`
		}
		inputBytes, _ := json.Marshal(input)
		if err := json.Unmarshal(inputBytes, &saveInput); err == nil {
			response = map[string]interface{}{
				"path_out":  fmt.Sprintf("stub: saved to %s", saveInput.Path),
				"timed_out": false,
				"error_out": "",
			}
		} else {
			return nil, fmt.Errorf("failed to parse save_vi input")
		}

	case "get_object_terminals":
		response = map[string]interface{}{
			"result":    "stub: terminals [0: input, 1: output, 2: error]",
			"timed_out": false,
			"error_out": "",
		}

	case "get_vi_error_list":
		response = map[string]interface{}{
			"result":    "stub: no errors",
			"timed_out": false,
			"error_out": "",
		}

	default:
		return nil, fmt.Errorf("unknown tool: %s", tool)
	}

	// Marshal response to JSON
	data, err := json.Marshal(response)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal canned response: %w", err)
	}

	return json.RawMessage(data), nil
}
