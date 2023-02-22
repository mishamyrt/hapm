"""HAPM CLI common utils"""
from __future__ import annotations

from typing import Dict, List

from libhapm.package import PackageDescription


def group_by_kind(packages: List[PackageDescription]) -> Dict[str, List[PackageDescription]]:
    """Groups packages by kind. Returns a dictionary,
    where the key is the kind and the value is an array of packages of the appropriate kind"""
    groups: Dict[str, List[PackageDescription]] = {}
    for package in packages:
        kind = package['kind']
        if kind in groups:
            groups[kind].append(package)
        else:
            groups[kind] = [package]
    return groups
