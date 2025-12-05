// Package mcpserver provides tool definitions for the LabVIEW MCP server.
// This file contains the input/output structures and handler functions for all
// MCP tools that will be exposed to Claude Desktop.
package mcpserver

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/MeKo-Christian/labview_mcp/internal/python"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Tool input and output structures follow the schema documented in docs/tools.md
// Each tool has:
// - An Input struct with JSON tags for MCP parameter binding
// - An Output struct with JSON tags for MCP result binding
// - A handler function with signature: func(context.Context, *mcp.CallToolRequest, InputType) (*mcp.CallToolResult, OutputType, error)

// globalPythonExecutor is set by the server during initialization
var globalPythonExecutor *python.Executor

// callPython is a helper function that calls the Python executor if available,
// otherwise returns the provided stub response.
func callPython[T any](ctx context.Context, toolName string, input interface{}, stubOutput T) (T, error) {
	if globalPythonExecutor != nil {
		respPayload, err := globalPythonExecutor.Call(ctx, toolName, input)
		if err != nil {
			var zero T
			return zero, fmt.Errorf("Python execution failed: %w", err)
		}

		var output T
		if err := json.Unmarshal(respPayload, &output); err != nil {
			var zero T
			return zero, fmt.Errorf("failed to unmarshal response: %w", err)
		}

		return output, nil
	}

	// Return stub when no executor (shouldn't happen, but safe fallback)
	return stubOutput, nil
}

// ============================================================================
// Module Management Tools
// ============================================================================

// StartModuleInput defines the input parameters for start_module tool
type StartModuleInput struct {
	// No input parameters required for start_module
}

// StartModuleOutput defines the output for start_module tool
type StartModuleOutput struct {
	Status               string `json:"status" jsonschema:"Status of the module startup"`
	ModuleAlreadyRunning bool   `json:"module_already_running,omitempty" jsonschema:"Whether the module was already running"`
	ErrorOut             string `json:"error_out,omitempty" jsonschema:"Error message if any"`
}

// StartModuleTool initializes the LabVIEW DQMH server module.
// This tool must be called before any other LabVIEW-specific functions.
func StartModuleTool(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input StartModuleInput,
) (*mcp.CallToolResult, StartModuleOutput, error) {
	output, err := callPython(ctx, "start_module", input, StartModuleOutput{
		Status:               "stub: module started",
		ModuleAlreadyRunning: false,
		ErrorOut:             "",
	})
	return nil, output, err
}

// StopModuleInput defines the input parameters for stop_module tool
type StopModuleInput struct {
	// No input parameters required for stop_module
}

// StopModuleOutput defines the output for stop_module tool
type StopModuleOutput struct {
	Status   string `json:"status" jsonschema:"Status of the module shutdown"`
	ErrorOut string `json:"error_out,omitempty" jsonschema:"Error message if any"`
}

// StopModuleTool stops the LabVIEW server and cleans up all temporary data.
// Objects cannot be accessed after stopping.
func StopModuleTool(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input StopModuleInput,
) (*mcp.CallToolResult, StopModuleOutput, error) {
	output, err := callPython(ctx, "stop_module", input, StopModuleOutput{
		Status:   "stub: module stopped",
		ErrorOut: "",
	})
	return nil, output, err
}

// ============================================================================
// VI Creation & Manipulation Tools
// ============================================================================

// NewVIInput defines the input parameters for new_vi tool
type NewVIInput struct {
	// No input parameters required for new_vi
}

// NewVIOutput defines the output for new_vi tool
type NewVIOutput struct {
	VIID     int    `json:"vi_id" jsonschema:"Integer reference to the created VI"`
	Result   string `json:"result" jsonschema:"Result message"`
	TimedOut bool   `json:"timed_out,omitempty" jsonschema:"Whether the operation timed out"`
	ErrorOut string `json:"error_out,omitempty" jsonschema:"Error message if any"`
}

// NewVITool creates a new VI in the LabVIEW IDE and returns a reference for further modifications.
func NewVITool(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input NewVIInput,
) (*mcp.CallToolResult, NewVIOutput, error) {
	output, err := callPython(ctx, "new_vi", input, NewVIOutput{
		VIID:     0,
		Result:   "stub: VI created",
		TimedOut: false,
		ErrorOut: "",
	})
	return nil, output, err
}

// AddObjectInput defines the input parameters for add_object tool
type AddObjectInput struct {
	DiagramID  int    `json:"diagram_id" jsonschema:"Reference to the diagram or VI"`
	ObjectName string `json:"object_name" jsonschema:"Name of the LabVIEW object to add (e.g., 'Add', 'While Loop #1', 'Numeric Control')"`
	PositionX  int    `json:"position_x" jsonschema:"X coordinate position"`
	PositionY  int    `json:"position_y" jsonschema:"Y coordinate position"`
}

// AddObjectOutput defines the output for add_object tool
type AddObjectOutput struct {
	ObjectID int    `json:"object_id" jsonschema:"Integer reference to the created object"`
	Result   string `json:"result" jsonschema:"Result message"`
	TimedOut bool   `json:"timed_out,omitempty" jsonschema:"Whether the operation timed out"`
	ErrorOut string `json:"error_out,omitempty" jsonschema:"Error message if any"`
}

// AddObjectTool adds an object to the block diagram, structure sub-diagram, or front panel.
// Over 500 object types are supported including structures, controls/indicators, and functions.
func AddObjectTool(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input AddObjectInput,
) (*mcp.CallToolResult, AddObjectOutput, error) {
	output, err := callPython(ctx, "add_object", input, AddObjectOutput{
		ObjectID: 0,
		Result:   fmt.Sprintf("stub: object '%s' added at (%d,%d)", input.ObjectName, input.PositionX, input.PositionY),
		TimedOut: false,
		ErrorOut: "",
	})
	return nil, output, err
}

