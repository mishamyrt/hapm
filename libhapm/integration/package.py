"""HAPM integration package module"""
from __future__ import annotations

from distutils.dir_util import copy_tree
from os import listdir
from os.path import isdir, join
from shutil import rmtree

from git import Repo

from libhapm.package import BasePackage
from libhapm.versions import find_latest_stable


class IntegrationPackage(BasePackage):
    """IntegrationPackage represent custom_components packages"""
    repo: Repo

    kind = "integrations"

    def initialize(self) -> None:
        """Method will be called if the entity is created for the first time.
        It should initialise the files on the system"""
        self.repo = Repo.clone_from(self.url, self.path, branch=self.version)

    def load(self) -> None:
        """Method will be called if the entity is created for the first time.
        It should initialise the files on the system"""
        self.repo = Repo(self.path)

    def switch(self, version: str) -> None:
        """Method should switch the version of the installed package"""
        self.repo.git.reset('--hard')
        self.repo.remote().fetch()
        self.repo.git.checkout(version)
        self.version = version

    def export(self, path: str) -> None:
        """Method should offload the package payload to the specified folder"""
        components_path = join(self.path, "custom_components")
        components = listdir(components_path)
        if len(components) < 1:
            raise FileNotFoundError(f"There are no components in the '{components_path}' folder")
        source = join(components_path, components[0])
        destination = join(path, components[0])
        if isdir(destination):
            rmtree(destination)
        if isdir(source):
            copy_tree(source, destination)

    def latest_version(self) -> str:
        """The method should look for package updates and return latest stable version"""
        self.repo.remote().fetch()
        tags = []
        for tag in self.repo.tags:
            tags.append(str(tag))
        return find_latest_stable(tags)
