package directorytree

import (
	"fmt"
	"path/filepath"
	"strings"
)

func newDiff() *Diff {
	return &Diff{
		AddedDirectories:   make([]Change, 0),
		RemovedDirectories: make([]Change, 0),
		AddedFiles:         make([]Change, 0),
		RemovedFiles:       make([]Change, 0),
		MovedDirectories:   make([]Change, 0),
		MovedFiles:         make([]Change, 0),
	}
}

type Diff struct {
	AddedDirectories   []Change
	RemovedDirectories []Change
	AddedFiles         []Change
	RemovedFiles       []Change
	MovedDirectories   []Change
	MovedFiles         []Change
}

func (s *Diff) HasChanges() bool {
	return len(s.AddedDirectories) > 0 ||
		len(s.RemovedDirectories) > 0 ||
		len(s.AddedFiles) > 0 ||
		len(s.RemovedFiles) > 0 ||
		len(s.MovedDirectories) > 0 ||
		len(s.MovedFiles) > 0
}

func (d *Diff) ChangesTotal() int {
	return len(d.AddedDirectories) +
		len(d.RemovedDirectories) +
		len(d.AddedFiles) +
		len(d.RemovedFiles) +
		len(d.MovedDirectories) +
		len(d.MovedFiles)
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
	result := make([]string, 0)
	result = append(result, d.ChangesToString(d.AddedDirectories)...)
	result = append(result, d.ChangesToString(d.RemovedDirectories)...)
	result = append(result, d.ChangesToString(d.AddedFiles)...)
	result = append(result, d.ChangesToString(d.RemovedFiles)...)
	result = append(result, d.ChangesToString(d.MovedDirectories)...)
	result = append(result, d.ChangesToString(d.MovedFiles)...)
	return strings.Join(result, "\n")
}

func (d *Diff) ChangesToString(changes []Change) []string {
	strs := make([]string, len(changes))
	for i, v := range changes {
		strs[i] = v.String()
	}

	return strs
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

func (i *indexedChanges) addToDiff(diff *Diff) {
	diff.AddedDirectories = append(diff.AddedDirectories, i.mapToChanges(DIRECTORY_ADDED, i.addedDirs)...)
	diff.RemovedDirectories = append(diff.RemovedDirectories, i.mapToChanges(DIRECTORY_REMOVED, i.removedDirs)...)
	diff.AddedFiles = append(diff.AddedFiles, i.mapToChanges(FILE_ADDED, i.addedFiles)...)
	diff.RemovedFiles = append(diff.RemovedFiles, i.mapToChanges(FILE_REMOVED, i.removedFiles)...)
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
