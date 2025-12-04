# LabVIEW MCP Server (Go Implementation)

A high-performance Model Context Protocol (MCP) server written in Go that enables Large Language Models (LLMs) like Claude to interact with LabVIEW through a clean, well-defined API.

## 🎯 Project Status

**This repository is currently being transformed from Python to Go.** The transformation follows a phased approach documented in [PLAN.md](PLAN.md).

- ✅ **Phase A**: Go MCP server with stubbed tools
- 🚧 **Phase B**: IPC layer and LabVIEW integration
- ⏳ **Phase C**: Production deployment and optimization

See [PLAN.md](PLAN.md) for detailed implementation status and roadmap.

---

## 🚀 Features

- **High Performance**: Go-based MCP server with minimal overhead
- **Clean Architecture**: Modular design with clear separation of concerns
- **Flexible Integration**: Multiple LabVIEW integration options (COM, TCP, or Python adapter)
- **Type Safety**: Strongly-typed tool definitions with comprehensive schema validation
- **Cross-Platform**: Builds for Windows, Linux, and macOS
- **Developer Friendly**: Comprehensive documentation and debugging tools

---

## 📋 Prerequisites

- **Go 1.21+** - For building the MCP server
- **LabVIEW** - With DQMH modules configured
- **Claude Desktop** - Or any MCP-compatible client
- **Windows** - Required for COM integration (or use TCP/HTTP alternative)

---

## 🛠️ Installation

### 1. Clone the Repository

```bash
git clone https://github.com/MeKo-Christian/labview_mcp.git
cd labview_mcp
```

### 2. Build the MCP Server

```bash
# Initialize Go modules
go mod download

# Build for your platform
go build -o bin/mcp-server ./cmd/mcp-server

# Or build for Windows (cross-compile)
GOOS=windows GOARCH=amd64 go build -o bin/mcp-server.exe ./cmd/mcp-server
```

### 3. Build the LabVIEW Bridge (Optional for Phase B+)

```bash
go build -o bin/labview-bridge ./cmd/labview-bridge
```

### 4. Configure Claude Desktop

Add the server to your Claude Desktop configuration file:

**Windows**: `%APPDATA%\Claude\claude_desktop_config.json`

**macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`

**Linux**: `~/.config/Claude/claude_desktop_config.json`

```json
{
  "mcpServers": {
    "labview": {
      "command": "/path/to/labview_mcp/bin/mcp-server.exe",
      "args": []
    }
  }
}
```

### 5. Restart Claude Desktop

Restart Claude Desktop to load the new MCP server. The LabVIEW tools should now appear in Claude's available tools.

---

## 🔧 Architecture

### High-Level Overview

```text
┌─────────────────┐
│  Claude Desktop │
│   (MCP Client)  │
└────────┬────────┘
         │ stdio (JSON-RPC)
         ▼
┌─────────────────────────┐
│   MCP Server (Go)       │
│  - Tool registration    │
│  - Schema validation    │
│  - Request routing      │
└────────┬────────────────┘
         │ IPC (JSON over pipes)
         ▼
┌─────────────────────────┐
│  LabVIEW Bridge (Go)    │
│  - IPC message handling │
│  - LabVIEW integration  │
└────────┬────────────────┘
         │ COM / TCP / Python
         ▼
