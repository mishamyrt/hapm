"""HAPM lockfile module"""
from __future__ import annotations

from json import dump, load
from os.path import isfile
from typing import List

from hapm.package import PackageDescription


class Lockfile:
    """Represents an entity that manages a file
    with descriptions of current package versions"""
    _path: str
    _encoding: str

    def __init__(self, path: str, encoding="utf-8"):
        self._path = path
        self._encoding = encoding

    def exists(self) -> bool:
        """Checks if lockfile is exists"""
        return isfile(self._path)

    def dump(self, descriptions: List[PackageDescription]) -> None:
        """Writes a list of packages to a file"""
        with open(self._path, "w", encoding=self._encoding) as file:
            dump(descriptions, file)

    def load(self) -> List[PackageDescription]:
        """Loads a list of packages from a file"""
        with open(self._path, "r", encoding=self._encoding) as file:
            return load(file)
