package fs

import (
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

func includeDirWithParents(drw model.DirectoryReaderWriter, path string) error {
	if directories.NormalizeDirectoryPath(path) == model.ROOT_DIRECTORY_PATH {
		return nil
	}

	if err := includeDirWithParents(drw, filepath.Dir(path)); err != nil {
		return err
	}

	err := directories.IncludeOrCreateDirectory(drw, path)
	if err != nil {
		return err
	}

	return nil
}

func includeHierarchy(drw model.DirectoryReaderWriter, path string, depth int) error {
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
			if err := directories.IncludeOrCreateDirectory(drw, subdir); err != nil {
				return err
			}

			if err := includeHierarchy(drw, subdir, depth-1); err != nil {
				return err
			}
		}
	}

	return nil
}

func excludeHierarchy(drw model.DirectoryReaderWriter, path string, depth int) error {
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
			dir, err := directories.GetDirectory(drw, subdir)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					continue
				}
				return err
			}

			if directories.IsExcluded(dir) {
				continue
			}

			if err := excludeHierarchy(drw, subdir, depth-1); err != nil {
				return err
			}

			if err := directories.ExcludeDirectory(drw, subdir); err != nil {
				return err
			}
		}
	}

	return nil
}

func ExcludeDir(drw model.DirectoryReaderWriter, path string) error {
	if err := excludeHierarchy(drw, path, -1); err != nil {
		return err
	}

	return directories.ExcludeDirectory(drw, path)
}

func IncludeDir(drw model.DirectoryReaderWriter, path string, subdirs bool, hierarchy bool) error {
	if err := includeDirWithParents(drw, path); err != nil {
		return err
	}

	if subdirs {
		if err := includeHierarchy(drw, path, 1); err != nil {
			return err
		}
	}

	if hierarchy {
		if err := includeHierarchy(drw, path, -1); err != nil {
			return err
		}
	}

	return nil
}
