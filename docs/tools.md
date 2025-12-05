# LabVIEW MCP Tools Schema Documentation

This document provides a comprehensive schema definition for all tools exposed by the LabVIEW MCP server. This information is extracted from the Python implementation and will serve as the specification for the Go implementation.

## Tool Categories

The 29 tools are organized into the following categories:

1. **Module Management** (2 tools)
2. **VI Creation & Manipulation** (12 tools)
3. **Selection Management** (3 tools)
4. **Debugging & Inspection** (9 tools)
5. **Utility** (3 tools)

---

## Common Output Pattern

All tools follow a consistent output pattern. The LabVIEW VI Call2 method returns a VARIANT array (`param_values`) that contains both inputs and outputs. Common output fields include:

- **error out**: Error information (string)
- **timed out?**: Boolean indicating if the operation timed out
- **result**: String result message
- Various tool-specific outputs

---

## 1. Module Management Tools

### 1.1 start_module

**Purpose**: Initializes the LabVIEW DQMH server module. Must be called before any other LabVIEW-specific functions.

**Inputs**: None

**LabVIEW Parameters** (in param_names/param_values order):
- error out (output)
- Wait for Event Sync? (output)
- Scripting Server Broadcast Events (output)
- Module Was Already Running? (output)
- Module Name (input - empty string)
- Show Main VI Diagram on Init (F) (input - False)
- error in (input - empty string)

**Outputs**:
- Returns param_values array containing module initialization status

**VI Path**: `LabVIEW_Server/Scripting Server/Start Module.vi`

---

### 1.2 stop_module

**Purpose**: Stops the LabVIEW server and cleans up all temporary data. Objects cannot be accessed after stopping.

**Inputs**: None

**LabVIEW Parameters**:
- error out (output)
- Timeout to Wait for Stop (s) (-1: no timeout) (input - -1)
- error in (input - empty string)
- Wait for Module to Stop? (F) (input - True)
- Origin (input - empty string)

**Outputs**:
- Returns param_values array containing stop status

**VI Path**: `LabVIEW_Server/Scripting Server/Stop Module.vi`

---

## 2. VI Creation & Manipulation Tools

### 2.1 new_vi

**Purpose**: Creates a new VI in the LabVIEW IDE and returns a reference for further modifications.

**Inputs**: None

**LabVIEW Parameters**:
- error out (output)
- timed out? (output)
- result (output - string)
- vi_id (output - integer, starts at 0)
- wait for reply (T) (input - True)
- error in (input - empty string)

**Outputs**:
- vi_id: Integer reference to the created VI
- result: String result message
- timed out?: Boolean timeout indicator

**VI Path**: `LabVIEW_Server/Scripting Server/new_vi.vi`

---

### 2.2 add_object

**Purpose**: Adds an object to the block diagram, structure sub-diagram, or front panel.

**Inputs**:
- **position_y**: int - Y coordinate position
- **object_name**: string - Name of the LabVIEW object to add (see Allowed Objects list below)
- **diagram_id**: int - Reference to the diagram or VI
- **position_x**: int - X coordinate position

**LabVIEW Parameters**:
- error out (output)
- timed out? (output)
- result (output - string)
- object_id (output - integer, starts at 0)
- wait for reply (T) (input - True)
- position_y (input - int)
- error in (input - empty string)
- object_name (input - string)
- diagram_id (input - int)
- position_x (input - int)

**Outputs**:
- object_id: Integer reference to the created object
- result: String result message

**Allowed Objects**: Over 500 object types including:
- Structures: While Loop #1, While Loop #2, For Loop, Case Structure, Sequence Structure, Event Structure
- Controls/Indicators: Numeric Control, String Control, Boolean, Array, Cluster, etc.
- Functions: Add, Subtract, Multiply, Divide, Bundle, Unbundle, etc.
- See full list in function docstring at [main.py:95-96](main.py:95-96)

**VI Path**: `LabVIEW_Server/Scripting Server/add_object.vi`

**Notes**: Consider using `create_control` for controls and `enclose_selection` for structures around existing objects.

