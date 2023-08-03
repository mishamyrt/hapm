"""HAPM coloring util"""
from __future__ import annotations

from typing import List

from .constants import ANSI_RESET


def ink(text: str | int, color=ANSI_RESET, effects: List[str] = None) -> str:
    """Colors the text to be output to the console"""
    prefix = f"\033[{color}"
    if effects is not None:
        prefix += ";" + ";".join(effects)
    return f"{prefix}m{str(text)}\033[{ANSI_RESET}m"
