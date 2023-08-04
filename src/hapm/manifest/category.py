"""HAPM manifest parsing utils"""
from __future__ import annotations

from typing import Dict, List

from hapm.package import PackageDescription

from .location import parse_location


def parse_category(manifest: Dict[str, List[str]], key: str) -> List[PackageDescription]:
    """Parses the manifest, turning it into a list of packages"""
    if key not in manifest:
        raise TypeError(f"Key {key} is not found in repo")
    items: List[PackageDescription] = []
    for entry in manifest[key]:
        location = parse_location(entry)
        if location["full_name"] is None:
            raise TypeError(f"Wrong entity: {entry}")
        if location["version"] == "latest":
            raise TypeError(f"Version is missing: {entry}")
        items.append({
            "full_name": location["full_name"],
            "version": location["version"],
            "kind": key
        })
    return items
