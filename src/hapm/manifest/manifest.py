"""HAPM manifest controller"""
from typing import List

from ruamel.yaml import safe_load

from hapm.package import PackageDescription
from hapm.report import report_exception

from .category import parse_category


class Manifest:
    """HAPM manifest controller"""
    _encoding: str

    values: List[PackageDescription] = []

    def __init__(self, path: str, encoding="utf-8"):
        self.path = path
        self._encoding = encoding

    def load(self) -> List[PackageDescription]:
        """Reads the manifest file and parses its contents"""
        with open(self.path, "r", encoding="utf-8") as stream:
            raw = safe_load(stream)
        try:
            for key in raw:
                self.values.extend(parse_category(raw, key))
        except TypeError as exception:
            report_exception("parsing manifest", exception)
            exit(1)
