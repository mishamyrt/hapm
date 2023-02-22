"""HAPM manager module"""
from os.path import isdir, join
from typing import Dict, List

from libhapm.integration import IntegrationPackage
from libhapm.package import BasePackage, PackageDescription
from libhapm.versions import is_newer

from .diff import PackageDiff
from .lockfile import Lockfile

PACKAGE_HANDLERS = {
    IntegrationPackage.kind: IntegrationPackage
}


class PackageManager:
    """The controller that manages the packets in the storage"""

    _packages: Dict[str, BasePackage] = {}

    def __init__(self, path: str, lockfile_name="lock.json"):
        self._path = path

        lock_path = join(self._path, lockfile_name)
        self._lock = Lockfile(lock_path)

        if isdir(self._path) and self._lock.exists():
            self._boot_from_lock()

    def _boot_from_lock(self):
        descriptions = self._lock.load()
        if len(descriptions) == 0:
            return
        for description in descriptions:
            package = PACKAGE_HANDLERS[description["kind"]](
                description, self._path)
            package.load()
            self._packages[package.url] = package

    def diff(self, update: List[PackageDescription]) -> List[PackageDiff]:
        """Finds the difference between the current state and the list of packets received.
        Returns the modified package description"""
        update_urls: List[str] = []
        diffs: List[PackageDiff] = []
        for description in update:
            url, version = description["url"], description["version"]
            update_urls.append(url)
            diff: PackageDiff = description.copy()
            if url in self._packages:
                current_version = self._packages[url].version
                if current_version == version:
                    continue
                else:
                    diff["current_version"] = current_version
                    diff["operation"] = "switch"
            else:
                diff["operation"] = "add"
            diffs.append(diff)
        for (url, integration) in self._packages.items():
            try:
                if update_urls.index(url):
                    continue
            except ValueError:
                diff: PackageDiff = integration.description()
                diff["operation"] = "delete"
                diffs.append(diff)
        return diffs

    def apply(self, diffs: List[PackageDiff]):
        """Applies the new configuration.
        Important: this method will make changes to the file system.
        Returns False if no changes were made."""
        urls_to_remove = []
        for diff in diffs:
            operation = diff["operation"]
            url = diff["url"]
            if operation == "add":
                package = PACKAGE_HANDLERS[diff["kind"]](
                    diff, self._path)
                package.initialize()
                self._packages[url] = package
            elif operation == "delete":
                self._packages[url].destroy()
                urls_to_remove.append(url)
            else:
                self._packages[url].switch(diff["version"])
        # Delete keys in a separate loop so as not to change the iterated list
        for url in urls_to_remove:
            self._packages.pop(url, None)
        self._lock.dump(self.descriptions())

    def export(self, kind: str, path: str):
        """Deletes the package from the file system"""
        for (_, integration) in self._packages.items():
            if integration.kind == kind:
                integration.export(path)

    def updates(self) -> List[PackageDescription]:
        """Searches for updates for packages, returns list of available updates."""
        updates: List[PackageDescription] = []
        for (_, package) in self._packages.items():
            latest_version = package.latest_version()
            if is_newer(package.version, latest_version):
                updates.append({
                    "url": package.url,
                    "kind": package.kind,
                    "version": latest_version
                })
        return updates

    def descriptions(self) -> List[PackageDescription]:
        """Collects a list of current package descriptions"""
        return [package.description() for _, package in self._packages.items()]
