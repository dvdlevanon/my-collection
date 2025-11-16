package fssync

import (
	"context"
	"my-collection/server/pkg/bl/directories"
	"my-collection/server/pkg/bl/items"
	"my-collection/server/pkg/bl/tags"
	"my-collection/server/pkg/directorytree"
	"my-collection/server/pkg/model"
	"my-collection/server/pkg/relativasor"
	"path/filepath"
	"strings"

	"github.com/go-errors/errors"
	"gorm.io/gorm"
)

func newFsSyncer(ctx context.Context, path string, db model.Database, dig model.DirectoryItemsGetter, filter directorytree.FilesFilter) (*fsSyncer, error) {
	exists, err := directories.DirectoryExists(ctx, db, path)
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

	rootdb, err := directorytree.BuildFromDb(ctx, db, dig)
	if err != nil {
		return nil, err
	}

	diff := directorytree.Compare(rootfs, rootdb)
	stales := directorytree.FindStales(rootdb)

	return &fsSyncer{
		diff:   diff,
		stales: stales,
		ctx:    ctx,
	}, nil
}

type fsSyncer struct {
	stales *directorytree.Stale
	diff   *directorytree.Diff
	ctx    context.Context
}

func (f *fsSyncer) hasFsChanges() bool {
	return f.stales.HasChanges() || f.diff.HasChanges()
}

func (f *fsSyncer) sync(ctx context.Context, db model.Database, digs model.DirectoryItemsGetterSetter,
	datg model.DirectoryAutoTagsGetter, fmg model.FileMetadataGetter) (bool, []error) {

	// TODO:
	// > remove moved dirs and files from stale dirs and items
	// > tag title+parent must be unique (representing the full origin location)
	// > remove auto tags when removing directory
	// > update directory path when renaming dirs
	// > remove auto tags from items if they removed from their dir (how do we know its an auto tag?)
	// > redefine auto tags - what are they? maybe auto tags, allow more flexibility with them

	errs := make([]error, 0)
	errs = append(errs, addMissingDirectoryTags(ctx, db, db)...)

	if f.hasFsChanges() {
		f.debugPrint()
		errs = append(errs, removeStaleItems(ctx, digs, db, f.stales.Files)...)
		errs = append(errs, removeStaleDirs(ctx, db, db, f.stales.Dirs)...)
		errs = append(errs, addMissingDirs(ctx, db, f.diff.AddedDirectories)...)
		errs = append(errs, renameDirs(ctx, db, db, db, f.diff.MovedDirectories)...)
		errs = append(errs, renameFiles(ctx, db, db, db, f.diff.MovedFiles)...)
		errs = append(errs, removeDeletedDirs(ctx, db, db, f.diff.RemovedDirectories)...)
		errs = append(errs, removeDeletedFiles(ctx, digs, db, f.diff.RemovedFiles)...)
		errs = append(errs, addNewFiles(ctx, db, digs, datg, fmg, f.diff.AddedFiles)...)
	}

	anyItemChanged, tagsErrs := syncAutoTags(ctx, db, db, db, datg)
	errs = append(errs, tagsErrs...)
	return f.hasFsChanges() || anyItemChanged, errs
}

