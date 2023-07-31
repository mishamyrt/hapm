from __future__ import annotations
from github import Github, Repository, GithubException
from base64 import b64decode
from os import environ
import requests

plugin_url = "https://github.com/PiotrMachowski/lovelace-xiaomi-vacuum-map-card"
plugin_version = "v2.1.2"

# plugin_url = "https://github.com/custom-cards/spotify-card"
# plugin_version = "v2.4.0"

GITHUB_PREFIX = "https://github.com/"


else:
    token = environ['GITHUB_PAT']
api = Github(token)

def repo_name(url: str) -> str:
    """Extracts the repository name from the url"""
    parts = url.split('/')
    return parts[len(parts) - 1]

def check_file(repo: Repository, path: str, branch: str) -> str | None:
    try:
        content = repo.get_contents(path, ref=branch)
        return b64decode(content.content)
    except GithubException:
        return None

def check_dist(repo: Repository, filename: str, branch: str) -> str | None:
    return check_file(repo, f"dist/{filename}", branch=branch)

def check_release(repo: Repository, filename: str, branch: str) -> str | None:
    release = repo.get_release(branch)
    for asset in release.assets:
        if asset.name == filename:
            response = requests.get(asset.browser_download_url, timeout=10)
            return response.content
    return None

def check_root(repo: Repository, filename: str, branch: str) -> str | None:
    return check_file(repo, filename, branch=branch)


if not plugin_url.startswith(GITHUB_PREFIX):
    raise TypeError

repository = api.get_repo(plugin_url.replace(GITHUB_PREFIX, ""))
plugin_file = f"{repo_name(plugin_url)}.js"
if plugin_file.startswith("lovelace-"):
    plugin_file = plugin_file.replace("lovelace-", "")
print("dist", check_dist(repository, plugin_file, plugin_version) is not None)
print("release", check_release(repository, plugin_file, plugin_version) is not None)
print("root", check_root(repository, plugin_file, plugin_version) is not None)
# release = repo.get_release('v2.1.2')
# for asset in release.assets:
#     print(asset.name)
