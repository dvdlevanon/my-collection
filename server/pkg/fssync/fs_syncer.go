package fssync

import (
	"my-collection/server/pkg/bl/directories"
	"my-collection/server/pkg/bl/items"
	"my-collection/server/pkg/bl/tags"
	"my-collection/server/pkg/db"
	"my-collection/server/pkg/directorytree"
	"my-collection/server/pkg/model"
	"my-collection/server/pkg/relativasor"
	"path/filepath"
	"strings"

	"github.com/go-errors/errors"
	"k8s.io/utils/pointer"
)

func newFsSyncer(path string, db *db.Database, dig model.DirectoryItemsGetter, filter directorytree.FilesFilter) (*fsSyncer, error) {
	exists, err := directories.DirectoryExists(db, path)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.Errorf("directory not found in db %s", path)
	}

	rootfs, err := directorytree.BuildFromPath(path, filter)
	if err != nil {
		return nil, err
	}

	rootdb, err := directorytree.BuildFromDb(db, dig)
	if err != nil {
		return nil, err
	}

	diff := directorytree.Compare(rootfs, rootdb)
	stales := directorytree.FindStales(rootdb)

	return &fsSyncer{
		diff:   diff,
		stales: stales,
	}, nil
}

type fsSyncer struct {
	stales *directorytree.Stale
	diff   *directorytree.Diff
}

func (f *fsSyncer) hasFsChanges() bool {
	return f.stales.HasChanges() || f.diff.HasChanges()
}

func (f *fsSyncer) sync(db *db.Database, digs model.DirectoryItemsGetterSetter,
	dctg model.DirectoryConcreteTagsGetter, flmg model.FileLastModifiedGetter) []error {
	errors := make([]error, 0)
	if f.hasFsChanges() {
		f.debugPrint()
		errors = append(errors, addMissingDirectoryTags(db, db)...)
		errors = append(errors, removeStaleItems(digs, db, f.stales.Files)...)
		errors = append(errors, removeStaleDirs(db, db, f.stales.Dirs)...)
		errors = append(errors, addMissingDirs(db, f.diff.AddedDirectories)...)
		errors = append(errors, removeDeletedDirs(db, db, f.diff.RemovedDirectories)...)
		errors = append(errors, removeDeletedFiles(digs, db, f.diff.RemovedFiles)...)
		errors = append(errors, f.renameDirs()...)
		errors = append(errors, renameFiles(db, db, db, f.diff.MovedFiles)...)
		errors = append(errors, addNewFiles(db, digs, dctg, flmg, f.diff.AddedFiles)...)
	}
	// syncConcreteTags()
	return errors
}

func (f *fsSyncer) debugPrint() {
	if len(f.stales.Dirs) > 0 {
		logger.Debugf("Stale directories %d - [%s]", len(f.stales.Dirs), strings.Join(f.stales.Dirs, ", "))
	}

	if len(f.stales.Files) > 0 {
		logger.Debugf("Stale files %d - [%s]", len(f.stales.Files), strings.Join(f.stales.Files, ", "))
	}

	if len(f.diff.AddedDirectories) > 0 {
		logger.Debugf("Added directories %d - [%s]", len(f.diff.AddedDirectories), strings.Join(f.diff.ChangesToString(f.diff.AddedDirectories), ", "))
	}

	if len(f.diff.RemovedDirectories) > 0 {
		logger.Debugf("Removed directories %d - [%s]", len(f.diff.RemovedDirectories), strings.Join(f.diff.ChangesToString(f.diff.RemovedDirectories), ", "))
	}

	if len(f.diff.AddedFiles) > 0 {
		logger.Debugf("Added files %d - [%s]", len(f.diff.AddedFiles), strings.Join(f.diff.ChangesToString(f.diff.AddedFiles), ", "))
	}

	if len(f.diff.RemovedFiles) > 0 {
		logger.Debugf("Removed files %d - [%s]", len(f.diff.RemovedFiles), strings.Join(f.diff.ChangesToString(f.diff.RemovedFiles), ", "))
	}

	if len(f.diff.MovedDirectories) > 0 {
		logger.Debugf("Moved directories %d - [%s]", len(f.diff.MovedDirectories), strings.Join(f.diff.ChangesToString(f.diff.MovedDirectories), ", "))
	}

	if len(f.diff.MovedFiles) > 0 {
		logger.Debugf("Moved files %d - [%s]", len(f.diff.MovedFiles), strings.Join(f.diff.ChangesToString(f.diff.MovedFiles), ", "))
	}
}

