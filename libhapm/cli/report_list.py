"""HAPM CLI list reporter"""
from __future__ import annotations

from typing import List

from libhapm.package import PackageDescription, repo_name

from .ink import ANSI_DIM, ink
from .utils import group_by_kind


def _format_kind(kind: str) -> str:
    return ink(kind.capitalize() + ":", effects=ANSI_DIM)


def _format_entry(package: PackageDescription) -> str:
    name = repo_name(package["full_name"])
    version = ink(package["version"], effects=ANSI_DIM)
    return f"  {name} {version}"


def report_list(diff: List[PackageDescription]):
    """Prints into stdout list of packages in a nice way"""
    groups = group_by_kind(diff)
    log = ""
    for kind, packages in groups.items():
        log += f"{_format_kind(kind)}\n"
        for package in packages:
            log += f"{_format_entry(package)}\n"
    print(log, end="\r")
