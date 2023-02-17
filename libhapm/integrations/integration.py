"""HAPM storage integration utils"""
from __future__ import annotations
from typing import TypedDict
from dataclasses import dataclass
from os import listdir
from os.path import join, isdir
from distutils.dir_util import copy_tree
from shutil import rmtree

from git import Repo

from libhapm.packages import repo_name, Package

# Integration sync statuses
IN_SYNC = 0
SHOULD_ADD = 1
SHOULD_SYNC = 2


class IntegrationLock(TypedDict):
    """Minimum object for lockfile from which an integration instance can be recovered"""
    name: str
    url: str
    version: str
    path: str


@dataclass
class Integration:
    name: str
    path: str
    url: str
    version: str
    repo: Repo

    def to_lock(self) -> IntegrationLock:
        lock: IntegrationLock = {
            "name": self.name,
            "url": self.url,
            "version": self.version,
            "path": self.path
        }
        return lock

    def to_package(self) -> Package:
        package: Package = {"url": self.url, "version": self.version}
        return package

    def switch_to(self, version: str) -> bool:
        self.repo.git.reset('--hard')
        data = self.repo.remote().fetch()
        for i in data:
            print(i.ref)
        self.repo.git.checkout(version)
        self.version = version

    def remove(self):
        """Deletes the package from the file system"""
        rmtree(self.path)

    def export(self, path: str):
        """Deletes the package from the file system"""
        components_path = join(self.path, "custom_components")
        components = listdir(components_path)
        for component in components:
            source = join(components_path, component)
            destination = join(path, component)
            if isdir(destination):
                rmtree(destination)
            if isdir(source):
                copy_tree(source, destination)

    @staticmethod
    def from_lock(lock: IntegrationLock) -> Integration:
        """Restores a copy from the lockfile"""
        return Integration(
            name=lock["name"],
            path=lock["path"],
            url=lock["url"],
            version=lock["version"],
            repo=Repo(lock["path"])
        )

    @staticmethod
    def from_package(package: Package, path: str) -> Integration:
        """Initiates a new integration in storage"""
        name = repo_name(package["url"])
        repo_path = f"{path}/{name}"
        print(f"repo_path: {repo_path}")
        return Integration(name=name,
                           path=repo_path,
                           url=package["url"],
                           version=package["version"],
                           repo=Repo.clone_from(package["url"],
                                                repo_path,
                                                branch=package["version"]))
