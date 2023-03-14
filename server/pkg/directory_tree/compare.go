package directorytree

import (
	"path/filepath"
)

func Compare(fs *DirectoryNode, db *DirectoryNode) *Diff {
	rawChanges := newIndexedChanges()
	compareDirectory(fs, db, rawChanges)
	changes := detectMoves(rawChanges)
	changes = append(changes, rawChanges.toChanges()...)
	changes = removeExcluded(db, changes)
	return &Diff{Changes: changes}
}

func compareDirectory(fs *DirectoryNode, db *DirectoryNode, changes *indexedChanges) {
	bothDirs, fsDirsOnly, dbDirsOnly := compareSubDirectories(fs, db)
	createDirsChanges(fsDirsOnly, dbDirsOnly, changes)

	fsFilesOnly, dbFilesOnly := compareFiles(fs, db)
	createFilesChanges(fsFilesOnly, dbFilesOnly, changes)

	for _, dir := range bothDirs {
		compareDirectory(dir[0], dir[1], changes)
	}

	for _, dir := range fsDirsOnly {
		compareDirectory(dir, nil, changes)
	}

	for _, dir := range dbDirsOnly {
		compareDirectory(nil, dir, changes)
	}
}

func createDirsChanges(fsDirsOnly []*DirectoryNode, dbDirsOnly []*DirectoryNode, changes *indexedChanges) {
	for _, dir := range fsDirsOnly {
		changes.dirAdded(dir.getPath())
	}

	for _, dir := range dbDirsOnly {
		changes.dirRemoved(dir.getPath())
	}
}

func createFilesChanges(fsFilesOnly []*FileNode, dbFilesOnly []*FileNode, changes *indexedChanges) {
	for _, file := range fsFilesOnly {
		changes.fileAdded(file.getPath())
	}

	for _, file := range dbFilesOnly {
		changes.fileRemoved(file.getPath())
	}
}

func compareSubDirectories(fs *DirectoryNode, db *DirectoryNode) ([][]*DirectoryNode, []*DirectoryNode, []*DirectoryNode) {
	both := make([][]*DirectoryNode, 0)
	fsOnly := make([]*DirectoryNode, 0)
	dbOnly := make([]*DirectoryNode, 0)
	if db == nil {
		return both, fs.Children, dbOnly
	}

	if fs == nil {
		return both, fsOnly, db.Children
	}

outFs:
	for _, fsDir := range fs.Children {
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
		for _, fsDir := range fs.Children {
			if dbDir.Title == fsDir.Title {
				continue outDb
			}
		}

		dbOnly = append(dbOnly, dbDir)
	}

	return both, fsOnly, dbOnly
}

func compareFiles(fs *DirectoryNode, db *DirectoryNode) ([]*FileNode, []*FileNode) {
	fsOnly := make([]*FileNode, 0)
	dbOnly := make([]*FileNode, 0)
	if db == nil {
		return fs.Files, dbOnly
	}
	if fs == nil {
		return fsOnly, db.Files
	}

outFs:
	for _, fsFile := range fs.Files {
		for _, dbFile := range db.Files {
			if fsFile.Title == dbFile.Title {
				continue outFs
			}
		}

		fsOnly = append(fsOnly, fsFile)
	}

outDb:
	for _, dbFile := range db.Files {
		for _, fsFile := range fs.Files {
			if dbFile.Title == fsFile.Title {
				continue outDb
			}
		}

		dbOnly = append(fsOnly, dbFile)
	}

	return fsOnly, dbOnly
}

func detectMoves(rawChanges *indexedChanges) []Change {
	changes := make([]Change, 0)
	changed := false

	for filename, l := range rawChanges.removedDirs {
		for _, p := range l {
			changes, changed = detectDirectoryMove(changes, rawChanges, p)
			if changed {
				rawChanges.removedDirs[filename] = make([]string, 0)
			}
		}
	}

	for filename, l := range rawChanges.removedFiles {
		for _, p := range l {
			changes, changed = detectFileMove(changes, rawChanges, p)
			if changed {
				rawChanges.removedFiles[filename] = make([]string, 0)
			}
		}
	}

	return changes
}

func detectDirectoryMove(changes []Change, rawChanges *indexedChanges, path string) ([]Change, bool) {
	filename := filepath.Base(path)
	l, ok := rawChanges.addedDirs[filename]
	if !ok || len(l) == 0 {
		return changes, false
	}

	rawChanges.addedDirs[filename] = make([]string, 0)
	return append(changes, Change{Path1: path, Path2: l[0], ChangeType: DIRECTORY_MOVED}), true
}

func detectFileMove(changes []Change, rawChanges *indexedChanges, path string) ([]Change, bool) {
	filename := filepath.Base(path)
	l, ok := rawChanges.addedFiles[filename]
	if !ok || len(l) == 0 {
		return changes, false
	}

	rawChanges.addedFiles[filename] = make([]string, 0)
	return append(changes, Change{Path1: path, Path2: l[0], ChangeType: FILE_MOVED}), true
}

func removeExcluded(db *DirectoryNode, changes []Change) []Change {
	counter := 0
	for {
		if counter >= len(changes) {
			return changes
		}

		change := changes[counter]
		if change.ChangeType != DIRECTORY_ADDED && change.ChangeType != FILE_ADDED {
			counter = counter + 1
			continue
		}

		if db.isExcluded(change.Path1) {
			changes = append(changes[:counter], changes[counter+1:]...)
			continue
		}

		counter = counter + 1
	}
}
