#!/usr/bin/env python3
"""
Python wrapper for connect_objects LabVIEW tool.
Reads JSON from stdin, calls LabVIEW via COM, writes JSON to stdout.
"""
import sys
import json
import os
import pythoncom
from win32com.client import VARIANT
import win32com.client


def main():
    try:
        # Read JSON input from stdin
        input_data = json.load(sys.stdin)

        # Extract input parameters
        to_object_terminal_index = input_data.get("to_object_terminal_index", 0)
        from_object_terminal_index = input_data.get("from_object_terminal_index", 0)
        to_object_reference = input_data.get("to_object_reference", 0)
        from_object_reference = input_data.get("from_object_reference", 0)
        vi_reference = input_data.get("vi_reference", 0)

        # Connect to LabVIEW
        lv_app = win32com.client.Dispatch("LabVIEW.Application")

        # Get VI path
        script_dir = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
        vi_path = os.path.join(script_dir, "LabVIEW_Server", "Scripting Server", "connect_objects.vi")

        # Get VI reference
        vi = lv_app.GetVIReference(vi_path, "", False, 0)
        vi._FlagAsMethod("Call2")

        # Prepare parameters
        param_names = VARIANT(
            pythoncom.VT_BYREF | pythoncom.VT_ARRAY | pythoncom.VT_BSTR,
            ("error out", "timed out?", "result", "to_object_terminal_index",
             "from_object_terminal_index", "wait for reply (T)", "to_object_reference",
             "error in", "from_object_reference", "vi_reference")
        )

        param_values = VARIANT(
            pythoncom.VT_BYREF | pythoncom.VT_ARRAY | pythoncom.VT_VARIANT,
            ("", False, "", to_object_terminal_index, from_object_terminal_index,
             True, to_object_reference, "", from_object_reference, vi_reference)
        )

        # Call LabVIEW VI
        vi.Call2(param_names, param_values, False, False, False, False)

        # Extract results
        error_out = str(param_values[0]) if param_values[0] else ""
        timed_out = bool(param_values[1])
        result = str(param_values[2]) if param_values[2] else ""

        # Return JSON response
        output = {
            "result": result,
            "timed_out": timed_out,
            "error_out": error_out
        }

        json.dump(output, sys.stdout)
        sys.stdout.flush()

    except Exception as e:
        # Return error as JSON
        error_result = {
            "result": "",
            "timed_out": False,
            "error_out": str(e)
        }
        json.dump(error_result, sys.stdout)
        sys.stdout.flush()
        sys.exit(1)


if __name__ == "__main__":
    main()