---

### 2.3 connect_objects

**Purpose**: Connects two terminals of two objects with a wire on the block diagram.

**Inputs**:
- **to_object_terminal_index**: int - Terminal index on the destination object
- **from_object_terminal_index**: int - Terminal index on the source object
- **to_object_reference**: int - Reference ID of destination object
- **from_object_reference**: int - Reference ID of source object
- **vi_reference**: int - Reference ID of the VI

**LabVIEW Parameters**:
- error out (output)
- timed out? (output)
- result (output - string)
- to_object_terminal_index (input - int)
- from_object_terminal_index (input - int)
- wait for reply (T) (input - True)
- to_object_reference (input - int)
- error in (input - empty string)
- from_object_reference (input - int)
- vi_reference (input - int)

**Outputs**:
- result: String result message indicating connection status

**VI Path**: `LabVIEW_Server/Scripting Server/connect_objects.vi`

**Notes**: If a front panel control/indicator reference is passed, it automatically uses the corresponding block diagram terminal. Use `get_object_terminals` to find terminal indices before wiring.

---

### 2.4 save_vi

**Purpose**: Saves a VI to the specified path.

**Inputs**:
- **path**: string - File path including filename (empty string saves to previous location)
- **vi_id**: int - Reference ID of the VI to save

**LabVIEW Parameters**:
- error out (output)
- timed out? (output)
- path_out (output - string)
- wait for reply (T) (input - True)
- error in (input - empty string)
- path (input - string)
- vi_id (input - int)

**Outputs**:
- path_out: String showing where the VI was saved

**VI Path**: `LabVIEW_Server/Scripting Server/save_vi.vi`

---

### 2.5 set_value

**Purpose**: Sets the value of a constant, control, or indicator.

**Inputs**:
- **value**: string - Value to set (can be number as string)
- **object_id**: int - Reference ID of the object

**LabVIEW Parameters**:
- error out (output)
- timed out? (output)
- wait for reply (T) (input - True)
- error in (input - empty string)
- value (input - string)
- object_id (input - int)

**Outputs**:
- Returns param_values with status

**VI Path**: `LabVIEW_Server/Scripting Server/set_value.vi`

---

### 2.6 create_control

**Purpose**: Creates a control, indicator, or constant on a specified terminal.

**Inputs**:
- **constant**: bool - If True, creates a constant (only for inputs)
- **terminal_index**: int - Terminal index to create control on
- **object_id**: int - Reference ID of the node

**LabVIEW Parameters**:
- error out (output)
- timed out? (output)
- created_object_id (output - int, starts at 0)
- wait for reply (T) (input - True)
- constant (input - bool)
- error in (input - empty string)
- terminal_index (input - int)
- object_id (input - int)

**Outputs**:
- created_object_id: Integer reference to the created control/indicator/constant

**VI Path**: `LabVIEW_Server/Scripting Server/create_control.vi`

**Notes**:
- For inputs: creates control (constant=False) or constant (constant=True)
- For outputs: creates indicator
- Returns front panel element reference or block diagram constant reference

---

### 2.7 add_subvi

**Purpose**: Adds a SubVI to a block diagram or structure sub-diagram.

**Inputs**:
- **subvi_path**: string - Path to the VI to add as SubVI
- **diagram_id**: int - Reference ID of the diagram

**LabVIEW Parameters**:
- error out (output)
- timed out? (output)
- subvi_id (output - int, starts at 0)
- wait for reply (T) (input - True)
- error in (input - empty string)
- subvi_path (input - string)
- diagram_id (input - int)

**Outputs**:
- subvi_id: Integer reference to the added SubVI

**VI Path**: `LabVIEW_Server/Scripting Server/add_subvi.vi`

---

### 2.8 open_vi

**Purpose**: Opens an existing VI from the allowed paths and returns a reference.

**Inputs**:
- **vi_path**: string - Path to the VI to open

**LabVIEW Parameters**:
- error out (output)
- timed out? (output)
- vi_id (output - string/int)
- wait for reply (T) (input - True)
- error in (input - empty string)
- vi_path (input - string)

