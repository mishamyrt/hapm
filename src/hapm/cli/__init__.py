"""HAPM CLI application"""
from argparse import BooleanOptionalAction
from os import environ

from arrrgs import arg, command, global_args, run

from hapm.manager import PackageManager
from hapm.manifest import Manifest

from .common import get_token

# Commands
from .export import export
from .install import install, sync
from .versions import list_packages, updates, versions

STORAGE_DIR = ".hapm"
MANIFEST_PATH = "hapm.yaml"

global_args(
    arg('--manifest', '-m', default=MANIFEST_PATH, help="Manifest path"),
    arg('--storage', '-s', default=STORAGE_DIR, help="Storage location"),
    arg('--dry', '-d',
        action=BooleanOptionalAction,
        help="Only print information. Do not make any changes to the files")
)

@command()
def init(args, store: PackageManager):
    """Manifest creation"""
    manifest = Manifest(args.manifest)
    manifest.init(store.supported_types())

def prepare(args):
    """Creates HAPM context"""
    token = get_token()
    return args, PackageManager(args.storage, token)

def start():
    """Application entrypoint"""
    run(prepare=prepare)
