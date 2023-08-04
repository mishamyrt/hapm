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

def report_wrong_format(location: str):
    report_error(f"Wrong location format: '{location}'")
    example = """Package Location can be specified in several formats.
* Root or tag URL of a repository on GitHub.
  - https://github.com/mishamyrt/myrt_desk_hass
  - https://github.com/mishamyrt/myrt_desk_hass/releases/tag/v0.2.4
* Package name with version separated by the @ symbol
  - mishamyrt/myrt_desk_hass
  - mishamyrt/myrt_desk_hass@v0.2.4
If no version is specified, then latest will be used."""
    report_warning(example)


def report_exception(action: str, exception: Exception):
    report_error(f"Error while {action}: {str(exception)}")

def report_warning(text: str | int):
    print(ink(text, ANSI_YELLOW))

def report_error(text: str | int):
    print(ink(text, ANSI_RED))
