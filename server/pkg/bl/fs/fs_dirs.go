package fs

import (
	"context"
	"errors"
	"my-collection/server/pkg/bl/directories"
	"my-collection/server/pkg/model"
	"my-collection/server/pkg/relativasor"
	"os"
	"path/filepath"
	"strings"

	"gorm.io/gorm"
)

func GetFsTree(path string, depth int) (*model.FsNode, error) {
	if depth < 0 {
		return nil, nil
	}

	if path == model.ROOT_DIRECTORY_PATH {
		path = relativasor.GetRootDirectory()
	}

	if strings.HasPrefix(filepath.Base(path), ".") {
		// ignore hidden files and dirs
		return nil, nil
	}

	fi, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	node := &model.FsNode{
		Path: path,
		Type: getNodeType(fi),
	}

	if !fi.IsDir() {
		return node, nil
	}

	children, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, child := range children {
		childNode, err := GetFsTree(filepath.Join(path, child.Name()), depth-1)
		if err != nil {
			return nil, err
		}
		if childNode == nil {
			continue
		}
		node.Children = append(node.Children, childNode)
	}

	return node, nil
}

func getNodeType(fi os.FileInfo) model.FsNodeType {
	if fi.IsDir() {
		return model.FS_NODE_DIR
	} else {
		return model.FS_NODE_FILE
	}
}

func includeDirWithParents(ctx context.Context, drw model.DirectoryReaderWriter, path string) error {
	if directories.NormalizeDirectoryPath(path) == model.ROOT_DIRECTORY_PATH {
		return nil
	}

	if err := includeDirWithParents(ctx, drw, filepath.Dir(path)); err != nil {
		return err
	}

	err := directories.IncludeOrCreateDirectory(ctx, drw, path)
	if err != nil {
		return err
	}

	return nil
}

func includeHierarchy(ctx context.Context, drw model.DirectoryReaderWriter, path string, depth int) error {
	if depth == 0 {
		return nil
	}

	dirs, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for _, dir := range dirs {
		subdir := filepath.Join(path, dir.Name())
		info, err := os.Stat(subdir)
		if err != nil {
			return err
		}

		if info.IsDir() {
			if err := directories.IncludeOrCreateDirectory(ctx, drw, subdir); err != nil {
				return err
			}

			if err := includeHierarchy(ctx, drw, subdir, depth-1); err != nil {
				return err
			}
		}
	}

	return nil
}

func excludeHierarchy(ctx context.Context, drw model.DirectoryReaderWriter, path string, depth int) error {
	if depth == 0 {
		return nil
	}

	dirs, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for _, dir := range dirs {
		subdir := filepath.Join(path, dir.Name())
		info, err := os.Stat(subdir)
		if err != nil {
			return err
		}

		if info.IsDir() {
			dir, err := directories.GetDirectory(ctx, drw, subdir)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					continue
				}
				return err
			}

			if directories.IsExcluded(dir) {
				continue
			}

			if err := excludeHierarchy(ctx, drw, subdir, depth-1); err != nil {
				return err
			}

			if err := directories.ExcludeDirectory(ctx, drw, subdir); err != nil {
				return err
			}
		}
	}

	return nil
}

func ExcludeDir(ctx context.Context, drw model.DirectoryReaderWriter, path string) error {
	if err := excludeHierarchy(ctx, drw, path, -1); err != nil {
		return err
	}

	return directories.ExcludeDirectory(ctx, drw, path)
}

func IncludeDir(ctx context.Context, drw model.DirectoryReaderWriter, path string, subdirs bool, hierarchy bool) error {
	if err := includeDirWithParents(ctx, drw, path); err != nil {
		return err
	}

	if subdirs {
		if err := directories.AutoIncludeChildren(ctx, drw, path); err != nil {
			return err
		}
		if err := includeHierarchy(ctx, drw, path, 1); err != nil {
			return err
		}
	}

	if hierarchy {
		if err := directories.AutoIncludeHierarchy(ctx, drw, path); err != nil {
			return err
		}
		if err := includeHierarchy(ctx, drw, path, -1); err != nil {
			return err
		}
	}

	return nil
}
