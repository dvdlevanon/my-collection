package fswatch

import (
	"my-collection/server/pkg/bl/directories"
	"my-collection/server/pkg/bl/items"
	"my-collection/server/pkg/bl/tags"
	"my-collection/server/pkg/directorytree"
	"my-collection/server/pkg/model"
	"path/filepath"
)

func newFsSync(path string, fileFilter directorytree.FilesFilter, dr model.DirectoryReader, dig model.DirectoryItemsGetter) (*fsSync, error) {
	rootfs, err := directorytree.BuildFromPath(path, fileFilter)
	if err != nil {
		return nil, err
	}

	rootdb, err := directorytree.BuildFromDb(dr, dig)
	if err != nil {
		return nil, err
	}

	diff := directorytree.Compare(rootfs, rootdb)
	stales := directorytree.FindStales(rootdb)

	return &fsSync{
		diff:   diff,
		stales: stales,
	}, nil
}

type fsSync struct {
	stales *directorytree.Stale
	diff   *directorytree.Diff
}

func (f *fsSync) debugPrint() {
	logger.Infof("Stales: \n%s", f.stales.String())
	logger.Infof("Diff: \n%s", f.diff.String())
}

func (f *fsSync) removeStaleItems(dig model.DirectoryItemsGetter, iw model.ItemWriter) []error {
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

func (f *fsSync) removeStaleDirs(trw model.TagReaderWriter, dw model.DirectoryWriter) []error {
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

func (f *fsSync) addMissingDirs(drw model.DirectoryReaderWriter) []error {
	errors := make([]error, 0)
	for _, change := range f.diff.AddedDirectories {
		err := directories.AddExcludedDirectoryIfMissing(drw, change.Path1)
		if err != nil {
			errors = append(errors, err)
		}
	}

	return errors
}

func (f *fsSync) addNewFiles(iw model.ItemWriter, digs model.DirectoryItemsGetterSetter, dctg model.DirectoryConcreteTagsGetter) []error {
	errors := make([]error, 0)
	for _, change := range f.diff.AddedFiles {
		item, err := digs.GetBelongingItem(filepath.Dir(change.Path1), filepath.Base(change.Path1))
		if err != nil {
			errors = append(errors, err)
			continue
		}

		if err := f.handleFile(iw, digs, dctg, item, change.Path1); err != nil {
			errors = append(errors, err)
		}
	}

	return errors
}

func (f *fsSync) handleFile(iw model.ItemWriter, digs model.DirectoryItemsGetterSetter, dctg model.DirectoryConcreteTagsGetter, item *model.Item, path string) error {
	concreteTags, err := dctg.GetConcreteTags(filepath.Dir(path))
	if err != nil {
		return err
	}

	if item != nil {
		return items.EnsureItemHaveTags(iw, item, concreteTags)
	} else {
		return f.handleNewFile(digs, path, concreteTags)
	}
}

func (f *fsSync) handleNewFile(digs model.DirectoryItemsGetterSetter, path string, concreteTags []*model.Tag) error {
	item, err := items.BuildItemFromPath(path)
	if err != nil {
		return err
	}

	item.Tags = concreteTags
	return digs.AddBelongingItem(item)
}

func (f *fsSync) removeDeletedDirs() []error {
	errors := make([]error, 0)

	return errors
}

func (f *fsSync) removeDeletedFiles() []error {
	errors := make([]error, 0)

	return errors
}

func (f *fsSync) renameDirs() []error {
	errors := make([]error, 0)

	return errors
}

func (f *fsSync) renameFiles() []error {
	errors := make([]error, 0)

	return errors
}
