"""Version format parsers"""
import re
from typing import Optional

ValueSegments = list[int]
SuffixSegments = Optional[list[str]]
VersionParts = tuple[ValueSegments, SuffixSegments]

_VERSION_RE = r"^v?(\d+(?:(?:\.\d+)+)?)(\.?[0-9A-Za-z-\.]+)?$"

class InvalidVersion(ValueError):
    """Raises when version format is invalid"""

def _parse_segments(segment_values: str) -> ValueSegments:
    segments = []
    for segment in segment_values.split("."):
        if len(segment) == 0:
            continue
        if not segment.isnumeric():
            raise InvalidVersion(f"Invalid segments value: {segment_values}")
        if int(segment) < 0:
            raise InvalidVersion(f"Invalid segments value: {segment_values}")
        segments.append(int(segment))
    return segments

def _parse_suffix(suffix: str) -> SuffixSegments:
    if suffix == "" or suffix is None:
        return None

    if suffix[0] == "-" or suffix[0] == ".":
        if len(suffix) == 1:
            raise InvalidVersion(f"Invalid suffix: {suffix}")
        suffix = suffix[1:]
    if "." not in suffix:
        return [suffix]
    return suffix.split(".")

def parse_version(version_expr: str) -> VersionParts:
    """
    Parses a version and returns its segments and type.
    Parser designed to be fail-safe and attempts
    to parse the version as much as possible
    """
    if not version_expr:
        raise InvalidVersion("Version is empty")
    if ".." in version_expr:
        raise InvalidVersion(f"Repeated dots in version: {version_expr}")
    match = re.match(_VERSION_RE, version_expr)
    if match is None:
        raise InvalidVersion(f"Invalid version format: {version_expr}")
    value_segments = _parse_segments(match.group(1))
    suffix_segments = _parse_suffix(match.group(2))
    return value_segments, suffix_segments
