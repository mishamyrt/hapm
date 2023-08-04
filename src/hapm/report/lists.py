"""HAPM CLI lists reporter"""
from __future__ import annotations

from typing import List

from hapm.color import ANSI_DIM, ANSI_YELLOW, ink
from hapm.git import repo_name
from hapm.package import PackageDescription

from .utils import group_by_kind


def _format_kind(kind: str) -> str:
    return ink(kind.capitalize() + ":", effects=ANSI_DIM)


def _format_package(package: PackageDescription) -> str:
    name = repo_name(package["full_name"])
    version = ink(package["version"], effects=ANSI_DIM)
    return f"  {name} {version}"

def _format_version(package: str, version: str) -> str:
    line = ink(f"- {package}@", effects=ANSI_DIM)
    line += ink(version, ANSI_YELLOW)
    return line


def report_packages(diff: List[PackageDescription]) -> None:
    """Prints into stdout list of packages in a nice way"""
    groups = group_by_kind(diff)
    log = ""
    for kind, packages in groups.items():
        log += f"{_format_kind(kind)}\n"
        for package in packages:
            log += f"{_format_package(package)}\n"
    print(log, end="\r")

def report_versions(package: str, versions: List[str]) -> None:
    """Prints package versions in a format suitable for installation"""
    content = ""
    for version in versions:
        content += _format_version(package, version) + "\n"
    print(content)
