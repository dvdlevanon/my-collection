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

		fullpath := filepath.Join(path, file.Name())
		_, err := os.Stat(fullpath)
		if err != nil {
			// ignore files that we cant read their stat (broken symlink e.g.)
			continue
		}

		if file.IsDir() {
			child, err := buildFromDir(ctx, node, fullpath, filter)
			if err != nil {
				return nil, err
			}

			node.Children = append(node.Children, child)
		} else if filter.Filter(ctx, fullpath) {
			node.Files = append(node.Files, createFileNode(node, file.Name()))
		}
	}

	return node, nil
}
