"""asd"""
from hapm.manager import PackageManager
from hapm.manifest import Manifest
from hapm.report import Progress, report_diff, report_latest, report_summary


def load_manifest(args) -> Manifest:
    """Loads manifest file"""
    manifest = Manifest(args.manifest)
    manifest.load()
    if len(manifest.has_latest) > 0:
        report_latest(manifest.has_latest)
    return manifest


def synchronize(args, store: PackageManager, manifest=None):
    """Synchronizes local versions of components with the manifest."""
    if manifest is None:
        manifest = load_manifest(args)
    diff = store.diff(manifest.values)
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
