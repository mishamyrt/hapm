"""HAPM manager module"""
from os.path import isdir, isfile, join
from typing import List, Dict

from libhapm.integration import IntegrationPackage
from libhapm.package import BasePackage, PackageDescription
from libhapm.versions import is_newer

from .lockfile import Lockfile

PACKAGE_HANDLERS = {
    IntegrationPackage.kind: IntegrationPackage
}


class PackageManager:

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

    def apply(self, update: List[PackageDescription]):
        """Applies the new configuration.
        Important: this method will make changes to the file system.
        Returns False if no changes were made."""
        changed = False
        updated_urls: List[str] = []
        for description in update:
            url, version = description["url"], description["version"]
            updated_urls.append(url)
            if url in self._packages:
                package = self._packages[url]
                if package.version == version:
                    continue
                else:
                    package.switch(version)
                    changed = True
            else:
                package = PACKAGE_HANDLERS[description["kind"]](
                    description, self._path)
                package.initialize()
                self._packages[url] = package
                changed = True
        deleted = self._clean_to(updated_urls)
        if changed or deleted:
            self._lock.dump(self.descriptions())
        return changed or deleted

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
        return [package.description() for _, package in self._packages.items()]

    def _clean_to(self, urls: List[str]) -> bool:
        """Deletes packages that are not on the list"""
        changed = False
        urls_to_remove = []
        for (url, integration) in self._packages.items():
            try:
                if urls.index(url) > -1:
                    continue
            except ValueError:
                integration.destroy()
                urls_to_remove.append(url)
                changed = True
        for url in urls_to_remove:
            self._packages.pop(url, None)
        return changed