**Outputs**:
- vi_id: Reference to the opened VI

**VI Path**: `LabVIEW_Server/Scripting Server/open_vi.vi`

---

### 2.9 rename_object

**Purpose**: Changes an object's Label.Text property, typically for front panel objects.

**Inputs**:
- **label_visible**: bool - Whether to display the label
- **new_label_name**: string - New label text
- **object_id**: int - Reference ID of the object

**LabVIEW Parameters**:
- error out (output)
- timed out? (output)
- wait for reply (T) (input - True)
- error in (input - empty string)
- label_visible (input - bool)
- new_label_name (input - string)
- object_id (input - int)

**Outputs**:
- Returns param_values with status

**VI Path**: `LabVIEW_Server/Scripting Server/rename_object.vi`

---

### 2.10 delete_object

**Purpose**: Removes an object from the block diagram or front panel.

**Inputs**:
- **object_id**: int - Reference ID of the object to delete

**LabVIEW Parameters**:
- error out (output)
- timed out? (output)
- wait for reply (T) (input - True)
- error in (input - empty string)
- object_id (input - int)

**Outputs**:
- Returns param_values with deletion status

**VI Path**: `LabVIEW_Server/Scripting Server/delete_object.vi`

---

### 2.11 connect_to_pane

**Purpose**: Connects a front panel control/indicator to the SubVI connector pane.

**Inputs**:
- **connector_pane_index**: int - Connector pane position (0-11, see layout below)
- **control_id**: int - Reference ID of the control/indicator

**Connector Pane Layout**:
```
Right side (outputs 0-3): bottom to top
Bottom connectors: lower left (4), upper left (5), lower right (6), upper right (7)
Left side (inputs 8-11): bottom to top
```

**LabVIEW Parameters**:
- error out (output)
- timed out? (output)
- wait for reply (T) (input - True)
- error in (input - empty string)
- connector_pane_index (input - int)
- control_id (input - int)

**Outputs**:
- Returns param_values with connection status

**VI Path**: `LabVIEW_Server/Scripting Server/connect_to_pane.vi`

**Notes**: Controls (VI inputs) should connect to left side (6-11), indicators (VI outputs) to right side (0-5).

---

### 2.12 create_project

**Purpose**: Creates a new LabVIEW project and returns a reference.

**Inputs**: None

**LabVIEW Parameters**:
- error out (output)
- timed out? (output)
- project_id (output - int, starts at 0)
- wait for reply (T) (input - True)
- error in (input - empty string)

**Outputs**:
- project_id: Integer reference to the created project

**VI Path**: `LabVIEW_Server/Scripting Server/create_project.vi`

---

## 3. Selection Management Tools

### 3.1 clear_selection_list

**Purpose**: Clears all elements from the selection list. Use before making a new selection.

**Inputs**:
- **vi_id**: int - Reference ID of the VI

**LabVIEW Parameters**:
- error out (output)
- timed out? (output)
- wait for reply (T) (input - True)
- error in (input - empty string)
- vi_id (input - int)

**Outputs**:
- Returns param_values with clear status

**VI Path**: `LabVIEW_Server/Scripting Server/clear_selection_list.vi`

---

### 3.2 add_to_selection

**Purpose**: Adds an object to the selection list for copy/paste/enclose operations.

**Inputs**:
- **object_id**: int - Reference ID of the object to add
- **vi_id**: int - Reference ID of the VI

**LabVIEW Parameters**:
- error out (output)
- timed out? (output)
- wait for reply (T) (input - True)
- error in (input - empty string)
- object_id (input - int)
- vi_id (input - int)

**Outputs**:
- Returns param_values with updated selection list

**VI Path**: `LabVIEW_Server/Scripting Server/add_to_selection.vi`

---

### 3.3 remove_from_selection

**Purpose**: Removes an object from the selection list.

**Inputs**:
- **object_id**: int - Reference ID of the object to remove
- **vi_id**: int - Reference ID of the VI

