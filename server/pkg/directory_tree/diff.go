package directorytree

import (
	"fmt"
	"path/filepath"
	"strings"
)

type Diff struct {
	Changes []Change
}

type ChangeType int

const (
	DIRECTORY_ADDED = iota
	DIRECTORY_REMOVED
	DIRECTORY_MOVED
	FILE_ADDED
	FILE_REMOVED
	FILE_MOVED
)

type Change struct {
	Path1      string
	Path2      string
	ChangeType ChangeType
}

func (c *Change) String() string {
	switch c.ChangeType {
	case DIRECTORY_ADDED:
		return fmt.Sprintf("+ %s", c.Path1)
	case DIRECTORY_REMOVED:
		return fmt.Sprintf("- %s", c.Path1)
	case DIRECTORY_MOVED:
		return fmt.Sprintf("%s -> %s", c.Path1, c.Path2)
	case FILE_ADDED:
		return fmt.Sprintf("+ %s", c.Path1)
	case FILE_REMOVED:
		return fmt.Sprintf("- %s", c.Path1)
	case FILE_MOVED:
		return fmt.Sprintf("%s -> %s", c.Path1, c.Path2)
	default:
		return "unknown"
	}
}

func (d *Diff) String() string {
	result := strings.Builder{}

	for _, change := range d.Changes {
		result.WriteString(change.String())
		result.WriteByte('\n')
	}

	return result.String()
}

func newIndexedChanges() *indexedChanges {
	return &indexedChanges{
		addedDirs:    make(map[string][]string),
		addedFiles:   make(map[string][]string),
		removedDirs:  make(map[string][]string),
		removedFiles: make(map[string][]string),
	}
}

type indexedChanges struct {
	addedDirs    map[string][]string
	addedFiles   map[string][]string
	removedDirs  map[string][]string
	removedFiles map[string][]string
}

func (i *indexedChanges) dirAdded(path string) {
	filename := filepath.Base(path)
	i.addedDirs[filename] = append(i.addedDirs[filename], path)
}

func (i *indexedChanges) dirRemoved(path string) {
	filename := filepath.Base(path)
	i.removedDirs[filename] = append(i.removedDirs[filename], path)
}

func (i *indexedChanges) fileAdded(path string) {
	filename := filepath.Base(path)
	i.addedFiles[filename] = append(i.addedFiles[filename], path)
}

func (i *indexedChanges) fileRemoved(path string) {
	filename := filepath.Base(path)
	i.removedFiles[filename] = append(i.removedFiles[filename], path)
}

func (i *indexedChanges) toChanges() []Change {
	changes := i.mapToChanges(DIRECTORY_ADDED, i.addedDirs)
	changes = append(changes, i.mapToChanges(DIRECTORY_REMOVED, i.removedDirs)...)
	changes = append(changes, i.mapToChanges(FILE_ADDED, i.addedFiles)...)
	changes = append(changes, i.mapToChanges(FILE_REMOVED, i.removedFiles)...)
	return changes
}

func (i *indexedChanges) mapToChanges(ct ChangeType, m map[string][]string) []Change {
	changes := make([]Change, 0)
	for _, l := range m {
		for _, f := range l {
			changes = append(changes, Change{ChangeType: ct, Path1: f})
		}
	}

	return changes
}
