# LabVIEW VI to Python Tool Mapping

This document maps each Python MCP tool to its corresponding LabVIEW VI file and describes the LabVIEW server architecture.

## Architecture Overview

The LabVIEW server uses the **DQMH (Delacor Queued Message Handler)** architecture pattern:

- **Request VIs**: Public API VIs that are called from Python (e.g., `new_vi.vi`, `add_object.vi`)
- **Do VIs**: Internal implementation VIs that perform the actual work (e.g., `Do Create Control.vi`)
- **Module VIs**: Lifecycle management (`Start Module.vi`, `Stop Module.vi`)
- **Helper VIs**: Utilities for reference management, error handling, etc.

### Directory Structure

```
LabVIEW_Server/
├── Scripting Server/          # Main DQMH module
│   ├── Start Module.vi        # Module initialization
│   ├── Stop Module.vi         # Module shutdown
│   ├── <tool_name>.vi         # Public API VIs (called from Python)
│   ├── Do <Action>.vi         # Internal implementation VIs
│   ├── Get <Resource>.vi      # Helper VIs for reference management
│   └── Test_*.vi              # Test VIs
└── Tools/                     # Code generation utilities
    └── PyServer/              # Python MCP server generation tools
```

---

## Python Tool → LabVIEW VI Mapping

### 1. Module Management

| Python Tool | LabVIEW VI | Internal Implementation | Notes |
|------------|-----------|------------------------|-------|
| `start_module()` | `Start Module.vi` | `Init Module.vi` | Initializes DQMH module, creates event queues |
| `stop_module()` | `Stop Module.vi` | `Close Module.vi`, `Handle Exit.vi` | Cleans up module, destroys events |

---

### 2. VI Creation & Manipulation

| Python Tool | LabVIEW VI | Internal Implementation | Notes |
|------------|-----------|------------------------|-------|
| `new_vi()` | `new_vi.vi` | `Create New VI.vi`, `Add New Reference.vi` | Creates VI and stores reference |
| `add_object()` | `add_object.vi` | `add vi object.vi` | Adds object to diagram/panel |
| `connect_objects()` | `connect_objects.vi` | `Wire Two Objects.vi` | Wires terminals together |
| `save_vi()` | `save_vi.vi` | `Do Save VI.vi`, `Is Allowed Path.vi` | Saves VI with path validation |
| `set_value()` | `set_value.vi` | `Do Set Value.vi`, `Get Control by ID.vi` | Sets control/constant value |
| `create_control()` | `create_control.vi` | `Do Create Control.vi` | Creates control/indicator/constant |
| `add_subvi()` | `add_subvi.vi` | `Do Add Subvi.vi` | Adds SubVI to diagram |
| `open_vi()` | `open_vi.vi` | `Do Open Vi.vi`, `Load VI.vi` | Opens existing VI from path |
| `rename_object()` | `rename_object.vi` | `Do Rename Object.vi` | Changes object label |
| `delete_object()` | `delete_object.vi` | `Do Delete Object.vi` | Removes object |
| `connect_to_pane()` | `connect_to_pane.vi` | `do connect to pane.vi` | Connects to connector pane |
| `create_project()` | `create_project.vi` | `Do create project.vi` | Creates new project |

---

### 3. Selection Management

| Python Tool | LabVIEW VI | Internal Implementation | Notes |
|------------|-----------|------------------------|-------|
| `clear_selection_list()` | `clear_selection_list.vi` | `Do Clear Selection List.vi` | Clears selection |
| `add_to_selection()` | `add_to_selection.vi` | `Do Add To Selection.vi` | Adds to selection |
| `remove_from_selection()` | `remove_from_selection.vi` | `Do Remove From Selection.vi` | Removes from selection |

---

### 4. Debugging & Inspection

