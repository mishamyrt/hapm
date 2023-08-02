"""HAPM CLI error reporter"""
from .ink import ANSI_RED, ink


def report_error(action: str, e: Exception):
    message = f"Error while {action}:\n"
    message += ink(e, ANSI_RED)
    print(message)