**LabVIEW Parameters**:
- error out (output)
- timed out? (output)
- wait for reply (T) (input - True)
- error in (input - empty string)
- object_id (input - int)
- vi_id (input - int)

**Outputs**:
- Returns param_values with updated selection list

**VI Path**: `LabVIEW_Server/Scripting Server/remove_from_selection.vi`

---

## 4. Debugging & Inspection Tools

### 4.1 get_object_terminals

**Purpose**: Returns terminal names/descriptions and their indices for use in connect_objects.

**Inputs**:
- **object_id**: int - Reference ID of the object

**LabVIEW Parameters**:
- error out (output)
- timed out? (output)
- result (output - string with terminal info)
- wait for reply (T) (input - True)
- error in (input - empty string)
- object_id (input - int)

**Outputs**:
- result: String containing terminal names, descriptions, and indices

**VI Path**: `LabVIEW_Server/Scripting Server/get_object_terminals.vi`

---

### 4.2 get_vi_error_list

**Purpose**: Returns the current VI error list (compile errors) in text format.

**Inputs**:
- **vi_reference**: int - Reference ID of the VI

**LabVIEW Parameters**:
- error out (output)
- timed out? (output)
- result (output - string with error list)
- wait for reply (T) (input - True)
- error in (input - empty string)
- vi_reference (input - int)

**Outputs**:
- result: String containing error descriptions and locations

**VI Path**: `LabVIEW_Server/Scripting Server/get_vi_error_list.vi`

**Notes**: Use this to verify if operations succeeded and identify issues.

---

### 4.3 cleanup_vi

**Purpose**: Auto-arranges the block diagram of a VI.

**Inputs**:
- **vi_reference**: int - Reference ID of the VI

**LabVIEW Parameters**:
- error out (output)
- timed out? (output)
- result (output - string)
- wait for reply (T) (input - True)
- error in (input - empty string)
- vi_reference (input - int)

**Outputs**:
- result: String result message

**VI Path**: `LabVIEW_Server/Scripting Server/cleanup_vi.vi`

---

### 4.4 run_vi

**Purpose**: Executes a VI and optionally brings the front panel to foreground.

**Inputs**:
- **open_frontpanel**: bool - Whether to show the front panel
- **vi_id**: int - Reference ID of the VI

**LabVIEW Parameters**:
- error out (output)
- timed out? (output)
- wait for reply (T) (input - True)
- error in (input - empty string)
- open_frontpanel (input - bool)
- vi_id (input - int)

**Outputs**:
- Returns param_values with execution status

**VI Path**: `LabVIEW_Server/Scripting Server/run_vi.vi`

---

### 4.5 get_object_help

**Purpose**: Returns the detailed LabVIEW help article for an object, including terminal descriptions.

**Inputs**:
- **object_id**: int - Reference ID of the object

**LabVIEW Parameters**:
- error out (output)
- timed out? (output)
- help_string (output - string with help text)
- wait for reply (T) (input - True)
- error in (input - empty string)
- object_id (input - int)

**Outputs**:
- help_string: String containing help documentation

**VI Path**: `LabVIEW_Server/Scripting Server/get_object_help.vi`

**Notes**: Use `get_object_terminals` for exact terminal indices.

---

### 4.6 get_loop_conditional_terminal

**Purpose**: Returns the conditional terminal reference of a while or for loop.

**Inputs**:
- **loop_id**: int - Reference ID of the loop

**LabVIEW Parameters**:
- error out (output)
- timed out? (output)
- conditional_terminal_id (output - int, starts at 0)
- wait for reply (T) (input - True)
- error in (input - empty string)
- loop_id (input - int)

**Outputs**:
- conditional_terminal_id: Integer reference to the terminal

**VI Path**: `LabVIEW_Server/Scripting Server/get_loop_conditional_terminal.vi`

**Notes**: Use with `connect_objects` (as to_object_id with index 0) or `create_control`.

---

### 4.7 show_conditional_terminal

**Purpose**: Shows or hides the conditional terminal of a for loop.

**Inputs**:
- **show**: bool - Whether to show the terminal
- **loop_id**: int - Reference ID of the loop

