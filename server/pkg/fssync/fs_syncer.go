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
	"gorm.io/gorm"
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

	// TODO:
	// > remove moved dirs and files from stale dirs and items
	// > tag title+parent must be unique (representing the full origin location)
	// > remove concrete tags when removing directory
	// > update directory path when renaming dirs
	// > remove concrete tags from items if they removed from their dir (how do we know its a concrete tag?)
	// > redefine concrete tags - what are they? maybe auto tags, allow more flexibility with them

	errs := make([]error, 0)
	errs = append(errs, addMissingDirectoryTags(db, db)...)

	if f.hasFsChanges() {
		f.debugPrint()
		errs = append(errs, removeStaleItems(digs, db, f.stales.Files)...)
		errs = append(errs, removeStaleDirs(db, db, f.stales.Dirs)...)
		errs = append(errs, addMissingDirs(db, f.diff.AddedDirectories)...)
		errs = append(errs, removeDeletedDirs(db, db, f.diff.RemovedDirectories)...)
		errs = append(errs, removeDeletedFiles(digs, db, f.diff.RemovedFiles)...)
		errs = append(errs, renameDirs(db, db, db, f.diff.MovedDirectories)...)
		errs = append(errs, renameFiles(db, db, db, f.diff.MovedFiles)...)
		errs = append(errs, addNewFiles(db, digs, dctg, flmg, f.diff.AddedFiles)...)
	}

	errs = append(errs, syncConcreteTags(db, db, db, dctg)...)
	return errs
}

func (f *fsSyncer) debugPrint() {
	if len(f.stales.Dirs) > 0 {
		logger.Debugf("%d Stale directories - [%s]", len(f.stales.Dirs), strings.Join(f.stales.Dirs, ", "))
	}

	if len(f.stales.Files) > 0 {
		logger.Debugf("%d Stale files - [%s]", len(f.stales.Files), strings.Join(f.stales.Files, ", "))
	}

	if len(f.diff.AddedDirectories) > 0 {
		logger.Debugf("%d Added directories - [%s]", len(f.diff.AddedDirectories), strings.Join(f.diff.ChangesToString(f.diff.AddedDirectories), ", "))
	}

	if len(f.diff.RemovedDirectories) > 0 {
		logger.Debugf("%d Removed directories - [%s]", len(f.diff.RemovedDirectories), strings.Join(f.diff.ChangesToString(f.diff.RemovedDirectories), ", "))
	}

	if len(f.diff.AddedFiles) > 0 {
		logger.Debugf("%d Added files - [%s]", len(f.diff.AddedFiles), strings.Join(f.diff.ChangesToString(f.diff.AddedFiles), ", "))
	}

	if len(f.diff.RemovedFiles) > 0 {
		logger.Debugf("%d Removed files - [%s]", len(f.diff.RemovedFiles), strings.Join(f.diff.ChangesToString(f.diff.RemovedFiles), ", "))
	}

	if len(f.diff.MovedDirectories) > 0 {
		logger.Debugf("%d Moved directories - [%s]", len(f.diff.MovedDirectories), strings.Join(f.diff.ChangesToString(f.diff.MovedDirectories), ", "))
	}

	if len(f.diff.MovedFiles) > 0 {
		logger.Debugf("%d Moved files - [%s]", len(f.diff.MovedFiles), strings.Join(f.diff.ChangesToString(f.diff.MovedFiles), ", "))
	}
}

func addMissingDirectoryTags(dr model.DirectoryReader, trw model.TagReaderWriter) []error {
	errs := make([]error, 0)
	allDirectories, err := dr.GetAllDirectories()
	if err != nil {
		return append(errs, err)
	}

	for _, dir := range *allDirectories {
		if err := addMissingDirectoryTag(dr, trw, &dir); err != nil {
			errs = append(errs, err)
		}
	}

	return errs
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
	errs := make([]error, 0)
	for _, file := range files {
		errs = append(errs, removeItem(dig, iw, file)...)
	}

	return errs
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
	errs := make([]error, 0)
	for _, path := range dirs {
		errs = append(errs, removeDir(trw, dw, path)...)
	}

	return errs
}

func removeDir(trw model.TagReaderWriter, dw model.DirectoryWriter, path string) []error {
	errs := make([]error, 0)
	dir := newFsDirectory(directories.NormalizeDirectoryPath(path))
	tag, err := dir.getTag(trw)
	if err != nil {
		errs = append(errs, err)
	} else if tag != nil {
		errs = append(errs, tags.RemoveTagAndItsAssociations(trw, tag)...)
	}

	if err := dw.RemoveDirectory(directories.NormalizeDirectoryPath(path)); err != nil {
		errs = append(errs, err)
	}

	return errs
}

