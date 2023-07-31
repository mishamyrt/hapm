"""HAPM integration package module"""
from __future__ import annotations

import tarfile
from os import listdir, mkdir, remove
from os.path import join
from shutil import copytree, rmtree

from libhapm.github import get_tarball
from libhapm.package import BasePackage

FOLDER_NAME = "custom_components"


class IntegrationPackage(BasePackage):
    """IntegrationPackage represent custom_components packages"""

    kind = "integrations"
    extension = "tar.gz"

    def _download_tarball(self, version: str):
        content = get_tarball(self.full_name, version)
        with open(self.path(version), "wb") as file:
            file.write(content)

    def initialize(self) -> None:
        """Method will be called if the entity is created for the first time.
        It should initialise the files on the system"""
        self._download_tarball(self.version)

    def switch(self, version: str) -> None:
        """Method should switch the version of the installed package"""
        self._download_tarball(version)
        remove(self.path())
        self.version = version

    def export(self, path: str) -> None:
        """Method should offload the package payload to the specified folder"""
        with tarfile.open(self.path()) as file:
            all_members = file.getmembers()
            root = file.getmembers()[0].path
            components_path = join(root, FOLDER_NAME)
            target_members = []
            for member in all_members:
                if member.path.startswith(components_path):
                    target_members.append(member)
            file.extractall(members=target_members, path=path)
        exported_components_path = join(path, components_path)
        component = listdir(exported_components_path)[0]
        copytree(join(exported_components_path, component), join(path, FOLDER_NAME, component))
        rmtree(join(path, root))

    @staticmethod
    def pre_export(path: str):
        """This method is called when you starting exporting packages of a certain kind"""
        mkdir(join(path, FOLDER_NAME))

    @staticmethod
    def post_export(path: str):
        """Do nothing"""
