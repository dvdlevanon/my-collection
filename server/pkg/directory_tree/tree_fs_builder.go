package directorytree

import (
	"os"
	"path/filepath"
)

func BuildFromPath(path string) (*Tree, error) {
	root, err := buildFromDir(nil, path)
	if err != nil {
		return nil, err
	}
	return &Tree{Root: root}, nil
}

func buildFromDir(parent *DirectoryNode, path string) (*DirectoryNode, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	node := createDirectoryNode(parent, filepath.Base(path))
	for _, file := range files {
		if file.IsDir() {
			child, err := buildFromDir(node, filepath.Join(path, file.Name()))
			if err != nil {
				return nil, err
			}

			node.Children = append(node.Children, child)
		} else {
			node.Files = append(node.Files, createFileNode(node, file.Name()))
		}
	}

	return node, nil
}
