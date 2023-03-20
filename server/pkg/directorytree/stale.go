package directorytree

import (
	"fmt"
	"strings"
)

type Stale struct {
	Dirs  []string
	Files []string
}

func (s *Stale) String() string {
	dirs := fmt.Sprintf("dirs: %s", strings.Join(s.Dirs, "\n\t"))
	files := fmt.Sprintf("files: %s", strings.Join(s.Files, "\n\t"))
	return fmt.Sprintf("%s\n%s", dirs, files)
}

func FindStales(db *DirectoryNode) *Stale {
	result := Stale{}
	findStale(db, &result)
	return &result
}

func findStale(node *DirectoryNode, result *Stale) {
	if node.Excluded {
		addToStale(node, result)
	} else {
		for _, child := range node.Children {
			findStale(child, result)
		}
	}
}

func addToStale(node *DirectoryNode, result *Stale) {
	if node.Parent != nil && node.getRoot().isExcluded(node.Parent.getPath()) {
		result.Dirs = append(result.Dirs, node.getPath())
	}

	for _, file := range node.Files {
		result.Files = append(result.Files, file.getPath())
	}

	for _, child := range node.Children {
		addToStale(child, result)
	}
}
