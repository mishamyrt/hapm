"""HAPM repository sync utils"""

from __future__ import annotations
from shutil import rmtree
from os.path import isdir
from typing import Callable, Optional

from git import Repo

from .tag import get_tag


def is_correct_branch(repo: Repo, target: str) -> bool:
    """Checks that the local repository matches the specified branch or tag on the remote origin"""
    if repo.is_dirty():
        return False
    current = None
    try:
        current = repo.active_branch
    except TypeError:
        current = get_tag(repo)

    if current is None or current != target:
        return False

    return True

def refresh(repo: Repo, target: str):
    repo.git.reset('--hard')
    repo.git.pull()

def sync(url: str,
             tag: str,
             path: str,
             progress: Optional[Callable] = None):
    """Synchronises the specified version of the repository tag with the local repository.
        If the version is different or the repository has not yet been downloaded, clones it."""
    if not isdir(path):
        return Repo.clone_from(url, path, progress=progress, branch=tag)
    
    repo = Repo(path)
    if is_correct_branch(repo, tag):
        return repo

    origin = repo.remote()
    [remote_branch] = origin.fetch(tag)
    if remote_branch.commit != repo.commit():
        repo.git.checkout(remote_branch.commit, force=True)
        return False
    rmtree(path)
