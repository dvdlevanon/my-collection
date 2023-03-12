package directorytree

import (
	"path/filepath"

	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("fswatch")

type DirectoryNode struct {
	Parent   *DirectoryNode
	Children []*DirectoryNode
	Files    []*FileNode
	Title    string
	Excluded bool
}

type FileNode struct {
	Parent *DirectoryNode
	Title  string
}

func createFileNode(parent *DirectoryNode, title string) *FileNode {
	return &FileNode{
		Parent: parent,
		Title:  title,
	}
}

func createDirectoryNode(parent *DirectoryNode, title string) *DirectoryNode {
	return &DirectoryNode{
		Parent:   parent,
		Title:    title,
		Children: make([]*DirectoryNode, 0),
		Files:    make([]*FileNode, 0),
	}
}

func (dn *DirectoryNode) getPath() string {
	if dn.Parent == nil {
		return dn.Title
	} else {
		return filepath.Join(dn.Parent.getPath(), dn.Title)
	}
}

func (fn *FileNode) getPath() string {
	if fn.Parent == nil {
		return fn.Title
	} else {
		return filepath.Join(fn.Parent.getPath(), fn.Title)
	}
}
