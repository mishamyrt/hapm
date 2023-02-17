"""Home Assistant Package module"""
from typing import TypedDict

class Package(TypedDict):
    """Basic type for a package"""
    url: str
    version: str
