"""HAPM CLI error reporter"""
from __future__ import annotations

from hapm.color import ANSI_RED, ANSI_YELLOW, ink

TOKEN_GENERATE_LINK = "https://github.com/settings/tokens"

def report_no_token(env: str):
    """Prints to stdout that the user needs to set a variable"""
    message = f"""${env} is not defined.
Open {TOKEN_GENERATE_LINK},
generate a personal token and set it in the ${env} variable.
Otherwise you will run into rate limit fairly quickly."""
    report_warning(message)

def report_error(action: str, exception: Exception):
    message = f"Error while {action}:\n"
    message += ink(exception, ANSI_RED)
    print(message)

def report_warning(text: str | int):
    print(ink(text, ANSI_YELLOW))