// ConnectObjectsInput defines the input parameters for connect_objects tool
type ConnectObjectsInput struct {
	VIReference             int `json:"vi_reference" jsonschema:"Reference ID of the VI"`
	FromObjectReference     int `json:"from_object_reference" jsonschema:"Reference ID of source object"`
	FromObjectTerminalIndex int `json:"from_object_terminal_index" jsonschema:"Terminal index on the source object"`
	ToObjectReference       int `json:"to_object_reference" jsonschema:"Reference ID of destination object"`
	ToObjectTerminalIndex   int `json:"to_object_terminal_index" jsonschema:"Terminal index on the destination object"`
}

// ConnectObjectsOutput defines the output for connect_objects tool
type ConnectObjectsOutput struct {
	Result   string `json:"result" jsonschema:"Result message indicating connection status"`
	TimedOut bool   `json:"timed_out,omitempty" jsonschema:"Whether the operation timed out"`
	ErrorOut string `json:"error_out,omitempty" jsonschema:"Error message if any"`
}

// ConnectObjectsTool connects two terminals of two objects with a wire on the block diagram.
// Use get_object_terminals to find terminal indices before wiring.
func ConnectObjectsTool(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input ConnectObjectsInput,
) (*mcp.CallToolResult, ConnectObjectsOutput, error) {
	output, err := callPython(ctx, "connect_objects", input, ConnectObjectsOutput{
		Result:   "stub: connected objects",
		TimedOut: false,
		ErrorOut: "",
	})
	return nil, output, err
}

// SaveVIInput defines the input parameters for save_vi tool
type SaveVIInput struct {
	VIID int    `json:"vi_id" jsonschema:"Reference ID of the VI to save"`
	Path string `json:"path" jsonschema:"File path including filename (empty string saves to previous location)"`
}

// SaveVIOutput defines the output for save_vi tool
type SaveVIOutput struct {
	PathOut  string `json:"path_out" jsonschema:"Path where the VI was saved"`
	TimedOut bool   `json:"timed_out,omitempty" jsonschema:"Whether the operation timed out"`
	ErrorOut string `json:"error_out,omitempty" jsonschema:"Error message if any"`
}

// SaveVITool saves a VI to the specified path.
func SaveVITool(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input SaveVIInput,
) (*mcp.CallToolResult, SaveVIOutput, error) {
	output, err := callPython(ctx, "save_vi", input, SaveVIOutput{
		PathOut:  "stub: " + input.Path,
		TimedOut: false,
		ErrorOut: "",
	})
	return nil, output, err
}

// ============================================================================
// Debugging & Inspection Tools
// ============================================================================

// GetObjectTerminalsInput defines the input parameters for get_object_terminals tool
type GetObjectTerminalsInput struct {
	ObjectID int `json:"object_id" jsonschema:"Reference ID of the object"`
}

// GetObjectTerminalsOutput defines the output for get_object_terminals tool
type GetObjectTerminalsOutput struct {
	Result   string `json:"result" jsonschema:"String containing terminal names, descriptions, and indices"`
	TimedOut bool   `json:"timed_out,omitempty" jsonschema:"Whether the operation timed out"`
	ErrorOut string `json:"error_out,omitempty" jsonschema:"Error message if any"`
}

// GetObjectTerminalsTool returns terminal names/descriptions and their indices for use in connect_objects.
func GetObjectTerminalsTool(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input GetObjectTerminalsInput,
) (*mcp.CallToolResult, GetObjectTerminalsOutput, error) {
	output, err := callPython(ctx, "get_object_terminals", input, GetObjectTerminalsOutput{
		Result:   "stub: terminal info for object",
		TimedOut: false,
		ErrorOut: "",
	})
	return nil, output, err
}

// GetVIErrorListInput defines the input parameters for get_vi_error_list tool
type GetVIErrorListInput struct {
	VIReference int `json:"vi_reference" jsonschema:"Reference ID of the VI"`
}

// GetVIErrorListOutput defines the output for get_vi_error_list tool
type GetVIErrorListOutput struct {
	Result   string `json:"result" jsonschema:"String containing error descriptions and locations"`
	TimedOut bool   `json:"timed_out,omitempty" jsonschema:"Whether the operation timed out"`
	ErrorOut string `json:"error_out,omitempty" jsonschema:"Error message if any"`
}

// GetVIErrorListTool returns the current VI error list (compile errors) in text format.
// Use this to verify if operations succeeded and identify issues.
func GetVIErrorListTool(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input GetVIErrorListInput,
) (*mcp.CallToolResult, GetVIErrorListOutput, error) {
	output, err := callPython(ctx, "get_vi_error_list", input, GetVIErrorListOutput{
		Result:   "stub: no errors",
		TimedOut: false,
		ErrorOut: "",
	})
	return nil, output, err
}

// ============================================================================
// Utility Tools
// ============================================================================

// EchoInput defines the input parameters for echo tool
type EchoInput struct {
	Text string `json:"text" jsonschema:"Text to echo back"`
}

// EchoOutput defines the output for echo tool
type EchoOutput struct {
	Response string `json:"response" jsonschema:"Echoed response"`
}

// EchoTool is a simple echo tool for testing MCP connectivity.
// This tool does not require LabVIEW integration.
func EchoTool(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input EchoInput,
) (*mcp.CallToolResult, EchoOutput, error) {
	return nil, EchoOutput{
		Response: "You said: " + input.Text,
	}, nil
}
