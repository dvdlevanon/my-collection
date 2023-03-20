package fssync

import (
	"my-collection/server/pkg/bl/directories"
	"my-collection/server/pkg/bl/items"
	"my-collection/server/pkg/bl/tags"
	"my-collection/server/pkg/db"
	"my-collection/server/pkg/directorytree"
	"my-collection/server/pkg/model"
	"path/filepath"

	"github.com/go-errors/errors"
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

func (f *fsSyncer) sync(db *db.Database, digs model.DirectoryItemsGetterSetter,
	dctg model.DirectoryConcreteTagsGetter, flmg model.FileLastModifiedGetter) []error {
	f.debugPrint()
	errors := make([]error, 0)
	errors = append(errors, addMissingDirectoryTags(db, db)...)
	errors = append(errors, f.removeStaleItems(digs, db)...)
	errors = append(errors, f.removeStaleDirs(db, db)...)
	errors = append(errors, addMissingDirs(db, f.diff.AddedDirectories)...)
	errors = append(errors, f.removeDeletedDirs(f.diff.RemovedDirectories)...)
	errors = append(errors, f.removeDeletedFiles()...)
	errors = append(errors, f.renameDirs()...)
	errors = append(errors, f.renameFiles()...)
	errors = append(errors, addNewFiles(db, digs, dctg, flmg, f.diff.AddedFiles)...)
	// syncConcreteTags()
	return errors
}

func (f *fsSyncer) debugPrint() {
	logger.Infof("Stales: \n%s", f.stales.String())
	logger.Infof("Diff: \n%s", f.diff.String())
}

func addMissingDirectoryTags(dr model.DirectoryReader, trw model.TagReaderWriter) []error {
	errors := make([]error, 0)
	allDirectories, err := dr.GetAllDirectories()
	if err != nil {
		return append(errors, err)
	}

	for _, dir := range *allDirectories {
		if directories.IsExcluded(&dir) {
			continue
		}

		title := directories.DirectoryNameToTag(dir.Path)
		_, err := tags.GetOrCreateChildTag(trw, directories.DIRECTORIES_TAG_ID, title)
		if err != nil {
			errors = append(errors, err)
		}
	}

	return errors
}

func (f *fsSyncer) removeStaleItems(dig model.DirectoryItemsGetter, iw model.ItemWriter) []error {
	errors := make([]error, 0)
	for _, file := range f.stales.Files {
		item, err := dig.GetBelongingItem(filepath.Dir(file), filepath.Base(file))
		if err != nil {
			errors = append(errors, err)
			continue
		}

		items.RemoveItemAndItsAssociations(iw, item)
	}

	return errors
}

func (f *fsSyncer) removeStaleDirs(trw model.TagReaderWriter, dw model.DirectoryWriter) []error {
	errors := make([]error, 0)
	for _, path := range f.stales.Dirs {
		dir := newFsDirectory(path)
		tag, err := dir.getTag(trw)
		if err != nil {
			errors = append(errors, err)
		} else {
			errors = append(errors, tags.RemoveTagAndItsAssociations(trw, tag)...)
		}

		if err := dw.RemoveDirectory(path); err != nil {
			errors = append(errors, err)
		}
	}

	return errors
}

func addMissingDirs(drw model.DirectoryReaderWriter, addedDirectories []directorytree.Change) []error {
	errors := make([]error, 0)
	for _, change := range addedDirectories {
		err := directories.AddExcludedDirectoryIfMissing(drw, change.Path1)
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
	item, err := items.BuildItemFromPath(dirpath, path, flmg)
	if err != nil {
		return err
	}

	item.Tags = concreteTags
	return digs.AddBelongingItem(item)
}

func (f *fsSyncer) removeDeletedDirs(deletedDirs []directorytree.Change) []error {
	errors := make([]error, 0)

	// for _, dir := range deletedDirs {

	// }

	return errors
}

func (f *fsSyncer) removeDeletedFiles() []error {
	errors := make([]error, 0)

	return errors
}

func (f *fsSyncer) renameDirs() []error {
	errors := make([]error, 0)

	return errors
}

func (f *fsSyncer) renameFiles() []error {
	errors := make([]error, 0)

	return errors
}