**LabVIEW Parameters**:
- error out (output)
- timed out? (output)
- wait for reply (T) (input - True)
- error in (input - empty string)
- show (input - bool)
- loop_id (input - int)

**Outputs**:
- Returns param_values with status

**VI Path**: `LabVIEW_Server/Scripting Server/show_conditional_terminal.vi`

---

### 4.8 get_loop_iteration_terminal

**Purpose**: Returns the iteration terminal reference of a loop.

**Inputs**:
- **loop_id**: int - Reference ID of the loop

**LabVIEW Parameters**:
- error out (output)
- timed out? (output)
- iteration_terminal_id (output - int, starts at 0)
- wait for reply (T) (input - True)
- error in (input - empty string)
- loop_id (input - int)

**Outputs**:
- iteration_terminal_id: Integer reference to the terminal

**VI Path**: `LabVIEW_Server/Scripting Server/get_loop_iteration_terminal.vi`

---

### 4.9 get_structure_diagram

**Purpose**: Returns the sub-diagram reference for a structure (inside of loops/case structures).

**Inputs**:
- **index**: int - Case index for Case Structures (ignored for loops)
- **structure_id**: int - Reference ID of the structure

**LabVIEW Parameters**:
- error out (output)
- timed out? (output)
- diagram_id (output - int, starts at 0)
- wait for reply (T) (input - True)
- error in (input - empty string)
- index (input - int)
- structure_id (input - int)

**Outputs**:
- diagram_id: Integer reference to the sub-diagram

**VI Path**: `LabVIEW_Server/Scripting Server/get_structure_diagram.vi`

**Notes**: Use this to add objects inside structures. Consider using `enclose_selection` for easier wiring.

---

## 5. Utility Tools

### 5.1 echo

**Purpose**: Simple echo tool for testing MCP connectivity.

**Inputs**:
- **text**: string - Text to echo back

**Outputs**:
- String: "You said: {text}"

**Notes**: This is a Python-only tool, not backed by LabVIEW.

---

### 5.2 enclose_selection

**Purpose**: Creates a structure (loop, case, etc.) around the current selection with automatic tunnels and wiring.

**Inputs**:
- **structure_type**: string - Type of structure to create (see allowed values below)
- **vi_id**: int - Reference ID of the VI

**Allowed Structure Types**:
- While Loop #1
- While Loop #2
- For Loop
- Case Structure
- Sequence Structure
- Event Structure

**LabVIEW Parameters**:
- error out (output)
- timed out? (output)
- object_id (output - int, starts at 0)
- wait for reply (T) (input - True)
- error in (input - empty string)
- structure_type (input - string)
- vi_id (input - int)

**Outputs**:
- object_id: Integer reference to the created structure

**VI Path**: `LabVIEW_Server/Scripting Server/enclose_selection.vi`

**Notes**: This is the most efficient way to create loops/structures. Wire objects first, then use add_to_selection to select them, then call this tool. Tunnels are created automatically.

---

### 5.3 get_allowed_paths

**Purpose**: Returns the directories that the server can write files/folders to.

**Inputs**: None

**LabVIEW Parameters**:
- error out (output)
- timed out? (output)
- Allowed Paths (output - string)
- wait for reply (T) (input - True)
- error in (input - empty string)

**Outputs**:
- Allowed Paths: String with list of allowed directory paths

**VI Path**: `LabVIEW_Server/Scripting Server/get_allowed_paths.vi`

**Notes**: Use this to build full paths when given relative paths. Writes outside these directories will fail.

---

## Implementation Notes for Go Migration

### Common Patterns

1. **All tools return VARIANT arrays** containing both inputs and outputs
2. **"wait for reply (T)"** is always set to True for synchronous operations
3. **Output parameters come first** in param_names/param_values, followed by inputs
4. **Error handling**: "error in" and "error out" are present in all tools
5. **Timeout indicator**: "timed out?" boolean is present in most tools

### Type Mappings for Go

