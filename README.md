# LabVIEW MCP Server

A Model Context Protocol (MCP) server that enables Large Language Models (LLMs) like Claude to interact with LabVIEW through programmatic scripting.

## Overview

This Go-based MCP server provides a clean, high-performance interface for LLMs to create, modify, and debug LabVIEW VIs. It uses a bridge architecture to integrate with LabVIEW's DQMH scripting modules via COM on Windows.

**Status:** ✅ Implementation complete, ready for Windows + LabVIEW validation

## Architecture

```text
Claude Desktop (MCP Client)
    ↓ stdio (JSON-RPC)
Go MCP Server (bin/mcp-server)
    ↓ spawns child process
LabVIEW Bridge (bin/labview-bridge)
    ↓ IPC: line-delimited JSON over stdin/stdout
    ↓ On Windows: exec.Command("python", "wrapper/<tool>.py")
Python Wrapper Scripts (wrapper/*.py)
    ↓ WIN32COM (pywin32)
LabVIEW Application
    ↓ COM Call2 API
LabVIEW DQMH Module (LabVIEW_Server/Scripting Server/*.vi)
```

### Platform Behavior

- **Linux/macOS**: Bridge returns canned responses (for development/testing)
- **Windows**: Bridge executes Python wrappers that communicate with LabVIEW via COM

## Prerequisites

- **Go 1.21+** - For building the MCP server and bridge
- **LabVIEW 2020+** - With DQMH scripting modules (Windows only)
- **Python 3.8+** - For Windows COM integration (wrapper scripts)
- **uv** - Fast Python package manager (recommended) or pip
- **pywin32** - Python Windows COM library (auto-installed via uv/pip)
- **Claude Desktop** - Or any MCP-compatible client

## Installation

### 1. Clone and Build

```bash
git clone https://github.com/MeKo-Christian/labview_mcp.git
cd labview_mcp

# Download dependencies
go mod download

# Build both binaries
go build -o bin/mcp-server ./cmd/mcp-server
go build -o bin/labview-bridge ./cmd/labview-bridge

# For Windows (cross-compile)
GOOS=windows GOARCH=amd64 go build -o bin/mcp-server.exe ./cmd/mcp-server
GOOS=windows GOARCH=amd64 go build -o bin/labview-bridge.exe ./cmd/labview-bridge
```

### 2. Windows Setup (for LabVIEW integration)

```bash
# Install Python dependencies using uv (recommended)
uv sync

# Or use pip
pip install pywin32

# Deploy to Windows
# Copy these files to your Windows machine:
#   - bin/mcp-server.exe
#   - bin/labview-bridge.exe
#   - wrapper/*.py
#   - LabVIEW_Server/
#   - pyproject.toml (for uv sync on Windows)
#   - uv.lock (for reproducible builds)
```

### 3. Configure Claude Desktop

Add to your Claude Desktop configuration:

**Windows**: `%APPDATA%\Claude\claude_desktop_config.json`
**macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`
**Linux**: `~/.config/Claude/claude_desktop_config.json`

```json
{
  "mcpServers": {
    "labview": {
      "command": "C:\\path\\to\\labview_mcp\\bin\\mcp-server.exe",
      "args": []
    }
  }
}
```

### 4. Restart Claude Desktop

The LabVIEW tools will appear in Claude's available tools list.

## Available Tools

### Module Management

#### `start_module`

Start the LabVIEW DQMH scripting server. **Must be called first** before using other tools.

```json
No parameters required
```

Returns: Module status and configuration

#### `stop_module`

Stop the LabVIEW scripting server and clean up resources.

```json
No parameters required
```

Returns: Shutdown confirmation

### VI Creation

#### `new_vi`

Create a new VI in LabVIEW IDE.

```json
No parameters required
```

Returns: `vi_id` (integer reference to the created VI)

#### `save_vi`

Save a VI to disk.

```json
{
  "vi_id": 123,
  "path": "C:\\path\\to\\MyVI.vi"  // Empty string saves to previous location
}
```

Returns: Saved file path

### VI Modification

#### `add_object`

Add an object (control, function, structure) to a VI's block diagram or front panel.

```json
{
  "diagram_id": 123,
  "object_name": "While Loop",
  "position_x": 100,
  "position_y": 50
}
```

Supports 500+ object types - see [docs/tools.md](docs/tools.md) for complete list.

Returns: `object_id` (integer reference)

#### `connect_objects`

Connect two terminals with a wire on the block diagram.

```json
{
  "vi_reference": 123,
  "from_object_reference": 456,
  "from_object_terminal_index": 0,
  "to_object_reference": 789,
  "to_object_terminal_index": 1
}
```

Returns: Connection confirmation

### Debugging

#### `get_object_terminals`

Get terminal names, descriptions, and indices for wiring.

```json
{
  "object_id": 456
}
```

Returns: List of terminals with indices and descriptions

#### `get_vi_error_list`

Get compile errors (equivalent to clicking the broken run arrow).

```json
{
  "vi_reference": 123
}
```

Returns: Text list of errors and their locations

## Directory Structure

```text
labview_mcp/
├── cmd/
│   ├── mcp-server/         # MCP server entrypoint
│   └── labview-bridge/     # Bridge process for LabVIEW integration
├── internal/
│   ├── mcpserver/          # MCP server implementation
│   │   ├── server.go
│   │   └── tools.go        # 8 tool definitions
│   └── ipc/                # IPC protocol (line-delimited JSON)
│       ├── protocol.go
│       └── client.go
├── wrapper/                # Python COM wrapper scripts (8 tools)
│   ├── start_module.py
│   ├── stop_module.py
│   ├── new_vi.py
│   ├── add_object.py
│   ├── connect_objects.py
│   ├── get_object_terminals.py
│   ├── save_vi.py
│   └── get_vi_error_list.py
├── LabVIEW_Server/         # LabVIEW DQMH modules (29 Request VIs)
│   └── Scripting Server/
├── docs/
│   ├── tools.md            # Complete tool schema reference (29 tools)
│   ├── vi-mapping.md       # LabVIEW VI → Python tool mapping
│   └── ipc-protocol.md     # IPC protocol specification
├── bin/                    # Build artifacts
├── go.mod
├── go.sum
└── README.md
```

## Development

### Running Locally

```bash
# Run MCP server (standalone mode with stubs)
go run ./cmd/mcp-server

