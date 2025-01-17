"""Version utilities"""
from .parse import parse_version


class Version:
    """
    Version representation.
    Allows to easily compare versions.
    """
    _original: str
    _value = []
    _suffix = []

    def __init__(self, version: str):
        self._original = version
        self._value, self._suffix = parse_version(version)

    @property
    def is_stable(self) -> bool:
        """Checks if the version is stable"""
        return self._suffix is None

    def __str__(self):
        return self._original

    def __lt__(self, other: "Version") -> bool:
        if len(self._value) != len(other._value):
            return len(self._value) < len(other._value)
        if self._value != other._value:
            return self._value < other._value
        if self._suffix is None:
            return False
        if other._suffix is None:
            return True
        if len(self._suffix) != len(other._suffix):
            return len(self._suffix) < len(other._suffix)
        return self._suffix < other._suffix

    def __eq__(self, other: "Version") -> bool:
        return self._value == other._value and self._suffix == other._suffix

    def __gt__(self, other: "Version") -> bool:
        return not self < other and not self == other

    def __le__(self, other: "Version") -> bool:
        return self < other or self == other

    def __ge__(self, other: "Version") -> bool:
        return self > other or self == other

    def __ne__(self, other: "Version") -> bool:
        return not self == other
