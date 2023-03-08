package relativasor

import (
	"path/filepath"
	"strings"
)

func New(rootDirectory string) *PathRelativasor {
	return &PathRelativasor{
		rootDirectory: rootDirectory,
	}
}

type PathRelativasor struct {
	rootDirectory string
}

func (g *PathRelativasor) GetRelativePath(url string) string {
	if !strings.HasPrefix(url, g.rootDirectory) {
		return url
	}

	if url == g.rootDirectory {
		return url
	}

	relativePath := strings.TrimPrefix(url, g.rootDirectory)
	return strings.TrimPrefix(relativePath, string(filepath.Separator))
}

func (g *PathRelativasor) GetAbsoluteFile(url string) string {
	if strings.HasPrefix(url, string(filepath.Separator)) {
		return url
	} else {
		return filepath.Join(g.rootDirectory, url)
	}
}
