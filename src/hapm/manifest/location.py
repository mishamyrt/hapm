"""HAPM manifest parsing utils"""
from __future__ import annotations

from re import match
from typing import TypedDict
from urllib.parse import ParseResult, urlparse


class PackageLocation(TypedDict):
    """Dict describing the Home Assistant package repository location"""
    full_name: str
    version: str

def safe_urlparse(url: str) -> ParseResult | None:
    """Safe version of urlparse. Returns None on exceptions"""
    try:
        return urlparse(url)
    # pylint: disable-next=bare-except
    except: # noqa: E722
        return None

def parse_github_url(url: str) -> ParseResult | None:
    """Parses URLs to GitHub.
    If the URL is incorrect or does not point to github.com - returns None"""
    result = safe_urlparse(url)
    if result is None or not all([result.scheme, result.netloc]):
        return None
    if not result.netloc == "github.com":
        return None
    return result

def parse_location_url(url: str) -> PackageLocation | None:
    result = parse_github_url(url)
    if result is None:
        return None
    url_match = match(r'/(.*)/(.*)/releases/tag/(.*)', result.path)
    if url_match is not None:
        return {
            "full_name": f"{url_match.group(1)}/{url_match.group(2)}",
            "version": url_match.group(3)
        }
    if result.path.count('/') != 2:
        return None
    return {
        "full_name": result.path[1:],
        "version": "latest"
    }

def parse_package_name(location: str) -> PackageLocation | None:
    """Parses the manifest entry to the address and version"""
    result = match(r"(.[^@]*)(@(.*))?", location)
    if result is None:
        return None
    full_name = result.group(1)
    if full_name is None:
        return None
    version = result.group(3)
    if version is None:
        version = "latest"
    return {
        "full_name": full_name,
        "version": version
    }

def parse_location(package: str) -> PackageLocation | None:
    """Parses the string of the packet.
    The `package` input can be a repository reference, tag or full name,
    in the format username/repository@version.
    If no version is specified or it could not be determined from the reference,
    the latest version will be used"""
    if package.startswith("github.com"):
        package = "https://{url}"
    parsers = [
        parse_location_url,
        parse_package_name
    ]
    location = None
    for parse in parsers:
        location = parse(package)
        if location is not None:
            return location
    return location

