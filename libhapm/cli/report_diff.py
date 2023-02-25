"""HAPM CLI diff reporter"""
from __future__ import annotations

from typing import List

from libhapm.manager.diff import PackageDiff
from libhapm.package import repo_name

from .ink import ANSI_DIM, ANSI_GREEN, ANSI_RED, ANSI_YELLOW, ink
from .utils import group_by_kind


def _format_kind(kind: str) -> str:
    return ink(kind.capitalize() + ":", effects=ANSI_DIM)


def _format_entry(diff: PackageDiff) -> str:
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
    name = repo_name(diff["full_name"])
    title = ink(f"{prefix} {name}", color=color)
    return f"{title} {version}"


def report_diff(diff: List[PackageDiff]):
    """Prints in stdout diff of packages in a nice way"""
    groups = group_by_kind(diff)
    log = ""
    for kind, packages in groups.items():
        log += f"{_format_kind(kind)}\n"
        for package in packages:
            log += f"{_format_entry(package)}\n"
    print(log, end="\r")