Python Type → Go Type:
- `str` → `string`
- `int` → `int` (or `int32` for COM compatibility)
- `bool` → `bool`
- Return type is always the VARIANT array, which needs parsing

### Parameter Order Pattern

For Go tool implementations, the pattern is:
```go
type ToolNameInput struct {
    Param1 type1 `json:"param1" jsonschema:"description"`
    Param2 type2 `json:"param2" jsonschema:"description"`
}

type ToolNameOutput struct {
    Result string `json:"result" jsonschema:"description"`
    ObjectID int `json:"object_id,omitempty" jsonschema:"description"`
    TimedOut bool `json:"timed_out" jsonschema:"description"`
}
```

### Recommended v1 Tools (Phase A2)

Based on usage patterns and dependencies, recommended tools for initial Go implementation:

**Essential Core** (must have first):
1. `start_module` - Required to initialize the system
2. `new_vi` - Basic VI creation
3. `add_object` - Core object placement

**High Value** (implement next):
4. `connect_objects` - Wiring functionality
5. `get_object_terminals` - Required for correct wiring
6. `save_vi` - Persist work

**Quality of Life** (nice to have):
7. `get_vi_error_list` - Debugging
8. `cleanup_vi` - Code quality

### Critical Implementation Details

1. **VI Paths**: All LabVIEW VIs are in `LabVIEW_Server/Scripting Server/` directory
2. **COM Invocation**: Uses `GetVIReference()` and `Call2()` methods
3. **VARIANT Arrays**: Both param_names and param_values must be COM VARIANT arrays with `VT_BYREF | VT_ARRAY` flags
4. **Synchronous Calls**: The "wait for reply" pattern ensures request-response semantics
5. **Reference Management**: Object IDs are integers starting at 0, managed by the LabVIEW server

---

## Summary Statistics

- **Total Tools**: 29
- **Module Management**: 2
- **VI Creation & Manipulation**: 12
- **Selection Management**: 3
- **Debugging & Inspection**: 9
- **Utility**: 3

**Average Parameters per Tool**: 2-4 inputs, 2-5 outputs

**Most Complex Tools** (by parameter count):
1. `add_object` - 4 inputs + large enum
2. `connect_objects` - 5 inputs
3. `start_module` - 7 LabVIEW parameters

**Simplest Tools**:
1. `echo` - 1 input (Python-only)
2. `stop_module` - 0 inputs
3. `new_vi` - 0 inputs

---

## V1 Implementation Priority (Phase A3)

Based on dependency analysis and usage patterns, tools are prioritized for the initial Go implementation.

### Priority Tier 1: Foundation (MUST HAVE) 🔴

These tools are **absolutely required** for any meaningful LabVIEW scripting:

1. **`start_module`** 🔴 **V1 REQUIRED**
   - **Priority**: P0 - Blocking
   - **Rationale**: Required before ANY other operation. Initializes the DQMH server.
   - **Dependencies**: None
   - **Complexity**: Medium (7 parameters, module lifecycle)
   - **Implementation**: Phase A3, Task 3.2

2. **`new_vi`** 🔴 **V1 REQUIRED**
   - **Priority**: P0 - Blocking
   - **Rationale**: Creates a new VI. Foundation for all VI manipulation.
   - **Dependencies**: Requires `start_module`
   - **Complexity**: Low (0 inputs, returns vi_id)
   - **Implementation**: Phase A3, Task 3.2

3. **`add_object`** 🔴 **V1 REQUIRED**
   - **Priority**: P0 - Blocking
   - **Rationale**: Core functionality - adds nodes, controls, structures to VIs
   - **Dependencies**: Requires `new_vi` or `get_structure_diagram`
   - **Complexity**: High (4 inputs + 500+ object enum)
   - **Implementation**: Phase A3, Task 3.2

### Priority Tier 2: Core Functionality (SHOULD HAVE) 🟡

These tools enable basic VI construction workflows:

