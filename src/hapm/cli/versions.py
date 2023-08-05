"""Package versions related function for the HAPM application"""
import sys

from arrrgs import arg, command

from hapm.manager import PackageManager
from hapm.manifest import parse_location
from hapm.report import (
    Progress,
    report_diff,
    report_packages,
    report_versions,
    report_warning,
    report_wrong_format,
)

from .common import unstable_arg


@command(
    unstable_arg,
)
def updates(args, store: PackageManager):
    """Prints available packages updates."""
    stable_only = not args.allow_unstable
    if not stable_only:
        report_warning("Search includes unstable versions")
    progress = Progress()
    progress.start("Looking for package updates")
    diff = store.updates(not args.allow_unstable)
    progress.stop()
    if len(diff) == 0:
        print("All packages is up to date")
        return
    report_diff(diff, True, True)

@command(
    arg('location', type=str, help="Package name or package URL")
)
def versions(args, store: PackageManager):
    """Prints the available versions for the package"""
    location = parse_location(args.location)
    if location is None:
        report_wrong_format(args.location)
        sys.exit(1)
    progress = Progress()
    progress.start("Looking for package versions")
    tags = store.get_versions(location)
    progress.stop()
    report_versions(location["full_name"], tags)


@command(name="list")
def list_packages(_, store: PackageManager):
    """Print current version of components."""
    packages = store.descriptions()
    report_packages(packages)