# Run bridge standalone
go run ./cmd/labview-bridge

# Test IPC round-trip
echo '{"type":"request","id":"test1","tool":"ping","payload":{"message":"hello"}}' | go run ./cmd/labview-bridge
```

### Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...
```

### Adding New Tools

1. Add tool schema to [docs/tools.md](docs/tools.md)
2. Add handler in `internal/mcpserver/tools.go`
3. Create Python wrapper in `wrapper/<tool>.py`
4. Register in `cmd/mcp-server/main.go`
5. Rebuild binaries

## How It Works

### IPC Protocol

The Go MCP server and bridge communicate via **line-delimited JSON** over stdin/stdout:

**Request Format:**

```json
{
  "type": "request",
  "id": "unique-id",
  "tool": "tool_name",
  "payload": { "param1": "value1" }
}
```

**Response Format:**

```json
{
  "type": "response",
  "id": "unique-id",
  "payload": { "result": "value" },
  "error": ""
}
```

### Python Wrapper Pattern

Each wrapper script ([wrapper/](wrapper/)):

1. Reads JSON from stdin
2. Connects to LabVIEW: `win32com.client.Dispatch("LabVIEW.Application")`
3. Gets VI reference: `GetVIReference("LabVIEW_Server/Scripting Server/<tool>.vi")`
4. Calls VI via COM: `vi.Call2(param_names, param_values, ...)`
5. Extracts results from `param_values` array
6. Returns JSON to stdout

### LabVIEW DQMH Architecture

The LabVIEW server uses the **Delacor Queued Message Handler (DQMH)** pattern:

- **Request VIs**: Public API (e.g., `new_vi.vi`, `add_object.vi`)
- **Do VIs**: Internal implementation (e.g., `Do Create Control.vi`)
- **Module VIs**: Lifecycle management (`Start Module.vi`, `Stop Module.vi`)

All calls use the **Request-and-Wait-for-Reply** pattern with the `Call2` API.

See [docs/vi-mapping.md](docs/vi-mapping.md) for complete VI documentation.

## Implementation Status

✅ **Completed:**

- Phase A1-A3: Go MCP server with 8 core tools
- Phase B1: IPC protocol and client library
- Phase B2: Bridge process with canned responses
- Phase B3: Python wrapper integration for Windows

⏳ **Pending:**

- Windows + LabVIEW validation testing
- Remaining 21 tools (optional)
- Performance optimization

See [PLAN.md](PLAN.md) for detailed progress tracking.

## Troubleshooting

### Tools not appearing in Claude

1. Check Claude Desktop logs: `~/.config/Claude/logs/`
2. Verify MCP server path in `claude_desktop_config.json`
3. Ensure binary has execute permissions
4. Restart Claude Desktop

### Bridge fails on Windows

1. Verify Python is in PATH: `python --version`
2. Check pywin32 is installed: `pip show pywin32`
3. Ensure LabVIEW is running
4. Check bridge logs (stderr output)

### LabVIEW COM errors

1. Ensure LabVIEW 2020+ is installed
2. Enable scripting: Tools → Options → VI Server → Enable VI Scripting
3. Check DQMH module is loaded: `LabVIEW_Server/Scripting Server/`
4. Verify VI paths in wrapper scripts

## Documentation

- [docs/tools.md](docs/tools.md) - Complete tool schema reference (all 29 tools)
- [docs/vi-mapping.md](docs/vi-mapping.md) - LabVIEW VI architecture and mapping
- [docs/ipc-protocol.md](docs/ipc-protocol.md) - IPC protocol specification
- [PLAN.md](PLAN.md) - Implementation roadmap and progress

## Contributing

Contributions are welcome! This project is actively maintained.

1. Check [PLAN.md](PLAN.md) for current status
2. Pick a task or feature
3. Follow standard Go conventions (`gofmt`, `golint`)
4. Write tests for new functionality
5. Submit a PR

## License

[To be specified]

## Acknowledgments

- Built on [Model Context Protocol](https://modelcontextprotocol.io/)
- Uses [Go MCP SDK](https://github.com/modelcontextprotocol/go-sdk)
- LabVIEW integration based on [DQMH framework](https://labviewwiki.org/wiki/Delacor_Queued_Message_Handler)
- Original Python implementation by [CalmyJane](https://github.com/CalmyJane/labview_assistant)

## Support

- **Issues**: [GitHub Issues](https://github.com/MeKo-Christian/labview_mcp/issues)
- **Organization**: [MeKo-Tech](https://github.com/MeKo-Tech)
- **Author**: [@MeKo-Christian](https://github.com/MeKo-Christian)