func addMissingDirectoryTags(dr model.DirectoryReader, trw model.TagReaderWriter) []error {
	errors := make([]error, 0)
	allDirectories, err := dr.GetAllDirectories()
	if err != nil {
		return append(errors, err)
	}

	for _, dir := range *allDirectories {
		errors = append(errors, addMissingDirectoryTag(dr, trw, &dir))
	}

	return errors
}

func addMissingDirectoryTag(dr model.DirectoryReader, trw model.TagReaderWriter, dir *model.Directory) error {
	if directories.IsExcluded(dir) {
		return nil
	}

	title := directories.DirectoryNameToTag(dir.Path)
	_, err := tags.GetOrCreateChildTag(trw, directories.DIRECTORIES_TAG_ID, title)
	return err
}

func removeStaleItems(dig model.DirectoryItemsGetter, iw model.ItemWriter, files []string) []error {
	errors := make([]error, 0)
	for _, file := range files {
		errors = append(errors, removeItem(dig, iw, file)...)
	}

	return errors
}

func removeItem(dig model.DirectoryItemsGetter, iw model.ItemWriter, file string) []error {
	dirpath := directories.NormalizeDirectoryPath(filepath.Dir(file))
	item, err := dig.GetBelongingItem(dirpath, filepath.Base(file))
	if err != nil {
		return []error{err}
	}

	if item != nil {
		return items.RemoveItemAndItsAssociations(iw, item)
	}

	return []error{}
}

func removeStaleDirs(trw model.TagReaderWriter, dw model.DirectoryWriter, dirs []string) []error {
	errors := make([]error, 0)
	for _, path := range dirs {
		errors = append(errors, removeDir(trw, dw, path)...)
	}

	return errors
}

func removeDir(trw model.TagReaderWriter, dw model.DirectoryWriter, path string) []error {
	errors := make([]error, 0)
	dir := newFsDirectory(directories.NormalizeDirectoryPath(path))
	tag, err := dir.getTag(trw)
	if err != nil {
		errors = append(errors, err)
	} else if tag != nil {
		errors = append(errors, tags.RemoveTagAndItsAssociations(trw, tag)...)
	}

	if err := dw.RemoveDirectory(directories.NormalizeDirectoryPath(path)); err != nil {
		errors = append(errors, err)
	}

	return errors
}

func addMissingDirs(drw model.DirectoryReaderWriter, addedDirectories []directorytree.Change) []error {
	errors := make([]error, 0)
	for _, change := range addedDirectories {
		err := directories.AddDirectoryIfMissing(drw, change.Path1, true)
		if err != nil {
			errors = append(errors, err)
		}
	}

	return errors
}

func addNewFiles(iw model.ItemWriter, digs model.DirectoryItemsGetterSetter, dctg model.DirectoryConcreteTagsGetter,
	flmg model.FileLastModifiedGetter, addedFiles []directorytree.Change) []error {
	errors := make([]error, 0)
	for _, change := range addedFiles {
		dirpath := directories.NormalizeDirectoryPath(filepath.Dir(change.Path1))
		item, err := digs.GetBelongingItem(dirpath, filepath.Base(change.Path1))
		if err != nil {
			errors = append(errors, err)
			continue
		}

		if err := handleFile(iw, digs, dctg, flmg, item, dirpath, change.Path1); err != nil {
			errors = append(errors, err)
		}
	}

	return errors
}