4. **`connect_objects`** 🟡 **V1 RECOMMENDED**
   - **Priority**: P1 - High
   - **Rationale**: Wiring is essential for functional VIs
   - **Dependencies**: Requires `add_object` for objects to connect
   - **Complexity**: Medium (5 inputs)
   - **Best Practice**: Use with `get_object_terminals` for correct indices

5. **`get_object_terminals`** 🟡 **V1 RECOMMENDED**
   - **Priority**: P1 - High
   - **Rationale**: Required for correct wiring with `connect_objects`
   - **Dependencies**: Requires objects from `add_object`
   - **Complexity**: Low (1 input, returns string)
   - **Implementation**: Phase A3, Task 3.3

6. **`save_vi`** 🟡 **V1 RECOMMENDED**
   - **Priority**: P1 - High
   - **Rationale**: Persist work, essential for usability
   - **Dependencies**: Requires VI from `new_vi` or `open_vi`
   - **Complexity**: Low (2 inputs)
   - **Security**: Uses path validation

7. **`get_vi_error_list`** 🟡 **V1 RECOMMENDED**
   - **Priority**: P1 - High
   - **Rationale**: Debugging - verify operations succeeded
   - **Dependencies**: Requires VI reference
   - **Complexity**: Low (1 input, returns error list)
   - **Use Case**: Check for compile errors after modifications

8. **`create_control`** 🟡 **V1 RECOMMENDED**
   - **Priority**: P1 - High
   - **Rationale**: Simplifies control/indicator/constant creation
   - **Dependencies**: Requires node from `add_object`
   - **Complexity**: Medium (3 inputs)
   - **Alternative**: Can use `add_object` directly, but less ergonomic

### Priority Tier 3: Quality of Life (NICE TO HAVE) 🟢

These tools improve the development experience but aren't critical:

9. **`cleanup_vi`** 🟢 **V1 OPTIONAL**
   - **Priority**: P2 - Medium
   - **Rationale**: Auto-arrange improves readability
   - **Dependencies**: Requires VI reference
   - **Complexity**: Low (1 input)

10. **`set_value`** 🟢 **V1 OPTIONAL**
    - **Priority**: P2 - Medium
    - **Rationale**: Useful for setting default values
    - **Dependencies**: Requires control/constant reference
    - **Complexity**: Low (2 inputs)

11. **`run_vi`** 🟢 **V1 OPTIONAL**
    - **Priority**: P2 - Medium
    - **Rationale**: Execute and test created VIs
    - **Dependencies**: Requires runnable VI
    - **Complexity**: Low (2 inputs)

12. **`echo`** 🟢 **V1 OPTIONAL**
    - **Priority**: P2 - Medium
    - **Rationale**: Simple connectivity test (Python-only, no LabVIEW)
    - **Dependencies**: None
    - **Complexity**: Trivial (1 input)
    - **Note**: Keep for testing, doesn't require IPC

### Priority Tier 4: Advanced Features (DEFERRED) ⚪

These tools enable advanced workflows but can wait for v2:

13. **Selection Management** (clear, add, remove) ⚪ **V2+**
    - **Priority**: P3 - Low
    - **Rationale**: Used for `enclose_selection`, advanced workflow
    - **Defer To**: V2 when implementing structure creation

14. **`enclose_selection`** ⚪ **V2+**
    - **Priority**: P3 - Low
    - **Rationale**: Powerful but complex, requires selection management
    - **Alternative**: Users can add structures first, then objects inside

15. **Loop/Structure Helpers** ⚪ **V2+**
    - Tools: `get_loop_conditional_terminal`, `show_conditional_terminal`, `get_loop_iteration_terminal`, `get_structure_diagram`
    - **Priority**: P3 - Low
    - **Rationale**: Useful for complex structures, not essential initially

16. **`add_subvi`** ⚪ **V2+**
    - **Priority**: P3 - Low
    - **Rationale**: Requires existing VIs, advanced use case

17. **`open_vi`** ⚪ **V2+**
    - **Priority**: P3 - Low
    - **Rationale**: Editing existing VIs, focus on creation first

18. **`connect_to_pane`** ⚪ **V2+**
    - **Priority**: P3 - Low
    - **Rationale**: SubVI creation, advanced workflow