| Python Tool | LabVIEW VI | Internal Implementation | Notes |
|------------|-----------|------------------------|-------|
| `get_object_terminals()` | `get_object_terminals.vi` | `get terminal info.vi`, `Get Terminals.vi` | Returns terminal info |
| `get_vi_error_list()` | `get_vi_error_list.vi` | `generate vi error list.vi` | Returns compile errors |
| `cleanup_vi()` | `cleanup_vi.vi` | `clean up vis blockdiagram.vi` | Auto-arranges diagram |
| `run_vi()` | `run_vi.vi` | `Do Run VI.vi` | Executes VI |
| `get_object_help()` | `get_object_help.vi` | `Do Get Object Help.vi` | Returns help text |
| `get_loop_conditional_terminal()` | `get_loop_conditional_terminal.vi` | `Do Get Cond Terminal.vi` | Gets loop conditional |
| `show_conditional_terminal()` | `show_conditional_terminal.vi` | `do show conditional terminal.vi` | Shows/hides terminal |
| `get_loop_iteration_terminal()` | `get_loop_iteration_terminal.vi` | `do get loop iteration terminal.vi` | Gets iteration terminal |
| `get_structure_diagram()` | `get_structure_diagram.vi` | `Do Get Structure Diagram.vi` | Gets sub-diagram |

---

### 5. Utility

| Python Tool | LabVIEW VI | Internal Implementation | Notes |
|------------|-----------|------------------------|-------|
| `echo()` | N/A | N/A | Python-only test tool |
| `enclose_selection()` | `enclose_selection.vi` | `Do Enclose Selection.vi` | Creates structure around selection |
| `get_allowed_paths()` | `get_allowed_paths.vi` | `Set Allowed Paths.vi`, `Is Allowed Path.vi` | Returns allowed directories |

---

## Helper VIs (Internal Use)

These VIs are not directly called from Python but are essential infrastructure:

### Reference Management
- `Add New Reference.vi` - Registers new object references
- `Get VI By ID.vi` - Retrieves VI reference by ID
- `Get GObject by ID.vi` - Retrieves graphic object by ID
- `Get Node by ID.vi` - Retrieves node reference by ID
- `Get Control by ID.vi` - Retrieves control reference by ID
- `Get Terminal by ID.vi` - Retrieves terminal reference by ID
- `Get Diagram By ID.vi` - Retrieves diagram reference by ID

### Module Management
- `Init Module.vi` - Module initialization
- `Module Did Init.vi` - Initialization complete handler
- `Module Did Stop.vi` - Stop complete handler
- `Get Module Execution Status.vi` - Status query
- `Update Module Execution Status.vi` - Status update
- `Get Module Main VI Information.vi` - Module info

### Synchronization
- `Obtain Request Events.vi` - Gets request event refnums
- `Obtain Broadcast Events.vi` - Gets broadcast event refnums
- `Obtain Broadcast Events for Registration.vi` - For event registration
- `Destroy Request Events.vi` - Cleanup request events
- `Destroy Broadcast Events.vi` - Cleanup broadcast events
- `Destroy Sync Refnums.vi` - Cleanup sync refnums
- `Get Sync Refnums.vi` - Gets synchronization refnums
- `Wait on Event Sync.vi` - Wait for event
- `Wait on Module Sync.vi` - Wait for module
- `Synchronize Caller Events.vi` - Sync caller
- `Synchronize Module Events.vi` - Sync module

### Semaphore Management
- `Obtain Module Semaphore.vi` - Gets semaphore
- `Acquire Module Semaphore.vi` - Acquires lock
- `Release Module Semaphore.vi` - Releases lock

### Error Handling
- `Error Reported.vi` - Error reporting
- `Module Not Running--error.vi` - Error constant
- `Module Not Stopped--error.vi` - Error constant
- `Module Not Synced--error.vi` - Error constant
- `Request and Wait for Reply Timeout--error.vi` - Timeout error
- `generate forbidden path error.vi` - Path validation error

### Path Validation
- `Is Allowed Path.vi` - Validates file paths
- `Set Allowed Paths.vi` - Configures allowed paths

### UI Management
- `Show Panel.vi` - Shows front panel
- `Hide Panel.vi` - Hides front panel
- `Open VI Panel.vi` - Opens VI panel
- `Hide VI Panel.vi` - Hides VI panel
- `Show Diagram.vi` - Shows block diagram

### Utilities
- `Get Datatype.vi` - Gets LabVIEW datatype
- `Is Same Object.vi` - Compares object references
- `Load VI.vi` - Loads VI into memory
- `Read Config Data.vi` - Reads configuration
- `Write INI.vi` - Writes INI file
- `Status Updated.vi` - Status update handler

