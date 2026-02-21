package repository

import "strings"

func RepoName(fullName string) string {
	parts := strings.Split(fullName, "/")
	if len(parts) == 0 {
		return ""
	}
	return parts[len(parts)-1]
}
