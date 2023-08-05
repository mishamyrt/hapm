"""HAPM manifest controller"""
import sys
from typing import List

from ruamel.yaml import YAML, safe_load

from hapm.package import PackageDescription
from hapm.report import report_exception

from .category import parse_category

dumper = YAML()
dumper.indent(mapping=2, sequence=4, offset=2)

class Manifest:
    """HAPM manifest controller"""
    _encoding: str

    values: List[PackageDescription] = []
    has_latest: List[str] = []

    def __init__(self, path: str, encoding="utf-8"):
        self.path = path
        self._encoding = encoding

    def set(self, full_name: str, version: str, kind=None):
        """Add or update package"""
        for (i, _) in enumerate(self.values):
            if self.values[i]["full_name"] == full_name:
                self.values[i]["version"] = version
                return
        if kind is None:
            raise TypeError("Package type is not declared")
        self.values.append({
            "full_name": full_name,
            "version": version,
            "kind": kind
        })

    def dump(self) -> None:
        """Save manifest to file"""
        content = {}
        for package in self.values:
            full_name = package["full_name"]
            kind = package["kind"]
            version = package["version"]
            location = f"{full_name}@{version}"
            if kind in content:
                content[package["kind"]].append(location)
            else:
                content[package["kind"]] = [location]
        with open(self.path, "w", encoding="utf-8") as file:
            dumper.dump(content, file)

    def load(self) -> List[PackageDescription]:
        """Reads the manifest file and parses its contents"""
        with open(self.path, "r", encoding="utf-8") as stream:
            raw = safe_load(stream)
        try:
            for key in raw:
                values = parse_category(raw, key)
                for value in values:
                    if value["version"] == "latest":
                        self.has_latest.append(value["full_name"])
                self.values.extend(values)
        except TypeError as exception:
            report_exception("parsing manifest", exception)
            sys.exit(1)
