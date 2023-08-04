"""Update search function for the HAPM application"""
from argparse import BooleanOptionalAction

from arrrgs import arg, command

from hapm.manager import PackageManager
from hapm.report import Progress, report_diff, report_warning


@command(
    arg('--allow-unstable', '-u',
        action=BooleanOptionalAction,
        help="Removes the restriction to stable versions when searching for updates"),
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
    report_diff(diff)
