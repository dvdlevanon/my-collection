package directorytree

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("tree")

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

func (dn *DirectoryNode) isExcluded(path string) bool {
	if path == "" || dn.Excluded {
		return true
	}

	parts := strings.SplitN(path, string(os.PathSeparator), 2)
	firstDir := parts[0]
	remainingDirs := ""
	if len(parts) > 1 {
		remainingDirs = parts[1]
	}

	for _, child := range dn.Children {
		if child.Title == firstDir {
			return child.isExcluded(remainingDirs)
		}
	}

	return dn.Excluded
}

func (dn *DirectoryNode) String(depth int) string {
	var result strings.Builder
	result.WriteString(fmt.Sprintf("%s%s", strings.Repeat("  ", depth), dn.Title))
	result.WriteByte('\n')

	sort.SliceStable(dn.Files, func(i, j int) bool { return dn.Files[i].Title < dn.Files[j].Title })
	for _, file := range dn.Files {
		result.WriteString(file.String(depth + 1))
		result.WriteByte('\n')
	}

	sort.SliceStable(dn.Children, func(i, j int) bool { return dn.Children[i].Title < dn.Children[j].Title })
	for _, child := range dn.Children {
		result.WriteString(child.String(depth + 1))
		result.WriteByte('\n')
	}

	return result.String()
}

func (fn *FileNode) String(depth int) string {
	return fmt.Sprintf("%s%s", strings.Repeat("  ", depth), fn.Title)
}
