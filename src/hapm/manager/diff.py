"""HAPM package diff helpers"""
from typing import Literal, Optional

from hapm.package import PackageDescription


class PackageDiff(PackageDescription):
    """Dict describing the Home Assistant package diff"""
    operation: Literal["switch", "delete", "add"]
    current_version: Optional[str]
