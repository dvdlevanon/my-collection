package directorytree

import (
	"context"
	"os"
	"path/filepath"
	"strings"
)

type FilesFilter interface {
	Filter(ctx context.Context, path string) bool
}

func BuildFromPath(ctx context.Context, path string, filter FilesFilter) (*DirectoryNode, error) {
	root, err := buildFromDir(ctx, nil, path, filter)
	if err != nil {
		return nil, err
	}
	return root, nil
}

func buildFromDir(ctx context.Context, parent *DirectoryNode, path string, filter FilesFilter) (*DirectoryNode, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	node := createDirectoryNode(parent, filepath.Base(path))
	if parent == nil {
		node.Title = ""
	}

	for _, file := range files {
		if strings.HasPrefix(file.Name(), ".") {
			// ignore hidden files and dirs
			continue
		}

		if file.IsDir() {
			child, err := buildFromDir(ctx, node, filepath.Join(path, file.Name()), filter)
			if err != nil {
				return nil, err
			}

			node.Children = append(node.Children, child)
		} else if filter.Filter(ctx, filepath.Join(path, file.Name())) {
			node.Files = append(node.Files, createFileNode(node, file.Name()))
		}
	}

	return node, nil
}
