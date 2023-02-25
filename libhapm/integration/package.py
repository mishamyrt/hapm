"""HAPM integration package module"""
from __future__ import annotations

import tarfile
from os import listdir, remove
from os.path import join
from shutil import copytree, rmtree

from requests import get

from libhapm.package import BasePackage


class IntegrationPackage(BasePackage):
    """IntegrationPackage represent custom_components packages"""

    kind = "integrations"
    extension = "tar.gz"

    def _download_archive(self, version: str):
        response = get(f"https://github.com/{self.full_name}/tarball/{self.version}", allow_redirects=True, timeout=60)
        with open(self.path(version), "wb") as file:
            file.write(response.content)

    def initialize(self) -> None:
        """Method will be called if the entity is created for the first time.
        It should initialise the files on the system"""
        self._download_archive(self.version)

    def switch(self, version: str) -> None:
        """Method should switch the version of the installed package"""
        self._download_archive(version)
        remove(self.path())
        self.version = version

    def export(self, path: str) -> None:
        """Method should offload the package payload to the specified folder"""
        with tarfile.open(self.path()) as file:
            all_members = file.getmembers()
            root = file.getmembers()[0].path
            components_path = join(root, "custom_components")
            target_members = []
            for member in all_members:
                if member.path.startswith(components_path):
                    target_members.append(member)
            file.extractall(members=target_members, path=path)
        exported_components_path = join(path, components_path)
        component = listdir(exported_components_path)[0]
        copytree(join(exported_components_path, component), join(path, component))
        rmtree(join(path, root))

    def latest_version(self) -> str:
        """The method should look for package updates and return latest stable version"""
        return 'v1.0.0'
        # self.repo.remote().fetch()
        # tags = []
        # for tag in self.repo.tags:
        #     tags.append(str(tag))
        # return find_latest_stable(tags)
