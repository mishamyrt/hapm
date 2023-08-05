"""Export functions for HAPM application"""
from arrrgs import arg, command

from hapm.manager import PackageManager


@command(
    arg('path', default=None, help="Output path")
)
def export(args, store: PackageManager):
    """Synchronizes local versions of components with the manifest."""
    store.export(args.path)
