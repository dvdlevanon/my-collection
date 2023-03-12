package directorytree

import (
	"math"
	"path/filepath"
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

func Compare(fs *DirectoryNode, db *DirectoryNode) *Diff {
	rawChanges := compareDirectory(fs, db)
	changes := detectMoves(rawChanges)
	return &Diff{Changes: changes}
}

func detectMoves(changes []Change) []Change {
	for i, change := range changes {
		switch change.ChangeType {
		case DIRECTORY_REMOVED:
			changes = detectDirectoryMove(changes, i)
		case FILE_REMOVED:
			changes = detectFileMove(changes, i)
		}
	}

	return changes
}

func detectDirectoryMove(changes []Change, i int) []Change {
	for curI, cur := range changes {
		if cur.ChangeType != DIRECTORY_ADDED {
			continue
		}

		if filepath.Base(changes[i].Path1) != filepath.Base(cur.Path1) {
			continue
		}

		moveChange := Change{Path1: changes[i].Path1, Path2: changes[curI].Path1, ChangeType: DIRECTORY_MOVED}
		first := int(math.Max(float64(i), float64(curI)))
		second := int(math.Min(float64(i), float64(curI)))
		changes = append(changes[:first], changes[first+1:]...)
		changes = append(changes[:second], changes[second+1:]...)
		changes = append(changes, moveChange)
	}

	return changes
}

func detectFileMove(changes []Change, i int) []Change {
	for curI, cur := range changes {
		if cur.ChangeType != FILE_ADDED {
			continue
		}

		if filepath.Base(changes[i].Path1) != filepath.Base(cur.Path1) {
			continue
		}

		moveChange := Change{Path1: changes[i].Path1, Path2: changes[curI].Path1, ChangeType: FILE_MOVED}
		changes = append(changes[:i], changes[i+1:]...)
		changes = append(changes[:curI], changes[curI+1:]...)
		changes = append(changes, moveChange)
	}

	return changes
}

func compareDirectory(fs *DirectoryNode, db *DirectoryNode) []Change {
	changes := make([]Change, 0)

	// if db.Excluded {
	// 	return changes
	// }

	bothDirs, fsDirsOnly, dbDirsOnly := compareSubDirectories(fs.Children, db)
	changes = append(changes, createDirsChanges(fsDirsOnly, dbDirsOnly)...)

	fsFilesOnly, dbFilesOnly := compareFiles(fs.Files, db)
	changes = append(changes, createFilesChanges(fsFilesOnly, dbFilesOnly)...)

	for _, dir := range bothDirs {
		changes = append(changes, compareDirectory(dir[0], dir[1])...)
	}

	for _, dir := range fsDirsOnly {
		changes = append(changes, compareDirectory(dir, nil)...)
	}

	return changes
}

func createDirsChanges(fsDirsOnly []*DirectoryNode, dbDirsOnly []*DirectoryNode) []Change {
	changes := make([]Change, 0)

	for _, dir := range fsDirsOnly {
		changes = append(changes, Change{Path1: dir.getPath(), ChangeType: DIRECTORY_ADDED})
	}

	for _, dir := range dbDirsOnly {
		changes = append(changes, Change{Path1: dir.getPath(), ChangeType: DIRECTORY_REMOVED})
	}

	return changes
}

func createFilesChanges(fsFilesOnly []*FileNode, dbFilesOnly []*FileNode) []Change {
	changes := make([]Change, 0)

	for _, file := range fsFilesOnly {
		changes = append(changes, Change{Path1: file.getPath(), ChangeType: FILE_ADDED})
	}

	for _, file := range dbFilesOnly {
		changes = append(changes, Change{Path1: file.getPath(), ChangeType: FILE_REMOVED})
	}

	return changes
}

func compareSubDirectories(fs []*DirectoryNode, db *DirectoryNode) ([][]*DirectoryNode, []*DirectoryNode, []*DirectoryNode) {
	both := make([][]*DirectoryNode, 0)
	fsOnly := make([]*DirectoryNode, 0)
	dbOnly := make([]*DirectoryNode, 0)
	if db == nil {
		return both, fs, dbOnly
	}

outFs:
	for _, fsDir := range fs {
		for _, dbDir := range db.Children {
			if fsDir.Title == dbDir.Title {
				both = append(both, []*DirectoryNode{fsDir, dbDir})
				continue outFs
			}
		}

		fsOnly = append(fsOnly, fsDir)
	}

outDb:
	for _, dbDir := range db.Children {
		for _, fsDir := range fs {
			if dbDir.Title == fsDir.Title {
				continue outDb
			}
		}

		dbOnly = append(dbOnly, dbDir)
	}

	return both, fsOnly, dbOnly
}

func compareFiles(fs []*FileNode, db *DirectoryNode) ([]*FileNode, []*FileNode) {
	fsOnly := make([]*FileNode, 0)
	dbOnly := make([]*FileNode, 0)
	if db == nil {
		return fs, dbOnly
	}

outFs:
	for _, fsFile := range fs {
		for _, dbFile := range db.Files {
			if fsFile.Title == dbFile.Title {
				continue outFs
			}
		}

		fsOnly = append(fsOnly, fsFile)
	}

outDb:
	for _, dbFile := range db.Files {
		for _, fsFile := range fs {
			if dbFile.Title == fsFile.Title {
				continue outDb
			}
		}

		dbOnly = append(fsOnly, dbFile)
	}

	return fsOnly, dbOnly
}
