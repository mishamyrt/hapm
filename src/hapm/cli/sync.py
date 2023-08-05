"""Syncing function for the HAPM application"""

from arrrgs import command

from hapm.manager import PackageManager

from .common import synchronize


@command()
def sync(args, store: PackageManager):
    """Synchronizes local versions of components with the manifest."""
    synchronize(args, store)
