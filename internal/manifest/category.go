package manifest

import (
	"fmt"

	"github.com/mishamyrt/hapm/internal/hapkg"
)

func ParseCategory(manifest map[string]any, key string) ([]hapkg.PackageDescription, error) {
	value, ok := manifest[key]
	if !ok {
		return nil, fmt.Errorf("key %s is not found in repo", key)
	}
	entries, ok := value.([]any)
	if !ok {
		return nil, fmt.Errorf("category %s must be a list", key)
	}
	items := make([]hapkg.PackageDescription, 0, len(entries))
	for _, entry := range entries {
		item, ok := entry.(string)
		if !ok {
			return nil, fmt.Errorf("wrong entity: %v", entry)
		}
		location, ok := ParseLocation(item)
		if !ok || location.FullName == "" {
			return nil, fmt.Errorf("wrong entity: %s", item)
		}
		items = append(items, hapkg.PackageDescription{
			FullName: location.FullName,
			Version:  location.Version,
			Kind:     key,
		})
	}
	return items, nil
}
