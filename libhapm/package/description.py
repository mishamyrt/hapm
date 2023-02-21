"""HAPM package description helpers"""
from typing import TypedDict


class PackageDescription(TypedDict):
    """Dict describing the Home Assistant package"""
    url: str
    version: str
    kind: str
