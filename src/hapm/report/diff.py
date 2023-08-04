"""HAPM CLI diff reporter"""
from __future__ import annotations

from typing import List

from hapm.color import ANSI_DIM, ANSI_GREEN, ANSI_RED, ANSI_YELLOW, ink
from hapm.manager.diff import PackageDiff
from hapm.repository import repo_name

from .utils import group_by_kind


def _format_kind(kind: str) -> str:
    return ink(kind.capitalize() + ":", effects=ANSI_DIM)

def _format_update(diff: PackageDiff) -> str:
    full_name = diff["full_name"]
    next_version = ink(diff["version"], ANSI_YELLOW)
    current_version = ink(f"({diff['current_version']})", effects=ANSI_DIM)
    delimiter = ink("@", effects=ANSI_DIM)
    return f"{full_name}{delimiter}{next_version} {current_version}"

def _format_entry(diff: PackageDiff, full_name: bool) -> str:
    version = diff["version"]
    if diff["operation"] == "add":
        prefix = "+"
        color = ANSI_GREEN
        version = ink(version, effects=ANSI_DIM)
    elif diff["operation"] == "switch":
        prefix = "*"
        color = ANSI_YELLOW
        version = f"{ink(diff['current_version'], effects=ANSI_DIM)} → {version}"
    else:
        prefix = "-"
        color = ANSI_RED
        version = ink(version, effects=ANSI_DIM)
    if full_name:
        name = diff["full_name"]
    else:
        name = repo_name(diff["full_name"])
    title = ink(f"{prefix} {name}", color=color)
    return f"{title}@{version}"


def report_diff(diff: List[PackageDiff], full_name=False, updates_only=False):
    """Prints in stdout diff of packages in a nice way"""
    groups = group_by_kind(diff)
    log = ""
    for kind, packages in groups.items():
        log += f"{_format_kind(kind)}\n"
        for package in packages:
            if updates_only:
                log += f"  {_format_update(package)}\n"
            else:
                log += f"{_format_entry(package, full_name)}\n"
    print(log, end="\r")
