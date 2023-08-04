"""Repository utils"""

def repo_name(full_name: str) -> str:
    """Extracts the repository name from the full_name."""
    parts = full_name.split('/')
    return parts[len(parts) - 1]
