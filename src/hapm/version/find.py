"""Find version utils"""

from .parse import InvalidVersion
from .version import Version


def find_latest(tags: list[str], stable_only: bool) -> Version:
    """Finds the latest version in the list.
    Excludes unstable releases if the flag is specified"""
    latest = Version("0.0.0")
    versions = []
    for tag in tags:
        try:
            versions.append(Version(tag))
        except InvalidVersion:
            continue
        versions.append(Version(tag))
    for version in versions:
        if not version.is_stable and stable_only:
            continue
        if version > latest:
            latest = version
    return str(latest)
