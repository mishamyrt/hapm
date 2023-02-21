"""HAPM manifest controller"""
from typing import Dict, List

from ruamel.yaml import safe_load

from libhapm.package import PackageDescription

from .parse import parse_category
from .types import ManifestDict


class Manifest:
    """HAPM manifest controller"""
    _encoding: str

    values: List[PackageDescription] = []

    def __init__(self, path: str, encoding="utf-8"):
        self.path = path
        self._encoding = encoding
    
    def load(self) -> str:
        with open(self.path, "r", encoding="utf-8") as stream:
            raw = safe_load(stream)
        for key in raw:
            self.values.extend(parse_category(raw, key))
