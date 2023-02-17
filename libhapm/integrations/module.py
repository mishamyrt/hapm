"""HAPM integrations module"""
from typing import List, Dict

from ..packages import PackagesModule, Package
from .integration import Integration, IntegrationLock


class IntegrationsModule(PackagesModule):
    _items: Dict[str, Integration] = {}

    def __init__(self, path: str, lock: List[IntegrationLock] = None):
        self._path = path
        if lock is not None:
            for i_lock in lock:
                self._items[i_lock["url"]] = Integration.from_lock(i_lock)

    @staticmethod
    def package_type() -> str:
        """Returns module package type"""
        return "integrations"

    def lock(self) -> List[dict]:
        """Returns list of integrations in lock format."""
        return  [self._items[url].to_lock() for url in self._items]

    def apply(self, update: List[Package]) -> bool:
        """Applies the new configuration.
        Important: this method can make changes to the file system.
        Returns False if no changes were made."""
        changed = False
        updated_urls: List[str] = []
        for package in update:
            url = package["url"]
            version = package["version"]
            updated_urls.append(url)
            if url in self._items:
                integration = self._items[url]
                if integration.version == version:
                    continue
                else:
                    integration.switch_to(version)
                    changed = True
            else:
                self._items[url] = Integration.from_package(package, self._path)
                changed = True
        deleted = self._clean_to(updated_urls)
        return changed or deleted

    def values(self) -> List[Package]:
        """Returns current items list."""
        return  [integration.to_package() for _, integration in self._items.items()]

    def export(self, path: str):
        """Deletes the package from the file system"""
        for (_, integration) in self._items.items():
            integration.export(path)
        # content = listdir(f"{self.path}/custom_components")
        # copy_tree()

    def _clean_to(self, urls: List[str]) -> bool:
        """Deletes integrations that are not on the list"""
        changed = False
        urls_to_remove = []
        for (url, integration) in self._items.items():
            try:
                if urls.index(url) > -1:
                    continue
            except ValueError:
                integration.remove()
                urls_to_remove.append(url)
                changed = True
        for url in urls_to_remove:
            self._items.pop(url, None)
        return changed