### Constants
- `Module Name--constant.vi` - Module name
- `Module Timeout--constant.vi` - Timeout value
- `Null Broadcast Events--constant.vi` - Null events constant

---

## Test VIs

The following test VIs demonstrate how to use the API:

- `Test Scripting Server API.vi` - Main test VI
- `Test_Create_Control.vi` - Control creation test
- `Test_Connect_Connector_Pane.vi` - Connector pane test
- `Test_EncloseSelection.vi` - Enclose selection test
- `Test_RandomNumberGenerator.vi` - Example: creates random number generator
- `Test_Rename_And_Set_Value.vi` - Rename/value test
- `Test_StructureSubdiagramsAndDeleteObject.vi` - Structure test
- `Test_WireLoops.vi` - Loop wiring test

---

## DQMH Request-Reply Pattern

The module uses DQMH's **Request and Wait for Reply** pattern:

1. **Python calls Request VI** (e.g., `new_vi.vi`)
2. **Request VI uses Call2** with `wait for reply (T) = True`
3. **Internal Do VI executes** the actual work
4. **Response is returned** via the param_values array
5. **Python receives result** and extracts outputs

### Call2 Parameters

All Request VIs use the same Call2 invocation pattern:

```labview
vi.Call2(
    param_names,      // VARIANT array of parameter names
    param_values,     // VARIANT array of parameter values (modified in place)
    False,            // open FP?
    False,            // close FP after call?
    False,            // suspend on call?
    False             // bring LabVIEW to front?
)
```

### Parameter Order Convention

**Output parameters first, then input parameters:**

```
param_names = [
    "error out",           // OUTPUT
    "timed out?",          // OUTPUT
    "result",              // OUTPUT
    "object_id",           // OUTPUT
    "wait for reply (T)",  // INPUT (always True)
    "input_param_1",       // INPUT
    "error in",            // INPUT (always empty)
    "input_param_2",       // INPUT
]
```

---

## Code Generation Tools

The `Tools/` directory contains VIs for generating the Python MCP server:

- `Update MCP Server.vi` - Main code generation VI
- `Add Tool From VI.vi` - Adds tool from VI specification
- `Get Python Code.vi` - Generates Python code
- `Get VI Function String.vi` - Gets function signature
- `Get VI In and Outputs.vi` - Extracts VI parameters
- `Create VI function name.vi` - Creates function name
- `Get Parameter String.vi` - Gets parameter string
- Various type checking utilities (`Is Int.vi`, `Is String.vi`, `Is Float.vi`)

### PyServer Subpackage

- `New.vi` - Creates new server skeleton
- `Add Tool.vi` - Adds tool to server
- `Create Tool String.vi` - Generates tool code
- `Get Python.vi` - Gets Python template

---

## Verification Status

All 29 Python tools have corresponding LabVIEW VI files:

✅ **All mapping verified** - Every Python tool has a matching `.vi` file
✅ **DQMH architecture confirmed** - Request/Reply pattern in use
✅ **Helper infrastructure identified** - Reference management, sync, semaphores
✅ **Test coverage exists** - Multiple test VIs demonstrate usage

---

## Notes for Go Implementation

### Reference IDs
- Object IDs are managed server-side as integers starting from 0
- References are stored in memory and persist until module stop
- Each type has its own helper (VI, GObject, Node, Control, Terminal, Diagram)

### Path Security
- All file operations go through `Is Allowed Path.vi` validation
- Attempts to write outside allowed paths generate errors
- Use `get_allowed_paths()` to query valid directories

### Synchronization
- The DQMH pattern ensures thread-safe operations
- "wait for reply (T)" must be True for synchronous calls
- Module must be started before any other operations

### Error Handling
- "error in" and "error out" present in all VIs
- Errors propagate through the param_values array
- Check "timed out?" boolean for timeout conditions

### For Go IPC Bridge (Phase B)
When implementing the Go bridge that calls these VIs:
1. Start with `Start Module.vi` call
2. Use COM or subprocess to invoke Request VIs
3. Parse param_values array for outputs
4. Map integer IDs to Go-side references if needed
5. Call `Stop Module.vi` on cleanup
