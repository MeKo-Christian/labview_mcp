# Design: Eliminate labview-bridge Binary

**Date:** 2025-12-05
**Status:** Approved
**Author:** Claude + MeKo-Christian

## Problem

The current architecture has three layers:
```
mcp-server (Go) → labview-bridge (Go) → Python wrappers
```

The `labview-bridge` is a thin wrapper that only:
- Detects platform (Windows vs non-Windows)
- Calls Python scripts via `exec.Command`
- Returns canned responses on non-Windows

This adds deployment complexity (two binaries), IPC overhead, and maintenance burden with minimal benefit.

## Solution

Merge bridge functionality into mcp-server by adding direct Python execution:

```
mcp-server (Go) → Python wrappers (on Windows)
```

## Architecture

### New Component: PythonExecutor

**Location:** `internal/python/executor.go`

**Interface:**
```go
type Executor struct {
    wrapperDir string
    isWindows  bool
}

func NewExecutor() *Executor
func (e *Executor) Call(ctx context.Context, tool string, input interface{}) (json.RawMessage, error)
```

**Responsibilities:**
- Platform detection using `runtime.GOOS`
- Python script path resolution (relative to binary location)
- Execute Python with JSON stdin/stdout
- Return canned responses on non-Windows
- Error handling and logging

### Modified Components

**1. internal/mcpserver/tools.go**
- Replace `callIPC()` with `callPython()`
- Remove global `globalIPCClient`, add global `globalPythonExecutor`

**2. internal/mcpserver/server.go**
- Replace `IPCClient *ipc.Client` with `PythonExecutor *python.Executor`
- Remove IPC client initialization

**3. cmd/mcp-server/main.go**
- Remove bridge spawning logic (lines 56-100)
- Remove IPC client creation
- Create PythonExecutor and pass to server

### Deleted Components

- `cmd/labview-bridge/` - entire directory
- `internal/ipc/` - entire package

## Implementation Details

### Python Execution Flow

**On Windows:**
1. Resolve wrapper script path: `<binary_dir>/../wrapper/<tool>.py`
2. Marshal input to JSON
3. Execute: `python <script_path>` with JSON stdin
4. Read JSON from stdout
5. Unmarshal to result type
6. Log stderr if present

**On non-Windows:**
- Return tool-specific canned response immediately

### Canned Responses

Moved from bridge to PythonExecutor with same structure:
- `start_module` → status: "stub: module started"
- `new_vi` → vi_id: 42, result: "stub: VI created"
- etc.

### Error Handling

**Script not found:**
```
Error: wrapper script not found: /path/to/wrapper/tool_name.py
```

**Python execution failed:**
```
Error: python execution failed: <error> (stderr: <stderr_output>)
```

**JSON parse failed:**
```
Error: failed to parse python output: <error> (output: <raw_output>)
```

## Benefits

1. **Simpler deployment** - One binary instead of two
2. **Reduced latency** - No IPC overhead between Go processes
3. **Less code** - Remove ~300 lines (bridge + IPC protocol)
4. **Easier debugging** - Direct execution, clearer error messages
5. **Same functionality** - Python wrappers unchanged

## Migration

**No user migration needed:**
- Configuration unchanged (still points to mcp-server binary)
- Python wrapper scripts unchanged
- Tool interface unchanged

**Deployment changes:**
- Remove `labview-bridge` binary from deployment
- Single binary: `mcp-server` (or `mcp-server.exe` on Windows)

## Testing

**Test scenarios:**
1. Non-Windows platform → verify stub responses
2. Windows without Python → verify clear error messages
3. Windows with Python → verify tool execution
4. Script not found → verify error handling
5. Invalid JSON from Python → verify error handling

## Rollback Plan

If issues arise:
1. Revert to commit before this change
2. Rebuild both binaries
3. Redeploy bridge-based architecture

Git history preserves all bridge code for reference.
