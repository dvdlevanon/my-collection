package directorytree

import (
	"my-collection/server/pkg/bl/directories"
	"my-collection/server/pkg/model"
	"os"
	"strings"
)

func BuildFromDb(dr model.DirectoryReader, dig model.DirectoryItemsGetter) (*DirectoryNode, error) {
	dirs, err := dr.GetAllDirectories()
	if err != nil {
		return nil, err
	}

	root := createDirectoryNode(nil, "")
	for _, dir := range *dirs {
		path := dir.Path
		if path == model.ROOT_DIRECTORY_PATH {
			path = ""
		}

		child := root.getOrCreateChild(path)
		child.Excluded = directories.IsExcluded(&dir)

		if err := child.readFilesFromDb(dig); err != nil {
			logger.Errorf("Error reading files from db %s", err)
		}
	}

	return root, nil
}

func (dn *DirectoryNode) getOrCreateChild(path string) *DirectoryNode {
	if path == "" {
		return dn
	}

	parts := strings.SplitN(path, string(os.PathSeparator), 2)
	firstDir := parts[0]
	remainingDirs := ""
	if len(parts) > 1 {
		remainingDirs = parts[1]
	}

	for _, child := range dn.Children {
		if child.Title == firstDir {
			return child.getOrCreateChild(remainingDirs)
		}
	}

	child := createDirectoryNode(dn, firstDir)
	dn.Children = append(dn.Children, child)
	return child.getOrCreateChild(remainingDirs)
}

func (dn *DirectoryNode) readFilesFromDb(dig model.DirectoryItemsGetter) error {
	items, err := dig.GetBelongingItems(dn.getPath())
	if err != nil {
		return err
	}

	for _, item := range *items {
		dn.Files = append(dn.Files, createFileNode(dn, item.Title))
	}

	return nil
}
