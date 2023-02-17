"""HAPM packages module base"""

from typing import TypeVar, Generic, List

from .package import Package

T = TypeVar('T')

class PackagesModule(Generic[T]):
    """This is an abstract package controller class.
    The class that implements it must be able to control a certain type of package."""

    def __init__(self, lock:List[dict]=None):
        """Creates new PackagesModule"""

    def lock(self) -> List[dict]:
        """Returns list of current packages in lock format."""

    def apply(self, update: List[Package]) -> bool:
        """Applies the new configuration.
        Important: this method can make changes to the file system.
        Returns False if no changes were made."""

    def values(self) -> List[Package]:
        """Returns current items list."""

    def export(self, path: str) -> bool:
        """Dumps the required data from the packages to the specified folder."""

    @staticmethod
    def package_type() -> str:
        """Returns module package type"""