func (f *fsSyncer) debugPrint() {
	logger.Debugf("Sync starting")

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

func addMissingDirectoryTags(ctx context.Context, dr model.DirectoryReader, trw model.TagReaderWriter) []error {
	errs := make([]error, 0)
	allDirectories, err := dr.GetAllDirectories(ctx)
	if err != nil {
		return append(errs, err)
	}

	for i, dir := range *allDirectories {
		if err := addMissingDirectoryTag(ctx, trw, &dir); err != nil {
			errs = append(errs, err)
		}
		if i+1%100 == 0 {
			logger.Debugf("Added %d/%d directory tags", i, len(*allDirectories))
		}
	}

	return errs
}

func addMissingDirectoryTag(ctx context.Context, trw model.TagReaderWriter, dir *model.Directory) error {
	if directories.IsExcluded(dir) {
		return nil
	}

	_, err := tags.GetOrCreateChildTag(ctx, trw, directories.GetDirectoriesTagId(), dir.Path)
	return err
}

func removeStaleItems(ctx context.Context, dig model.DirectoryItemsGetter, iw model.ItemWriter, files []string) []error {
	errs := make([]error, 0)
	for _, file := range files {
		errs = append(errs, removeItem(ctx, dig, iw, file)...)
	}

	return errs
}

func removeItem(ctx context.Context, dig model.DirectoryItemsGetter, iw model.ItemWriter, file string) []error {
	dirpath := directories.NormalizeDirectoryPath(filepath.Dir(file))
	item, err := dig.GetBelongingItem(ctx, dirpath, filepath.Base(file))
	if err != nil {
		return []error{err}
	}

	if item != nil {
		return items.RemoveItemAndItsAssociations(ctx, iw, item.Id)
	}

	return []error{}
}

func removeStaleDirs(ctx context.Context, trw model.TagReaderWriter, dw model.DirectoryWriter, dirs []string) []error {
	errs := make([]error, 0)
	for _, path := range dirs {
		errs = append(errs, removeDir(ctx, trw, dw, path)...)
	}

	return errs
}

func removeDir(ctx context.Context, trw model.TagReaderWriter, dw model.DirectoryWriter, path string) []error {
	errs := make([]error, 0)
	dir := newFsDirectory(directories.NormalizeDirectoryPath(path))
	tag, err := dir.getTag(ctx, wrapDb(trw, nil))
	if err != nil {
		errs = append(errs, err)
	} else if tag != nil {
		errs = append(errs, tags.RemoveTagAndItsAssociations(ctx, trw, tag)...)
	}

	if err := dw.RemoveDirectory(ctx, directories.NormalizeDirectoryPath(path)); err != nil {
		errs = append(errs, err)
	}

	return errs
}

func addMissingDirs(ctx context.Context, drw model.DirectoryReaderWriter, addedDirectories []directorytree.Change) []error {
	errs := make([]error, 0)
	for i, change := range addedDirectories {
		shouldInclude := directories.ShouldInclude(ctx, drw, change.Path1)
		err := directories.AddDirectoryIfMissing(ctx, drw, change.Path1, !shouldInclude)
		if err != nil {
			errs = append(errs, err)
		}
		if i+1%100 == 0 {
			logger.Debugf("Added %d/%d directories", i, len(addedDirectories))
		}
	}
	return errs
}

func addNewFiles(ctx context.Context, iw model.ItemWriter, digs model.DirectoryItemsGetterSetter, datg model.DirectoryAutoTagsGetter,
	fmg model.FileMetadataGetter, addedFiles []directorytree.Change) []error {
	errs := make([]error, 0)
	for i, change := range addedFiles {
		dirpath := directories.NormalizeDirectoryPath(filepath.Dir(change.Path1))
		item, err := digs.GetBelongingItem(ctx, dirpath, filepath.Base(change.Path1))
		if err != nil {
			errs = append(errs, err)
			continue
		}

		if err := handleFile(ctx, iw, digs, datg, fmg, item, dirpath, change.Path1); err != nil {
			errs = append(errs, err)
		}
		if i+1%100 == 0 {
			logger.Debugf("Added %d/%d files", i, len(addedFiles))
		}
	}

	return errs
}

func handleFile(ctx context.Context, iw model.ItemWriter, digs model.DirectoryItemsGetterSetter, datg model.DirectoryAutoTagsGetter,
	fmg model.FileMetadataGetter, item *model.Item, dirpath string, path string) error {
	autoTags, err := datg.GetAutoTags(ctx, dirpath)
	if err != nil {
		return err
	}

	if item != nil {
		_, err := items.EnsureItemHaveTags(ctx, iw, item, autoTags)
		return err
	} else {
		return handleNewFile(ctx, digs, fmg, dirpath, path, autoTags)
	}
}

func handleNewFile(ctx context.Context, digs model.DirectoryItemsGetterSetter, fmg model.FileMetadataGetter,
	dirpath string, path string, autoTags []*model.Tag) error {
	item, err := items.BuildItemFromPath(dirpath, relativasor.GetAbsoluteFile(path), fmg)
	if err != nil {
		return err
	}

	item.Tags = autoTags
	return digs.AddBelongingItem(ctx, item)
}

func removeDeletedDirs(ctx context.Context, trw model.TagReaderWriter, dw model.DirectoryWriter, deletedDirs []directorytree.Change) []error {
	errs := make([]error, 0)
	for _, dir := range deletedDirs {
		errs = append(errs, removeDir(ctx, trw, dw, dir.Path1)...)
	}

	return errs
}

func removeDeletedFiles(ctx context.Context, dig model.DirectoryItemsGetter, iw model.ItemWriter, deletedFiles []directorytree.Change) []error {
	errs := make([]error, 0)
	for _, file := range deletedFiles {
		errs = append(errs, removeItem(ctx, dig, iw, file.Path1)...)
	}

	return errs
}

func renameDirs(ctx context.Context, trw model.TagReaderWriter, drw model.DirectoryReaderWriter,
	irw model.ItemReaderWriter, movedDirs []directorytree.Change) []error {
	errs := make([]error, 0)
	for _, dir := range movedDirs {
		src := dir.Path1
		dst := dir.Path2

		if err := moveDir(ctx, trw, drw, irw, src, dst); err != nil {
			errs = append(errs, err...)
		}
	}

	return errs
}

func updateItemsLocation(ctx context.Context, trw model.TagReaderWriter, irw model.ItemReaderWriter, path string) []error {
	dirpath := directories.NormalizeDirectoryPath(path)
	errs := make([]error, 0)
	dir := newFsDirectory(dirpath)
	belongingItems, err := dir.getItems(ctx, wrapDb(trw, irw))
	if err != nil {
		return append(errs, err)
	}

	for _, item := range *belongingItems {
		if err := items.UpdateFileLocation(ctx, irw, &item, dirpath, filepath.Join(item.Origin, item.Title), ""); err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}

func moveDir(ctx context.Context, trw model.TagReaderWriter, drw model.DirectoryReaderWriter,
	irw model.ItemReaderWriter, src string, dst string) []error {
	errs := make([]error, 0)
	dstDirpath := directories.NormalizeDirectoryPath(dst)
	dstdir, err := directories.GetDirectory(ctx, drw, dstDirpath)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return append(errs, err)
	}
	if dstdir == nil {
		return removeDir(ctx, trw, drw, src)
	}

	srcDirpath := directories.NormalizeDirectoryPath(src)
	srcdir, err := directories.GetDirectory(ctx, drw, srcDirpath)
	if err != nil {
		return append(errs, err)
	}
	if srcdir == nil {
		return append(errs, errors.Errorf("source directory not exists %s", src))
	}

	if err := directories.UpdatePath(ctx, drw, srcdir, dst); err != nil {
		return append(errs, err)
	}

	if directories.IsExcluded(srcdir) {
		return errs
	}

	return append(errs, updateItemsLocation(ctx, trw, irw, dst)...)
}

func renameFiles(ctx context.Context, trw model.TagReaderWriter, drw model.DirectoryReaderWriter,
	irw model.ItemReaderWriter, movedFiles []directorytree.Change) []error {
	errs := make([]error, 0)
	for _, file := range movedFiles {
		src := file.Path1
		dst := file.Path2

		if err := moveFile(ctx, trw, drw, irw, src, dst); err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}

func moveFile(ctx context.Context, trw model.TagReaderWriter, drw model.DirectoryReaderWriter,
	irw model.ItemReaderWriter, src string, dst string) error {
	if err := validateReadyDirectory(ctx, trw, drw, dst); err != nil {
		return err
	}

	srcDirpath := directories.NormalizeDirectoryPath(filepath.Dir(src))
	srcDir := newFsDirectory(srcDirpath)
	item, err := srcDir.getItem(ctx, wrapDb(trw, irw), items.TitleFromFileName(src))
	if err != nil {
		return err
	}
	if item == nil {
		return errors.Errorf("original item not found %s", src)
	}

	dstDirpath := directories.NormalizeDirectoryPath(filepath.Dir(dst))
	if err := items.UpdateFileLocation(ctx, irw, item, dstDirpath, dst, ""); err != nil {
		return err
	}

	dstDir := newFsDirectory(dstDirpath)
	if err := dstDir.addItem(ctx, trw, irw, item); err != nil {
		return err
	}

	return srcDir.removeItem(ctx, trw, irw, item)
}

func validateReadyDirectory(ctx context.Context, trw model.TagReaderWriter, drw model.DirectoryReaderWriter, path string) error {
	dir, err := directories.ValidateReadyDirectory(ctx, drw, filepath.Dir(path))
	if err != nil {
		return err
	}

	return addMissingDirectoryTag(ctx, trw, dir)
}

func syncAutoTags(ctx context.Context, tr model.TagReader, irw model.ItemReaderWriter,
	dr model.DirectoryReader, datg model.DirectoryAutoTagsGetter) (bool, []error) {
	errs := make([]error, 0)
	allDirectories, err := dr.GetAllDirectories(ctx)
	if err != nil {
		return false, append(errs, err)
	}

	anyItemChanged := false
	for _, dir := range *allDirectories {
		autoTags, err := datg.GetAutoTags(ctx, dir.Path)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		itemChanged, syncErrs := syncAutoTagsForDir(ctx, tr, irw, autoTags, &dir)
		errs = append(errs, syncErrs...)
		anyItemChanged = anyItemChanged || itemChanged
	}

	return anyItemChanged, errs
}

func syncAutoTagsForDir(ctx context.Context, tr model.TagReader, irw model.ItemReaderWriter,
	autoTags []*model.Tag, dir *model.Directory) (bool, []error) {
	errs := make([]error, 0)
	fsdir := newFsDirectory(dir.Path)
	belongingItems, err := fsdir.getItems(ctx, wrapDb(tr, irw))
	if err != nil {
		return false, append(errs, err)
	}

	anyItemChanged := false
	for _, item := range *belongingItems {
		itemChanged, err := items.EnsureItemHaveTags(ctx, irw, &item, autoTags)
		if err != nil {
			errs = append(errs, err)
		}
		anyItemChanged = anyItemChanged || itemChanged
	}

	return anyItemChanged, errs
}
