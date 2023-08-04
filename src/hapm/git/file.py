"""HAPM Github file utils"""
from __future__ import annotations

from base64 import b64decode

from github import GithubException
from github.Repository import Repository
from requests import get


def get_tree_file(repo: Repository, branch: str, path: str) -> str | None:
    """Gets the file from the repository by path and returns its string content."""
    try:
        content = repo.get_contents(path, ref=branch)
        return b64decode(content.content)
    except GithubException:
        return None

def get_release_file(repo: Repository, branch: str, filename: str) -> str | None:
    """Gets the file from the release by name and returns its string content."""
    release = repo.get_release(branch)
    for asset in release.assets:
        if asset.name == filename:
            response = get(asset.browser_download_url, timeout=10)
            return response.content
    return None
