from typing import Tuple, Optional, List
from re import match
from urllib.parse import urlparse

from libhapm.packages import Package

from .types import ManifestDict


def parse_entry(entry: str) -> Tuple[Optional[str], Optional[str]]:
    """Parses the manifest entry to the address and version"""
    parts = match(r"(.[^@]*)(@(.*))?", entry)
    if parts is None:
        return (None, None)
    parse_result = urlparse(parts.group(1))
    if parse_result.scheme != '':
        url = f"{parse_result.scheme}://{parse_result.netloc}{parse_result.path}"
    else:
        url = f"https://{parse_result.path}"
    print(url)
    return (url, parts.group(3))


def parse_category(manifest: ManifestDict, key: str) -> List[Package]:
    if key not in manifest:
        return []
    items: List[Package] = []
    for entry in manifest[key]:
        (url, version) = parse_entry(entry)
        if url is None:
            raise TypeError(f"Wrong entity: {entry}")
        if version is None:
            raise TypeError(f"Version is missing: {entry}")
        items.append({
            "url": url,
            "version": version
        })
    return items
