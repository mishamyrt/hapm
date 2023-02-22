"""HAPM CLI summary reporter"""

from typing import List

from libhapm.manager.diff import PackageDiff

from .ink import ANSI_B_CYAN, ink


def report_summary(diff: List[PackageDiff]):
    """Prints the result of the successful work into the stdout"""
    if len(diff) == 0:
        print("There's nothing to do here")
        return
    adds = 0
    deletes = 0
    switches = 0
    for package in diff:
        operation = package["operation"]
        if operation == "add":
            adds += 1
        elif operation == "delete":
            deletes += 1
        elif operation == "switch":
            switches += 1
    results = []
    if adds > 0:
        results.append(f"installed {ink(adds, ANSI_B_CYAN)}")
    if deletes > 0:
        results.append(f"removed {ink(deletes, ANSI_B_CYAN)}")
    if switches > 0:
        results.append(f"switched {ink(switches, ANSI_B_CYAN)}")
    print(f"\nDone: {', '.join(results)}")
