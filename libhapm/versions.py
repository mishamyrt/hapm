"""Вспомогательные функции для работы с версиями"""
from re import match
from pkg_resources import parse_version

STABLE_VERSION_RE = r'^v?\d+\.\d+(\.\d+)?$'

def find_latest_stable(tags: dict) -> str:
    """Находит последнюю стабильную версию в списке"""
    latest = '0.0.0'
    for tag in tags:
        if not is_stable(tag):
            continue
        if is_newer(latest, tag):
            latest = tag
    return latest

def is_stable(version: str) -> bool:
    """Проверяет является ли версия стабильной"""
    return match(STABLE_VERSION_RE, version) is not None

def is_newer(current: str, new: str) -> bool:
    """Сравнивает версии"""
    try:
        return parse_version(new) > parse_version(current)
    # pylint: disable-next=broad-except
    except Exception():
        return False
