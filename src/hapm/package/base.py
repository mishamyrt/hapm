"""HAPM package"""
from __future__ import annotations

from os import remove
from os.path import join

from github import Github

from hapm.git import repo_name
from hapm.versions import find_latest_stable

from .description import PackageDescription


class BasePackage:
    """This is an abstract package controller class.
    The class that implements it must be able to control a certain type of package"""

    # Must be overridden by a child
    kind: str

    # Package properties
    full_name: str
    version: str
    basepath: str
    extension: str
    name: str
    _api_token: str
    _cache_path: str

    def __init__(self, description: PackageDescription, root_path: str, token: str):
        self.full_name = description["full_name"]
        self.version = description["version"]
        self._api_token = token

        self.name = repo_name(self.full_name)
        self.basepath = join(root_path, self.full_name.replace('/', '-'))

    # Built-in methods that will be useful in children

    def description(self) -> PackageDescription:
        """Returns the package description as a typed dict"""
        return {
            "full_name": self.full_name,
            "kind": self.kind,
            "version": self.version
        }

    def path(self, version=None) -> str:
        """Returns the path to the integration cache file"""
        if version is None:
            version = self.version
        return f"{self.basepath}@{version}.{self.extension}"

    def destroy(self) -> None:
        """Deletes the package from the file system"""
        remove(self.path())

    def latest_version(self) -> str:
        api = Github(self._api_token)
        repo = api.get_repo(self.full_name)
        tags_list = repo.get_tags()
        tags = []
        for tag in tags_list:
            tags.append(tag.name)
        return find_latest_stable(tags)
        

    # Abstract methods to be implemented by all types of package handlers

    def initialize(self) -> None:
        """Method will be called if the entity is created for the first time.
        It should initialise the files on the system"""

    def switch(self, version: str) -> None:
        """Method should switch the version of the installed package"""

    def export(self, path: str) -> None:
        """Method should offload the package payload to the specified folder"""

    # Lifecycle hooks

    @staticmethod
    def pre_export(path: str):
        """This method is called when you starting exporting packages of a certain kind"""

    @staticmethod
    def post_export(path: str):
        """This method is called after you finish exporting packages of a certain kind"""