func addMissingDirs(drw model.DirectoryReaderWriter, addedDirectories []directorytree.Change) []error {
	errs := make([]error, 0)
	for _, change := range addedDirectories {
		err := directories.AddDirectoryIfMissing(drw, change.Path1, true)
		if err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}

func addNewFiles(iw model.ItemWriter, digs model.DirectoryItemsGetterSetter, dctg model.DirectoryConcreteTagsGetter,
	flmg model.FileLastModifiedGetter, addedFiles []directorytree.Change) []error {
	errs := make([]error, 0)
	for _, change := range addedFiles {
		dirpath := directories.NormalizeDirectoryPath(filepath.Dir(change.Path1))
		item, err := digs.GetBelongingItem(dirpath, filepath.Base(change.Path1))
		if err != nil {
			errs = append(errs, err)
			continue
		}

		if err := handleFile(iw, digs, dctg, flmg, item, dirpath, change.Path1); err != nil {
			errs = append(errs, err)
		}
	}

	return errs
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
	errs := make([]error, 0)
	for _, dir := range deletedDirs {
		errs = append(errs, removeDir(trw, dw, dir.Path1)...)
	}

	return errs
}

func removeDeletedFiles(dig model.DirectoryItemsGetter, iw model.ItemWriter, deletedFiles []directorytree.Change) []error {
	errs := make([]error, 0)
	for _, file := range deletedFiles {
		errs = append(errs, removeItem(dig, iw, file.Path1)...)
	}

	return errs
}

func renameDirs(trw model.TagReaderWriter, drw model.DirectoryReaderWriter,
	irw model.ItemReaderWriter, movedDirs []directorytree.Change) []error {
	errs := make([]error, 0)
	for _, dir := range movedDirs {
		src := dir.Path1
		dst := dir.Path2

		if err := moveDir(trw, drw, irw, src, dst); err != nil {
			errs = append(errs, err...)
		}
	}

	return errs
}

func updateItemsLocation(trw model.TagReaderWriter, irw model.ItemReaderWriter, path string) []error {
	dirpath := directories.NormalizeDirectoryPath(path)
	errs := make([]error, 0)
	dir := newFsDirectory(dirpath)
	belongingItems, err := dir.getItems(trw, irw)
	if err != nil {
		return append(errs, err)
	}

	for _, item := range *belongingItems {
		if err := items.UpdateFileLocation(irw, &item, dirpath, filepath.Join(item.Origin, item.Title)); err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}

func moveDir(trw model.TagReaderWriter, drw model.DirectoryReaderWriter,
	irw model.ItemReaderWriter, src string, dst string) []error {
	errs := make([]error, 0)
	dstDirpath := directories.NormalizeDirectoryPath(dst)
	dstdir, err := directories.GetDirectory(drw, dstDirpath)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return append(errs, err)
	}
	if dstdir == nil {
		return removeDir(trw, drw, src)
	}

	srcDirpath := directories.NormalizeDirectoryPath(src)
	srcdir, err := directories.GetDirectory(drw, srcDirpath)
	if err != nil {
		return append(errs, err)
	}
	if srcdir == nil {
		return append(errs, errors.Errorf("source directory not exists %s", src))
	}

	if err := directories.UpdatePath(drw, srcdir, dst); err != nil {
		return append(errs, err)
	}

	if directories.IsExcluded(srcdir) {
		return errs
	}

	return append(errs, updateItemsLocation(trw, irw, dst)...)
}

func renameFiles(trw model.TagReaderWriter, drw model.DirectoryReaderWriter,
	irw model.ItemReaderWriter, movedFiles []directorytree.Change) []error {
	errs := make([]error, 0)
	for _, file := range movedFiles {
		src := file.Path1
		dst := file.Path2

		if err := moveFile(trw, drw, irw, src, dst); err != nil {
			errs = append(errs, err)
		}
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
	dir, err := directories.ValidateReadyDirectory(drw, filepath.Dir(path))
	if err != nil {
		return err
	}

	return addMissingDirectoryTag(drw, trw, dir)
}

func syncConcreteTags(tr model.TagReader, irw model.ItemReaderWriter,
	dr model.DirectoryReader, dctg model.DirectoryConcreteTagsGetter) []error {
	errs := make([]error, 0)
	allDirectories, err := dr.GetAllDirectories()
	if err != nil {
		return append(errs, err)
	}

	for _, dir := range *allDirectories {
		concreteTags, err := dctg.GetConcreteTags(dir.Path)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		errs = append(errs, syncConcreteTagsForDir(tr, irw, concreteTags, &dir)...)
	}

	return errs
}

func syncConcreteTagsForDir(tr model.TagReader, irw model.ItemReaderWriter,
	concreteTags []*model.Tag, dir *model.Directory) []error {
	errs := make([]error, 0)
	fsdir := newFsDirectory(dir.Path)
	belongingItems, err := fsdir.getItems(tr, irw)
	if err != nil {
		return append(errs, err)
	}

	for _, item := range *belongingItems {
		if err := items.EnsureItemHaveTags(irw, &item, concreteTags); err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}
