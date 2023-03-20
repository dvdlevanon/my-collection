package relativasor

import (
	"path/filepath"
	"strings"
)

var rootDirectory string

func Init(r string) {
	rootDirectory = r
}

func GetRelativePath(url string) string {
	if url == "." {
		return ""
	}

	if !strings.HasPrefix(url, rootDirectory) {
		return url
	}

	relativePath := strings.TrimPrefix(url, rootDirectory)
	return strings.TrimPrefix(relativePath, string(filepath.Separator))
}

func GetAbsoluteFile(url string) string {
	if strings.HasPrefix(url, string(filepath.Separator)) {
		return url
	} else {
		return filepath.Join(rootDirectory, url)
	}
}

func GetRootDirectory() string {
	return rootDirectory
}