19. **`rename_object`**, `delete_object` ⚪ **V2+**
    - **Priority**: P3 - Low
    - **Rationale**: Editing operations, less critical than creation

20. **`get_object_help`** ⚪ **V2+**
    - **Priority**: P3 - Low
    - **Rationale**: Documentation lookup, nice but not essential

21. **`get_allowed_paths`** ⚪ **V2+**
    - **Priority**: P3 - Low
    - **Rationale**: Security utility, handle internally initially

22. **`create_project`** ⚪ **V2+**
    - **Priority**: P3 - Low
    - **Rationale**: Project management, advanced feature

23. **`stop_module`** ⚪ **V2+**
    - **Priority**: P4 - Nice to have
    - **Rationale**: Cleanup, but Go can handle gracefully without explicit call
    - **Note**: Implement for completeness, but not critical for basic operation

---

## V1 Implementation Recommendation

### Minimal Viable Product (3 tools)

For the **absolute minimum** v1 to demonstrate end-to-end Go → LabVIEW functionality:

1. 🔴 `start_module` - Initialize system
2. 🔴 `new_vi` - Create VI
3. 🔴 `add_object` - Add content

**Result**: Can create a VI with objects via Go MCP server

### Recommended V1 (8 tools)

For a **usable** v1 that enables basic VI construction:

1. 🔴 `start_module` - Initialize
2. 🔴 `new_vi` - Create VI
3. 🔴 `add_object` - Add nodes/controls
4. 🟡 `get_object_terminals` - Get wiring info
5. 🟡 `connect_objects` - Wire objects
6. 🟡 `save_vi` - Persist work
7. 🟡 `get_vi_error_list` - Debug
8. 🟢 `echo` - Test connectivity (trivial, no LabVIEW)

**Result**: Can create, wire, save, and debug VIs via Go

### Dependency Graph for V1

```
start_module (P0) ────┐
                      ├──> new_vi (P0) ────┐
                      │                     ├──> add_object (P0) ────┬──> get_object_terminals (P1)
                      │                     │                        │
                      │                     │                        └──> connect_objects (P1)
                      │                     │
                      │                     └──> save_vi (P1)
                      │
                      └──> get_vi_error_list (P1)

echo (P2) ────> (standalone, no deps)
```

### Implementation Order for Phase A3

1. **Task 3.2**: Implement `start_module`, `new_vi`, `add_object` (Foundation Tier)
2. **Task 3.3**: Implement `echo` (simple test tool)
3. **Task 3.3**: Implement `get_object_terminals`, `connect_objects`, `save_vi`, `get_vi_error_list` (Core Tier)
4. **Task 3.4**: Build and deploy
5. **Task 3.5**: Verify all 8 tools in Claude Desktop

---

## Rationale Summary

### Why These Tools for V1?

**Foundation Requirements**:
- `start_module`: Mandatory pre-requisite
- `new_vi`: Can't do anything without a VI
- `add_object`: Core value proposition - add LabVIEW objects

**Usability Requirements**:
- `connect_objects`: Unwired objects are useless
- `get_object_terminals`: Required for correct wiring
- `save_vi`: Work persistence is essential
- `get_vi_error_list`: Debugging is critical for LLM iteration

**Testing Requirements**:
- `echo`: Simple connectivity verification

### Why Defer Others?

**Selection Management** (clear/add/remove/enclose):
- Complex workflow pattern
- Can add structures then objects inside (workaround exists)
- Defer to v2 when advanced structure manipulation is needed

**Loop/Structure Helpers**:
- Advanced scenarios
- Users can work around by adding structures first
- Focus v1 on basic VI construction

**File/Project Management** (open_vi, create_project):
- Focus on creation before editing
- Projects are advanced organizational feature

**Advanced Manipulation** (rename, delete, subvi, pane connection):
- Editing operations less critical than creation
- SubVI creation is advanced workflow

**Utility/Inspection** (get_help, get_allowed_paths):
- Nice to have but not blocking
- Can handle internally without exposing to LLM initially
