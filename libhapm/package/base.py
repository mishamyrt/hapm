"""HAPM package"""
from __future__ import annotations

from os.path import join
from shutil import rmtree

from .description import PackageDescription


def repo_name(url: str) -> str:
    """Extracts the repository name from the url"""
    parts = url.split('/')
    return parts[len(parts) - 1]

class BasePackage:
    """This is an abstract package controller class.
    The class that implements it must be able to control a certain type of package"""

    # Must be overlapped by a child
    kind: str

    # Package properties
    url: str
    version: str
    path: str
    name: str

    def __init__(self, description: PackageDescription, root_path: str):
        self.url = description["url"]
        self.version = description["version"]

        self.name = repo_name(self.url)
        self.path = join(root_path, self.name)

    # Built-in methods that will be useful in children

    def description(self) -> PackageDescription:
        """Returns the package description as a typed dict"""
        return {
            "url": self.url,
            "kind": self.kind,
            "version": self.version
        }

    def destroy(self) -> None:
        """Deletes the package from the file system"""
        rmtree(self.path)

    # Abstract methods to be implemented by all types of package handlers

    def initialize(self) -> None:
        """Method will be called if the entity is created for the first time.
        It should initialise the files on the system"""

    def load(self) -> None:
        """The method will be called if the entity already exists in the file system"""

    def switch(self, version: str) -> None:
        """Method should switch the version of the installed package"""

    def export(self, path: str) -> None:
        """Method should offload the package payload to the specified folder"""

    def latest_version(self) -> str:
        """The method should look for package updates and return latest stable version"""
