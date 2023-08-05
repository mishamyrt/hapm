"""HAPM manager module"""
from os import mkdir
from os.path import isdir, join
from shutil import rmtree
from typing import Dict, List

from hapm.git import get_versions
from hapm.integration import IntegrationPackage
from hapm.manifest import Manifest, PackageLocation
from hapm.package import BasePackage, PackageDescription
from hapm.plugin import PluginPackage
from hapm.versions import find_latest, is_newer

from .diff import PackageDiff
from .lockfile import Lockfile

PACKAGE_HANDLERS = {
    IntegrationPackage.kind: IntegrationPackage,
    PluginPackage.kind: PluginPackage
}


class PackageManager:
    """The controller that manages the packets in the storage"""

    _packages: Dict[str, BasePackage] = {}
    _api_token: str

    def __init__(self, path: str, token: str, lockfile_name="_lock.json"):
        self._path = path

        lock_path = join(self._path, lockfile_name)
        self._lock = Lockfile(lock_path)
        self._api_token = token

        if isdir(self._path):
            if self._lock.exists():
                self._boot_from_lock()
        else:
            mkdir(self._path)

    def get_versions(self, location: PackageLocation) -> List[str]:
        """Returns package version by location"""
        return get_versions(location["full_name"], self._api_token)

    def _boot_from_lock(self):
        descriptions = self._lock.load()
        if len(descriptions) == 0:
            return
        for description in descriptions:
            package = PACKAGE_HANDLERS[description["kind"]](
                description, self._path, self._api_token)
            self._packages[package.full_name] = package

    def diff(self, update: List[PackageDescription], stable_only=True) -> List[PackageDiff]:
        """Finds the difference between the current state and the list of packets received.
        Returns the modified package description"""
        update_full_names: List[str] = []
        diffs: List[PackageDiff] = []
        for description in update:
            if description["version"] == "latest":
                versions = get_versions(description["full_name"], self._api_token)
                description["version"] = find_latest(versions, stable_only)
            full_name, version = description["full_name"], description["version"]
            update_full_names.append(full_name)
            diff: PackageDiff = description.copy()
            if full_name in self._packages:
                current_version = self._packages[full_name].version
                if current_version != version:
                    diff["current_version"] = current_version
                    diff["operation"] = "switch"
            else:
                diff["operation"] = "add"
            if "operation" in diff:
                diffs.append(diff)
        for (full_name, integration) in self._packages.items():
            try:
                if update_full_names.index(full_name):
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
        full_names_to_remove = []
        for diff in diffs:
            operation = diff["operation"]
            full_name = diff["full_name"]
            version = diff["version"]
            if operation == "add":
                package = PACKAGE_HANDLERS[diff["kind"]](
                    diff, self._path, self._api_token)
                package.setup()
                self._packages[full_name] = package
            elif operation == "delete":
                self._packages[full_name].destroy()
                full_names_to_remove.append(full_name)
            else:
                package = self._packages[full_name]
                package.switch(version)
        # Delete keys in a separate loop so as not to change the iterated list
        for full_name in full_names_to_remove:
            self._packages.pop(full_name, None)
        self._lock.dump(self.descriptions())

    def export(self, path: str):
        """Deletes the package from the file system"""
        if isdir(path):
            rmtree(path)
        mkdir(path)
        kinds = []
        for (_, integration) in self._packages.items():
            if integration.kind not in kinds:
                kinds.append(integration.kind)
                PACKAGE_HANDLERS[integration.kind].pre_export(path)
            integration.export(path)
        for kind in kinds:
            PACKAGE_HANDLERS[kind].post_export(path)


    def updates(self, stable_only=True) -> List[PackageDiff]:
        """Searches for updates for packages, returns list of available updates."""
        updates: List[PackageDiff] = []
        for (_, package) in self._packages.items():
            latest_version = package.latest_version(stable_only)
            if is_newer(package.version, latest_version):
                updates.append({
                    "full_name": package.full_name,
                    "kind": package.kind,
                    "version": latest_version,
                    "current_version": package.version,
                    "operation": "switch"
                })
        return updates

    def descriptions(self) -> List[PackageDescription]:
        """Collects a list of current package descriptions"""
        return [package.description() for _, package in self._packages.items()]
