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
