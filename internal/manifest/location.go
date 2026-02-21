package manifest

import (
	"net/url"
	"regexp"
	"strings"
)

type PackageLocation struct {
	FullName string
	Version  string
}

var (
	releaseTagRe  = regexp.MustCompile(`^/(.*)/(.*)/releases/tag/(.*)$`)
	packageNameRe = regexp.MustCompile(`^(.*)/(.[^@]*)(@.{1,})?$`)
)

func safeURLParse(raw string) (*url.URL, bool) {
	parsed, err := url.Parse(raw)
	if err != nil {
		return nil, false
	}
	return parsed, true
}

func parseGitHubURL(raw string) (*url.URL, bool) {
	parsed, ok := safeURLParse(raw)
	if !ok {
		return nil, false
	}
	if parsed.Scheme == "" || parsed.Host == "" {
		return nil, false
	}
	if parsed.Host != "github.com" {
		return nil, false
	}
	return parsed, true
}

func ParseLocationURL(raw string) (*PackageLocation, bool) {
	parsed, ok := parseGitHubURL(raw)
	if !ok {
		return nil, false
	}
	if match := releaseTagRe.FindStringSubmatch(parsed.Path); match != nil {
		return &PackageLocation{FullName: match[1] + "/" + match[2], Version: match[3]}, true
	}
	if strings.Count(parsed.Path, "/") != 2 {
		return nil, false
	}
	return &PackageLocation{FullName: strings.TrimPrefix(parsed.Path, "/"), Version: "latest"}, true
}

func ParsePackageName(location string) (*PackageLocation, bool) {
	result := packageNameRe.FindStringSubmatch(location)
	if result == nil {
		return nil, false
	}
	user := result[1]
	repository := result[2]
	if user == "" || repository == "" {
		return nil, false
	}
	version := result[3]
	if version == "" {
		version = "latest"
	} else {
		version = strings.TrimPrefix(version, "@")
	}
	return &PackageLocation{FullName: user + "/" + repository, Version: version}, true
}

func ParseLocation(pkg string) (*PackageLocation, bool) {
	if strings.HasPrefix(pkg, "github.com") {
		pkg = "https://" + pkg
	}
	parsers := []func(string) (*PackageLocation, bool){
		ParseLocationURL,
		ParsePackageName,
	}
	for _, parser := range parsers {
		if parsed, ok := parser(pkg); ok {
			return parsed, true
		}
	}
	return nil, false
}
