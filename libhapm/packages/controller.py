"""HAPM packages controller module"""
from os.path import isdir, join
from os import mkdir
from typing import List, Dict

from libhapm.integrations import IntegrationsModule

from .manifest import Manifest
from .lock import LockFile
from .module import PackagesModule
from .package import Package

ACTIVE_MODULES = [
    IntegrationsModule
]


class PackagesController:
    """HAPM packages controller"""
    _lock: LockFile
    _path: str
    _modules: List[PackagesModule] = []

    def __init__(self, path: str):
        self._path = path

        lock_path = join(self._path, "lock.json")
        self._lock = LockFile(lock_path)

        if isdir(self._path):
            raw_lock = self._lock.load()
        else:
            raw_lock = {}

        for Module in ACTIVE_MODULES:
            package_type = Module.package_type()
            if package_type in raw_lock:
                lock = raw_lock[package_type]
            else:
                lock = None
            path = join(self._path, package_type)
            self._modules.append(Module(
                path=path,
                lock=lock
            ))

    def apply(self, manifest: Manifest) -> bool:
        """Applies the manifest to the current repository"""
        changed = False
        for module in self._modules:
            package_type = module.package_type()
            if package_type in manifest.values:
                module_changed = module.apply(manifest.values[package_type])
                changed = changed or module_changed
        if changed:
            self._dump_lock()
        return changed

    def values(self) -> Dict[str, Package]:
        values = {}
        for module in self._modules:
            values[module.package_type()] = module.values()
        return values

    def export(self, package_type: str, path: str):
        """Deletes the package from the file system"""
        if not isdir(path):
            mkdir(path)
        for module in self._modules:
            if module.package_type() == package_type:
                module.export(path)
                return


    def _dump_lock(self):
        self._lock.dump(self._modules)
