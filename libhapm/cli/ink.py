from __future__ import annotations

from typing import List

ANSI_RESET = "0"
ANSI_DIM = "2"
ANSI_RED = "31"
ANSI_GREEN = "32"
ANSI_YELLOW = "33"
ANSI_B_CYAN = "96"


def ink(text: str | int, color=ANSI_RESET, effects: List[str] = None):
    prefix = f"\033[{color}"
    if effects is not None:
        prefix += ";" + ";".join(effects)
    return f"{prefix}m{str(text)}\033[{ANSI_RESET}m"
