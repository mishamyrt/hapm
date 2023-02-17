"""HAPM manifest controller"""
from typing import Dict, List


from ruamel.yaml import safe_load

from ..package import Package
from .parse import parse_category
from .types import ManifestDict


class Manifest:
    """HAPM manifest controller"""
    _raw: ManifestDict

    values: Dict[str, List[Package]] = {}
    _lovelace: List[Package]
    _themes: List[Package]

    def __init__(self, path: str):
        self._path = path
        with open(path, "r", encoding="utf-8") as stream:
            raw = safe_load(stream)
        for key in raw:
            self.values[key] = parse_category(raw, key)
