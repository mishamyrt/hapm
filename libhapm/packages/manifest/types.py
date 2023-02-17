from typing import List, Optional, TypedDict

class ManifestDict(TypedDict):
    """List of packages divided into categories"""
    integrations: Optional[List[str]]
    lovelace: Optional[List[str]]
    themes: Optional[List[str]]