"""HAPM CLI application"""
from argparse import BooleanOptionalAction
from os import environ

from arrrgs import arg, command, global_args, run

from hapm.manager import PackageManager
from hapm.manifest import Manifest
from hapm.report import (
    Progress,
    report_diff,
    report_no_token,
    report_packages,
    report_summary,
)

from .versions import updates, versions

progress = Progress()

STORAGE_DIR = ".hapm"
MANIFEST_PATH = "hapm.yaml"
TOKEN_VAR = 'GITHUB_PAT'

global_args(
    arg('--manifest', '-m', default=MANIFEST_PATH, help="Manifest path"),
    arg('--storage', '-s', default=STORAGE_DIR, help="Storage location"),
    arg('--dry', '-d',
        action=BooleanOptionalAction,
        help="Only print information. Do not make any changes to the files")
)

@command()
def sync(args, store: PackageManager):
    """Synchronizes local versions of components with the manifest."""
    manifest = Manifest(args.manifest)
    manifest.load()
    diff = store.diff(manifest.values)
    report_diff(diff)
    if args.dry is True:
        return
    progress.start("Synchronizing the changes")
    store.apply(diff)
    progress.stop()
    report_summary(diff)


@command(
    arg('path', default=None, help="Output path")
)
def put(args, store: PackageManager):
    """Synchronizes local versions of components with the manifest."""
    store.export(args.path)



@command(name="list")
def list_packages(_, store: PackageManager):
    """Print current version of components."""
    packages = store.descriptions()
    report_packages(packages)


def prepare(args):
    """Creates HAPM context"""
    if TOKEN_VAR not in environ:
        report_no_token(TOKEN_VAR)
        token = None
    else:
        token = environ[TOKEN_VAR]
    return args, PackageManager(args.storage, token)

def start():
    """Application entrypoint"""
    run(prepare=prepare)
