# IPC Protocol Documentation

## Overview

The LabVIEW MCP server uses a simple JSON-based IPC (Inter-Process Communication) protocol to communicate between the Go MCP server and the LabVIEW bridge process.

## Transport

The protocol uses **line-delimited JSON** over stdin/stdout (or optionally TCP for future flexibility).

Each message is a complete JSON object followed by a newline character (`\n`).

## Message Format

### Request (Go MCP → LabVIEW Bridge)

```json
{
  "type": "request",
  "id": "1",
  "tool": "new_vi",
  "payload": {}
}
```

**Fields:**
- `type` (string): Always `"request"`
- `id` (string): Unique identifier for request/response correlation
- `tool` (string): Name of the LabVIEW tool to invoke (e.g., `"new_vi"`, `"add_object"`)
- `payload` (object): Tool-specific input parameters as JSON

### Response (LabVIEW Bridge → Go MCP)

```json
{
  "type": "response",
  "id": "1",
  "payload": {
    "vi_id": 42,
    "result": "VI created successfully"
  },
  "error": ""
}
```

**Fields:**
- `type` (string): Always `"response"`
- `id` (string): Matches the request ID for correlation
- `payload` (object, optional): Tool-specific output data as JSON
- `error` (string, optional): Error message if the operation failed (empty string on success)

## Tool Payloads

Each tool has specific input and output payload formats. See [tools.md](tools.md) for complete schemas.

### Example: start_module

**Request:**
```json
{
  "type": "request",
  "id": "1",
  "tool": "start_module",
  "payload": {}
}
```

**Response:**
```json
{
  "type": "response",
  "id": "1",
  "payload": {
    "status": "module started",
    "module_already_running": false,
    "error_out": ""
  }
}
```

### Example: new_vi

**Request:**
```json
{
  "type": "request",
  "id": "2",
  "tool": "new_vi",
  "payload": {}
}
```

**Response:**
```json
{
  "type": "response",
  "id": "2",
  "payload": {
    "vi_id": 42,
    "result": "VI created successfully",
    "timed_out": false,
    "error_out": ""
  }
}
```

### Example: add_object

**Request:**
```json
{
  "type": "request",
  "id": "3",
  "tool": "add_object",
  "payload": {
    "diagram_id": 42,
    "object_name": "Numeric Control",
    "position_x": 100,
    "position_y": 200
  }
}
```

**Response:**
```json
{
  "type": "response",
  "id": "3",
  "payload": {
    "object_id": 123,
    "result": "Object added successfully",
    "timed_out": false,
    "error_out": ""
  }
}
```

### Example: Error Response

**Request:**
```json
{
  "type": "request",
  "id": "4",
  "tool": "connect_objects",
  "payload": {
    "vi_reference": 42,
    "from_object_reference": 123,
    "from_object_terminal_index": 0,
    "to_object_reference": 999,
    "to_object_terminal_index": 0
  }
}
```

**Response (Error):**
```json
{
  "type": "response",
  "id": "4",
  "error": "Object 999 not found"
}
```

## Protocol Flow

1. **Go MCP Server** receives tool call from Claude Desktop via MCP protocol
2. **Go MCP Server** creates IPC Request with unique ID
3. **Go MCP Server** sends line-delimited JSON request to LabVIEW bridge (stdout)
4. **LabVIEW Bridge** receives request via stdin
5. **LabVIEW Bridge** parses JSON and invokes appropriate LabVIEW VI
6. **LabVIEW Bridge** creates IPC Response with matching ID
7. **LabVIEW Bridge** sends line-delimited JSON response back (stdout)
8. **Go MCP Server** receives response via stdin
9. **Go MCP Server** correlates response by ID and returns result to Claude Desktop

## Implementation Notes

- **Request IDs** are generated sequentially starting from 1
- **Concurrent requests** are supported - responses are matched by ID
- **Timeouts** can be implemented on the Go side using context cancellation
- **Line buffering** ensures complete JSON messages are transmitted atomically
- **Error handling** uses the `error` field in responses; empty string means success

## Future Enhancements

- TCP transport option for network-based bridge communication
- Bidirectional event streaming (bridge → MCP server notifications)
- Request cancellation protocol
- Compression for large payloads (VI screenshots, etc.)
