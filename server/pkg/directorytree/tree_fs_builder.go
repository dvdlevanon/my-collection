package directorytree

import (
	"os"
	"path/filepath"
)

type FilesFilter func(path string) bool

func BuildFromPath(path string, filter FilesFilter) (*DirectoryNode, error) {
	root, err := buildFromDir(nil, path, filter)
	if err != nil {
		return nil, err
	}
	return root, nil
}

func buildFromDir(parent *DirectoryNode, path string, filter FilesFilter) (*DirectoryNode, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	node := createDirectoryNode(parent, filepath.Base(path))
	if parent == nil {
		node.Title = ""
	}

	for _, file := range files {
		if file.IsDir() {
			child, err := buildFromDir(node, filepath.Join(path, file.Name()), filter)
			if err != nil {
				return nil, err
			}

			node.Children = append(node.Children, child)
		} else if filter(filepath.Join(path, file.Name())) {
			node.Files = append(node.Files, createFileNode(node, file.Name()))
		}
	}

	return node, nil
}
