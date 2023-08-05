"""asd"""
from argparse import BooleanOptionalAction

from arrrgs import arg

from hapm.manager import PackageManager
from hapm.manifest import Manifest
from hapm.report import Progress, report_diff, report_latest, report_summary

unstable_arg = arg('--allow-unstable', '-u',
        action=BooleanOptionalAction,
        help="Removes the restriction to stable versions when searching for updates")

def load_manifest(args) -> Manifest:
    """Loads manifest file"""
    manifest = Manifest(args.manifest)
    manifest.load()
    if len(manifest.has_latest) > 0:
        report_latest(manifest.has_latest)
    return manifest


def synchronize(args, store: PackageManager, stable_only=True, manifest=None):
    """Synchronizes local versions of components with the manifest."""
    if manifest is None:
        manifest = load_manifest(args)
    diff = store.diff(manifest.values, stable_only)
    report_diff(diff)
    if args.dry is True:
        return
    if len(diff) > 0:
        progress = Progress()
        progress.start("Synchronizing the changes")
        store.apply(diff)
        progress.stop()
    report_summary(diff)
    manifest.dump()
