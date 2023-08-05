"""HAPM CLI application"""
from argparse import BooleanOptionalAction
from os import environ

from arrrgs import arg, command, global_args, run

from hapm.manager import PackageManager
from hapm.manifest import Manifest
from hapm.report import report_no_token

# Commands
from .export import export
from .install import install, sync
from .versions import list_packages, updates, versions

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
def init(args, store: PackageManager):
    """Manifest creation"""
    manifest = Manifest(args.manifest)
    manifest.init(store.supported_types())



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
