#!/usr/bin/env python3
"""
Python wrapper for start_module LabVIEW tool.
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

        # Connect to LabVIEW
        lv_app = win32com.client.Dispatch("LabVIEW.Application")

        # Get VI path
        script_dir = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
        vi_path = os.path.join(script_dir, "LabVIEW_Server", "Scripting Server", "Start Module.vi")

        # Get VI reference
        vi = lv_app.GetVIReference(vi_path, "", False, 0)
        vi._FlagAsMethod("Call2")

        # Prepare parameters
        param_names = VARIANT(
            pythoncom.VT_BYREF | pythoncom.VT_ARRAY | pythoncom.VT_BSTR,
            ("error out", "Wait for Event Sync?", "Scripting Server Broadcast Events",
             "Module Was Already Running?", "Module Name", "Show Main VI Diagram on Init (F)", "error in")
        )

        param_values = VARIANT(
            pythoncom.VT_BYREF | pythoncom.VT_ARRAY | pythoncom.VT_VARIANT,
            ("", False, "", False, "", False, "")
        )

        # Call LabVIEW VI
        vi.Call2(param_names, param_values, False, False, False, False)

        # Extract results
        error_out = str(param_values[0]) if param_values[0] else ""
        module_already_running = bool(param_values[3])

        # Return JSON response
        result = {
            "status": "module started" if not error_out else "error",
            "module_already_running": module_already_running,
            "error_out": error_out
        }

        json.dump(result, sys.stdout)
        sys.stdout.flush()

    except Exception as e:
        # Return error as JSON
        error_result = {
            "status": "error",
            "error_out": str(e)
        }
        json.dump(error_result, sys.stdout)
        sys.stdout.flush()
        sys.exit(1)


if __name__ == "__main__":
    main()
