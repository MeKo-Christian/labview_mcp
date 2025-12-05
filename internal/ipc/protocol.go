// Package ipc provides inter-process communication abstractions for the LabVIEW MCP server.
// This package defines the protocol for communication between the Go MCP server and the LabVIEW bridge.
package ipc

import "encoding/json"

// MessageType identifies the type of IPC message
type MessageType string

const (
	// MessageTypeRequest is sent from Go MCP server to LabVIEW bridge
	MessageTypeRequest MessageType = "request"
	// MessageTypeResponse is sent from LabVIEW bridge back to Go MCP server
	MessageTypeResponse MessageType = "response"
)

// Request represents a request from the Go MCP server to the LabVIEW bridge.
// Requests are line-delimited JSON messages sent over stdin/stdout or TCP.
type Request struct {
	// Type is always "request"
	Type MessageType `json:"type"`
	// ID is a unique identifier for request/response correlation
	ID string `json:"id"`
	// Tool is the name of the LabVIEW tool to invoke (e.g., "new_vi", "add_object")
	Tool string `json:"tool"`
	// Payload contains the tool-specific input parameters as JSON
	Payload json.RawMessage `json:"payload"`
}

// Response represents a response from the LabVIEW bridge back to the Go MCP server.
type Response struct {
	// Type is always "response"
	Type MessageType `json:"type"`
	// ID matches the request ID for correlation
	ID string `json:"id"`
	// Payload contains the tool-specific output data as JSON
	Payload json.RawMessage `json:"payload,omitempty"`
	// Error contains error information if the operation failed
	Error string `json:"error,omitempty"`
}

// NewRequest creates a new IPC request with the given parameters.
func NewRequest(id, tool string, payload interface{}) (*Request, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return &Request{
		Type:    MessageTypeRequest,
		ID:      id,
		Tool:    tool,
		Payload: payloadBytes,
	}, nil
}

// NewResponse creates a new IPC response with the given parameters.
func NewResponse(id string, payload interface{}, errorMsg string) (*Response, error) {
	var payloadBytes json.RawMessage
	var err error

	if payload != nil {
		payloadBytes, err = json.Marshal(payload)
		if err != nil {
			return nil, err
		}
	}

	return &Response{
		Type:    MessageTypeResponse,
		ID:      id,
		Payload: payloadBytes,
		Error:   errorMsg,
	}, nil
}
