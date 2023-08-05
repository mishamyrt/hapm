
from arrrgs import arg, command

from hapm.manager import PACKAGE_HANDLERS, PackageManager
from hapm.manifest import parse_location
from hapm.report import report_exception, report_warning

from .common import load_manifest, synchronize


@command(
    arg('url', default=None, help="Output path", nargs='+'),
    arg('--type', '-t',
        default=None, type=str, choices=list(PACKAGE_HANDLERS.keys()),
        help="Packages type. Required parameter if a new package is installed"),
)
def install(args, store: PackageManager):
    """Synchronizes local versions of components with the manifest."""
    manifest = load_manifest(args)
    for url in args.url:
        location = parse_location(url)
        try:
            manifest.set(location["full_name"], location["version"], args.type)
        except TypeError as exception:
            report_exception("installing package", exception)
            if args.type is None:
                report_warning("--type parameter is not specified.\n"+
                               "This option is required when installing new packages")
    synchronize(args, store, manifest)