┌─────────────────────────┐
│   LabVIEW DQMH Module   │
│  - Request VIs          │
│  - Server module        │
└─────────────────────────┘
```

### Directory Structure

```text
labview_mcp/
├── cmd/
│   ├── mcp-server/          # Main MCP server entrypoint
│   │   └── main.go
│   └── labview-bridge/      # LabVIEW integration bridge
│       └── main.go
├── internal/
│   ├── mcpserver/           # MCP server implementation
│   │   ├── server.go
│   │   └── tools.go         # Tool definitions
│   ├── ipc/                 # IPC abstraction layer
│   │   └── ipc.go
│   ├── bridge/              # LabVIEW integration adapters
│   │   ├── python_adapter.go
│   │   ├── labview_com.go   # COM integration (Windows)
│   │   └── labview_client.go # TCP/HTTP client
│   └── logging/             # Structured logging
│       └── logger.go
├── docs/
│   ├── tools.md             # Tool schema documentation
│   ├── ipc-protocol.md      # IPC protocol specification
│   └── architecture.md      # Detailed architecture docs
├── bin/                     # Build artifacts (gitignored)
├── LabVIEW_Server/          # LabVIEW DQMH modules
├── scripts/                 # Python scripts (Phase B3 adapter)
├── main.py                  # Legacy Python server (deprecated)
├── go.mod
├── go.sum
├── PLAN.md                  # Implementation roadmap
└── README.md
```

---

## 🧰 Available Tools

The MCP server exposes LabVIEW functionality through tools. Each tool corresponds to a Request-and-wait-for-reply Event in the DQMH module.

### Core Module Management

#### `start_module`
Start the LabVIEW DQMH server module. **Call this first** before using any other LabVIEW tools.

```json
{
  "tool": "start_module",
  "parameters": {}
}
```

#### `stop_module`
Stop the LabVIEW server module. Cleans up all resources and temporary data.

```json
{
  "tool": "stop_module",
  "parameters": {}
}
```

### VI Creation and Manipulation

#### `new_vi`
Create a new VI in the LabVIEW IDE and return its reference.

```json
{
  "tool": "new_vi",
  "parameters": {}
}
```

**Returns**: `vi_id` (reference to the created VI)

#### `add_object`
Add an object (control, function, structure) to a VI's block diagram or front panel.

```json
{
  "tool": "add_object",
  "parameters": {
    "vi_reference": "12345",
    "object_name": "While Loop"
  }
}
```

**Returns**: `object_id` (reference to the created object)

#### `connect_objects`
Connect two terminals with a wire on the block diagram.

```json
{
  "tool": "connect_objects",
  "parameters": {
    "vi_reference": "12345",
    "from_object_reference": "67890",
    "from_object_terminal_index": 0,
    "to_object_reference": "11111",
    "to_object_terminal_index": 0
  }
}
```

#### `create_control`
Create a control, indicator, or constant on a terminal.

```json
{
  "tool": "create_control",
  "parameters": {
    "object_id": "67890",
    "terminal_index": 0,
    "constant": false
  }
}
```

#### `run_vi`
Execute a VI and optionally open its front panel.

```json
{
  "tool": "run_vi",
  "parameters": {
    "vi_id": "12345",
    "open_frontpanel": true
  }
}
```

### Debugging and Inspection

#### `get_object_terminals`
Get information about an object's terminals (names, indices, descriptions).

```json
{
  "tool": "get_object_terminals",
  "parameters": {
    "object_id": "67890"
  }
}
```

#### `get_vi_error_list`
Get the current error list for a VI (equivalent to clicking the broken run arrow).

```json
{
  "tool": "get_vi_error_list",
  "parameters": {
    "vi_reference": "12345"
  }
}
```

#### `cleanup_vi`
Clean up the block diagram of a VI (auto-arrange objects and wires).

```json
{
  "tool": "cleanup_vi",
  "parameters": {
    "vi_reference": "12345"
  }
}
```

For a complete list of tools and their schemas, see [docs/tools.md](docs/tools.md).

---

## 🔌 LabVIEW Integration Options

The Go bridge supports multiple integration approaches. Choose based on your requirements:

### Option 1: Python Adapter (Quickest)
Reuses existing Python COM integration code. Good for rapid development.

**Pros**: Fast to implement, leverages existing code
**Cons**: Requires Python runtime, additional process overhead

### Option 2: Go COM Integration (Native)
Direct COM integration using `go-ole` library. Best for production Windows deployments.

**Pros**: Native performance, single binary, no Python dependency
**Cons**: Windows-only, requires COM expertise

### Option 3: TCP/HTTP Server VI (Most Flexible)
LabVIEW VI listens on TCP/HTTP and processes JSON requests.

**Pros**: Cross-platform, language-agnostic, network-capable
**Cons**: Requires LabVIEW server VI development, network configuration

See [PLAN.md - Phase B3](PLAN.md#phase-b3--replace-stub-bridge-logic-with-real-labview-calls-later) for implementation details.

---

## 🧪 Development

### Running in Development Mode

```bash
# Run MCP server with debug logging
go run ./cmd/mcp-server --debug

# Test IPC output
go run ./cmd/mcp-server --ipc-test

# Run bridge standalone
go run ./cmd/labview-bridge
```

### Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/ipc/...
```

### Code Generation

The LabVIEW project includes a code generation VI:

1. Open `LabVIEW_Server/AI Assistant for LabVIEW.lvproj`
2. Run `Generate Go Code.vi` (when available)
3. Update tool definitions in `internal/mcpserver/tools.go`
4. Rebuild the MCP server

---

## 📚 Documentation

- [PLAN.md](PLAN.md) - Detailed transformation roadmap with todo-lists
- [docs/tools.md](docs/tools.md) - Complete tool schema reference
- [docs/ipc-protocol.md](docs/ipc-protocol.md) - IPC protocol specification
- [docs/architecture.md](docs/architecture.md) - Architecture deep-dive

---

## 🤝 Contributing

This is an active transformation project. Contributions are welcome!

### Development Workflow

1. Check [PLAN.md](PLAN.md) for current implementation status
2. Pick an uncompleted task or phase
3. Follow the todo-list for that phase
4. Submit a PR with your changes

### Code Style

- Follow standard Go conventions (`gofmt`, `golint`)
- Write tests for new functionality
- Document public APIs with godoc comments
- Update [PLAN.md](PLAN.md) task checkboxes as you progress

---

## 📝 License

[Insert License Here]

---

## 🙏 Acknowledgments

- Built on the [Model Context Protocol](https://modelcontextprotocol.io/) specification
- Uses the [Go MCP SDK](https://github.com/modelcontextprotocol/go-sdk)
- LabVIEW integration based on [DQMH framework](https://labviewwiki.org/wiki/Delacor_Queued_Message_Handler)

---

## 📧 Support

- **Issues**: [GitHub Issues](https://github.com/MeKo-Christian/labview_mcp/issues)
- **Organization**: [MeKo-Tech](https://github.com/MeKo-Tech)
- **Upstream**: Original Python implementation by [CalmyJane](https://github.com/CalmyJane/labview_assistant)

---

## 🗺️ Roadmap

- [x] Phase A1: Skeleton Go MCP server
- [ ] Phase A2: Tool schema extraction
- [ ] Phase A3: Stubbed tool implementation
- [ ] Phase B1: IPC layer implementation
- [ ] Phase B2: Bridge process with dummy responses
- [ ] Phase B3: Real LabVIEW integration
- [ ] Performance optimization
- [ ] Production deployment
- [ ] Deprecate Python implementation

See [PLAN.md](PLAN.md) for detailed progress tracking.
