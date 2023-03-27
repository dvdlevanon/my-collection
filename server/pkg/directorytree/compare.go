package directorytree

import (
	"path/filepath"
)

func Compare(fs *DirectoryNode, db *DirectoryNode) *Diff {
	rawChanges := newIndexedChanges()
	compareDirectory(fs, db, rawChanges)
	diff := detectMoves(rawChanges)
	rawChanges.addToDiff(diff)
	removeExcluded(db, diff)
	return diff
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

		dbOnly = append(dbOnly, dbFile)
	}

	return fsOnly, dbOnly
}

func detectMoves(rawChanges *indexedChanges) *Diff {
	diff := newDiff()
	changed := false

	for filename, l := range rawChanges.removedDirs {
		for _, p := range l {
			changed = detectDirectoryMove(diff, rawChanges, p)
			if changed {
				rawChanges.removedDirs[filename] = make([]string, 0)
			}
		}
	}

	for filename, l := range rawChanges.removedFiles {
		for _, p := range l {
			changed = detectFileMove(diff, rawChanges, p)
			if changed {
				rawChanges.removedFiles[filename] = make([]string, 0)
			}
		}
	}

	return diff
}

func detectDirectoryMove(diff *Diff, rawChanges *indexedChanges, path string) bool {
	filename := filepath.Base(path)
	l, ok := rawChanges.addedDirs[filename]
	if !ok || len(l) == 0 {
		return false
	}

	rawChanges.addedDirs[filename] = make([]string, 0)
	diff.MovedDirectories = append(diff.MovedDirectories, Change{Path1: path, Path2: l[0], ChangeType: DIRECTORY_MOVED})
	return true
}

func detectFileMove(diff *Diff, rawChanges *indexedChanges, path string) bool {
	filename := filepath.Base(path)
	l, ok := rawChanges.addedFiles[filename]
	if !ok || len(l) == 0 {
		return false
	}

	rawChanges.addedFiles[filename] = make([]string, 0)
	diff.MovedFiles = append(diff.MovedFiles, Change{Path1: path, Path2: l[0], ChangeType: FILE_MOVED})
	return true
}

func removeExcluded(db *DirectoryNode, diff *Diff) {
	diff.AddedDirectories = removeExcludedFromChanges(db, diff.AddedDirectories)
	diff.AddedFiles = removeExcludedFromChanges(db, diff.AddedFiles)
}

func removeExcludedFromChanges(db *DirectoryNode, changes []Change) []Change {
	counter := 0
	for {
		if counter >= len(changes) {
			return changes
		}

		change := changes[counter]
		if db.isExcluded(change.Path1) {
			changes = append(changes[:counter], changes[counter+1:]...)
			continue
		}

		counter = counter + 1
	}
}
