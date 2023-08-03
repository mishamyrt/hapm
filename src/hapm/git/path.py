"""HAPM Github path utils"""

def repo_name(full_name: str) -> str:
    """Extracts the repository name from the full_name."""
    parts = full_name.split('/')
    return parts[len(parts) - 1]

def repo_url(full_name: str) -> str:
    """Builds a link to a repository by name."""
    return f"https://github.com/{full_name}"
