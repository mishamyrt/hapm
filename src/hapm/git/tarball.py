"""HAPM Github tarball utils"""
from __future__ import annotations

from requests import get

from .path import repo_url


def get_tarball(full_name: str, branch: str) -> bytes:
    """Downloads and returns bytes of the tar.gz file of the specified branch."""
    response = get(f"{repo_url(full_name)}/tarball/{branch}",
                   allow_redirects=True,
                   timeout=60)
    return response.content