func handleFile(iw model.ItemWriter, digs model.DirectoryItemsGetterSetter, dctg model.DirectoryConcreteTagsGetter,
	flmg model.FileLastModifiedGetter, item *model.Item, dirpath string, path string) error {
	concreteTags, err := dctg.GetConcreteTags(dirpath)
	if err != nil {
		return err
	}

	if item != nil {
		return items.EnsureItemHaveTags(iw, item, concreteTags)
	} else {
		return handleNewFile(digs, flmg, dirpath, path, concreteTags)
	}
}

func handleNewFile(digs model.DirectoryItemsGetterSetter, flmg model.FileLastModifiedGetter,
	dirpath string, path string, concreteTags []*model.Tag) error {
	item, err := items.BuildItemFromPath(dirpath, relativasor.GetAbsoluteFile(path), flmg)
	if err != nil {
		return err
	}

	item.Tags = concreteTags
	return digs.AddBelongingItem(item)
}

func removeDeletedDirs(trw model.TagReaderWriter, dw model.DirectoryWriter, deletedDirs []directorytree.Change) []error {
	errors := make([]error, 0)
	for _, dir := range deletedDirs {
		errors = append(errors, removeDir(trw, dw, dir.Path1)...)
	}

	return errors
}

func removeDeletedFiles(dig model.DirectoryItemsGetter, iw model.ItemWriter, deletedFiles []directorytree.Change) []error {
	errors := make([]error, 0)
	for _, file := range deletedFiles {
		errors = append(errors, removeItem(dig, iw, file.Path1)...)
	}

	return errors
}

func (f *fsSyncer) renameDirs() []error {
	errors := make([]error, 0)

	return errors
}

func renameFiles(trw model.TagReaderWriter, drw model.DirectoryReaderWriter,
	irw model.ItemReaderWriter, movedFiles []directorytree.Change) []error {
	errs := make([]error, 0)
	for _, file := range movedFiles {
		src := file.Path1
		dst := file.Path2

		errs = append(errs, moveFile(trw, drw, irw, src, dst))
	}

	return errs
}

func moveFile(trw model.TagReaderWriter, drw model.DirectoryReaderWriter,
	irw model.ItemReaderWriter, src string, dst string) error {
	if err := validateReadyDirectory(trw, drw, dst); err != nil {
		return err
	}

	srcDirpath := directories.NormalizeDirectoryPath(filepath.Dir(src))
	srcDir := newFsDirectory(srcDirpath)
	item, err := srcDir.getItem(trw, irw, items.TitleFromFileName(src))
	if err != nil {
		return err
	}
	if item == nil {
		return errors.Errorf("original item not found %s", src)
	}

	dstDirpath := directories.NormalizeDirectoryPath(filepath.Dir(dst))
	if err := items.UpdateFileLocation(irw, item, dstDirpath, dst); err != nil {
		return err
	}

	dstDir := newFsDirectory(dstDirpath)
	if err := dstDir.addItem(trw, irw, item); err != nil {
		return err
	}

	return srcDir.removeItem(trw, irw, item)
}

func validateReadyDirectory(trw model.TagReaderWriter, drw model.DirectoryReaderWriter, path string) error {
	dirpath := directories.NormalizeDirectoryPath(filepath.Dir(path))

	if err := directories.AddDirectoryIfMissing(drw, dirpath, false); err != nil {
		return err
	}

	dir, err := directories.GetDirectory(drw, dirpath)
	if err != nil {
		return err
	}

	if directories.IsExcluded(dir) {
		dir.Excluded = pointer.Bool(false)
		if err := drw.CreateOrUpdateDirectory(dir); err != nil {
			return err
		}
	}

	return addMissingDirectoryTag(drw, trw, dir)
}
