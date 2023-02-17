from __future__ import annotations
from typing import List
from json import dump, load

from .module import PackagesModule

LOCK_ENCODING = "utf-8"


class LockFile:
    _path: str

    def __init__(self, path: str):
        self._path = path

    def dump(self, modules: List[PackagesModule]):
        content = {}
        for module in modules:
            content[module.package_type()] = module.lock()
        with open(self._path, "w", encoding=LOCK_ENCODING) as f:
            dump(content, f)

    def load(self) -> dict:
        with open(self._path, "r", encoding=LOCK_ENCODING) as f:
            return load(f)
