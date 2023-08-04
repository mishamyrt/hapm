"""Utility functions for working with versions"""
from re import match
from typing import List

from pkg_resources import parse_version

STABLE_VERSION_RE = r'^v?\d+\.\d+(\.\d+)?$'

def find_latest(tags: List[str], stable_only: bool) -> str:
    """Finds the latest version in the list.
    Excludes unstable releases if the flag is specified"""
    latest = '0.0.0'
    for tag in tags:
        if stable_only and not is_stable(tag):
            continue
        if is_newer(latest, tag):
            latest = tag
    return latest

def is_stable(version: str) -> bool:
    """Checks if the version is stable"""
    return match(STABLE_VERSION_RE, version) is not None

def is_newer(current: str, new: str) -> bool:
    """Compares versions"""
    try:
        return parse_version(new) > parse_version(current)
    # pylint: disable-next=broad-except
    except Exception():
        return False
