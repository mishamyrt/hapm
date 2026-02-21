package hapkg

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var versionRe = regexp.MustCompile(`^v?(\d+(?:(?:\.\d+)+)?)(\.?[0-9A-Za-z-\.]+)?$`)

type InvalidVersionError struct {
	msg string
}

func (e *InvalidVersionError) Error() string {
	return e.msg
}

type VersionParts struct {
	Value  []int
	Suffix []string
}

func ParseVersion(versionExpr string) (VersionParts, error) {
	if versionExpr == "" {
		return VersionParts{}, &InvalidVersionError{msg: "Version is empty"}
	}
	if strings.Contains(versionExpr, "..") {
		return VersionParts{}, &InvalidVersionError{msg: fmt.Sprintf("Repeated dots in version: %s", versionExpr)}
	}
	match := versionRe.FindStringSubmatch(versionExpr)
	if match == nil {
		return VersionParts{}, &InvalidVersionError{msg: fmt.Sprintf("Invalid version format: %s", versionExpr)}
	}
	values, err := parseSegments(match[1])
	if err != nil {
		return VersionParts{}, err
	}
	suffix, err := parseSuffix(match[2])
	if err != nil {
		return VersionParts{}, err
	}
	return VersionParts{Value: values, Suffix: suffix}, nil
}

func parseSegments(segmentValues string) ([]int, error) {
	segments := make([]int, 0)
	for _, segment := range strings.Split(segmentValues, ".") {
		if segment == "" {
			continue
		}
		if !regexp.MustCompile(`^\d+$`).MatchString(segment) {
			return nil, &InvalidVersionError{msg: fmt.Sprintf("Invalid segments value: %s", segmentValues)}
		}
		value, err := strconv.Atoi(segment)
		if err != nil || value < 0 {
			return nil, &InvalidVersionError{msg: fmt.Sprintf("Invalid segments value: %s", segmentValues)}
		}
		segments = append(segments, value)
	}
	return segments, nil
}

func parseSuffix(suffix string) ([]string, error) {
	if suffix == "" {
		return nil, nil
	}
	if suffix[0] == '-' || suffix[0] == '.' {
		if len(suffix) == 1 {
			return nil, &InvalidVersionError{msg: fmt.Sprintf("Invalid suffix: %s", suffix)}
		}
		suffix = suffix[1:]
	}
	if !strings.Contains(suffix, ".") {
		return []string{suffix}, nil
	}
	return strings.Split(suffix, "."), nil
}

type Version struct {
	Original string
	Value    []int
	Suffix   []string
}

func NewVersion(raw string) (Version, error) {
	parts, err := ParseVersion(raw)
	if err != nil {
		return Version{}, err
	}
	return Version{Original: raw, Value: parts.Value, Suffix: parts.Suffix}, nil
}

func MustNewVersion(raw string) Version {
	version, err := NewVersion(raw)
	if err != nil {
		panic(err)
	}
	return version
}

func (v Version) IsStable() bool {
	return v.Suffix == nil
}

func (v Version) Compare(other Version) int {
	if len(v.Value) != len(other.Value) {
		if len(v.Value) < len(other.Value) {
			return -1
		}
		return 1
	}
	for i := range v.Value {
		if v.Value[i] < other.Value[i] {
			return -1
		}
		if v.Value[i] > other.Value[i] {
			return 1
		}
	}
	if v.Suffix == nil && other.Suffix == nil {
		return 0
	}
	if v.Suffix == nil {
		return 1
	}
	if other.Suffix == nil {
		return -1
	}
	if len(v.Suffix) != len(other.Suffix) {
		if len(v.Suffix) < len(other.Suffix) {
			return -1
		}
		return 1
	}
	left := strings.Join(v.Suffix, ".")
	right := strings.Join(other.Suffix, ".")
	if left < right {
		return -1
	}
	if left > right {
		return 1
	}
	return 0
}

func FindLatestVersion(tags []string, stableOnly bool) string {
	latest, _ := NewVersion("0.0.0")
	versions := make([]Version, 0, len(tags))
	for _, tag := range tags {
		parsed, err := NewVersion(tag)
		if err != nil {
			continue
		}
		versions = append(versions, parsed)
	}
	for _, candidate := range versions {
		if stableOnly && !candidate.IsStable() {
			continue
		}
		if candidate.Compare(latest) > 0 {
			latest = candidate
		}
	}
	return latest.Original
}
