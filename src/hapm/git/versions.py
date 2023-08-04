"""HAPM Github versions utils"""
from typing import List, Optional

from github import Github, Tag


def get_versions(full_name: str, api_token: Optional[str]) -> List[Tag.Tag]:
    """Gets the available versions of the GitHub repository"""
    api = Github(api_token)
    repo = api.get_repo(full_name)
    tags_list = repo.get_tags()
    tags = []
    for tag in tags_list:
        tags.append(tag.name)
    return tags
