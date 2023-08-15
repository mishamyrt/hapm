"""asd"""
from __future__ import annotations

from argparse import BooleanOptionalAction
from os import environ

from arrrgs import arg

from hapm.manager import PackageManager
from hapm.manifest import Manifest
from hapm.report import (
    Progress,
    report_diff,
    report_latest,
    report_no_token,
    report_summary,
)

TOKEN_VAR = 'GITHUB_PAT'

unstable_arg = arg('--allow-unstable', '-u',
        action=BooleanOptionalAction,
        help="Removes the restriction to stable versions when searching for updates")

def load_manifest(args) -> Manifest:
    """Loads manifest file"""
    manifest = Manifest(args.manifest)
    manifest.load()
    return manifest


def synchronize(args, store: PackageManager, stable_only=True, manifest=None):
    """Synchronizes local versions of components with the manifest."""
    if manifest is None:
        manifest = load_manifest(args)
    progress = Progress()
    if len(manifest.has_latest) > 0:
        report_latest(manifest.has_latest)
        progress.start("Search for the latest versions")
    diff = store.diff(manifest.values, stable_only)
    if len(manifest.has_latest) > 0:
        progress.stop()
    report_diff(diff)
    if args.dry is True:
        return
    if len(diff) > 0:
        assert_warn_token()
        progress = Progress()
        progress.start("Synchronizing the changes")
        store.apply(diff)
        progress.stop()
    report_summary(diff)
    manifest.dump()

def get_token() -> str | None:
    """Returns token from environment variable or None if not set"""
    if TOKEN_VAR not in environ:
        return None
    return environ[TOKEN_VAR]

def assert_warn_token():
    """Checks if environment variable with token is set and prints a warning if not"""
    if get_token() is None:
        report_no_token(TOKEN_VAR)
