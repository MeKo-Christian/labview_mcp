package main

import (
	"context"
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// PingInput defines the input parameters for the ping tool
type PingInput struct {
	Message string `json:"message" jsonschema:"Message to echo back"`
}

// PingOutput defines the output for the ping tool
type PingOutput struct {
	Reply string `json:"reply" jsonschema:"Echoed reply message"`
}

// PingTool implements a simple echo/ping tool for testing
func PingTool(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input PingInput,
) (*mcp.CallToolResult, PingOutput, error) {
	reply := "Pong: " + input.Message
	return nil, PingOutput{Reply: reply}, nil
}

func main() {
	// Create MCP server with metadata
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "labview-mcp-go",
		Version: "0.1.0",
	}, nil)

	// Register the ping tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "ping",
		Description: "Echo test tool - responds with 'Pong: <message>'",
	}, PingTool)

	// Run server with stdio transport
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}
