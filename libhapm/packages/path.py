"""HAPM repository path helpers"""


def repo_name(url: str) -> str:
    """Extracts the repository name from the url."""
    parts = url.split('/')
    return parts[len(parts) - 1]
