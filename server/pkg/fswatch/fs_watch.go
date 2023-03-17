package fswatch

import (
	"my-collection/server/pkg/bl/directories"
	"my-collection/server/pkg/bl/tags"
	"my-collection/server/pkg/db"
	"my-collection/server/pkg/model"
	"my-collection/server/pkg/processor"
	"my-collection/server/pkg/relativasor"
	"my-collection/server/pkg/storage"
	"my-collection/server/pkg/utils"
	"time"

	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("fswatch")

type directoryProcessor interface {
	process() error
}

type Fswatch interface {
	Init() error
	DirectoryChanged(path string)
}

type FswatchMock struct{}

func (d *FswatchMock) Init() error                  { return nil }
func (d *FswatchMock) DirectoryChanged(path string) {}

type fswatchImpl struct {
	db                  *db.Database
	storage             *storage.Storage
	processor           processor.Processor
	changeChannel       chan string
	trustFileExtenssion bool
}

func New(db *db.Database, storage *storage.Storage, processor processor.Processor) Fswatch {
	logger.Infof("FS watch initialized")

	return &fswatchImpl{
		db:                  db,
		storage:             storage,
		processor:           processor,
		changeChannel:       make(chan string),
		trustFileExtenssion: true,
	}
}

func (d *fswatchImpl) Init() error {
	if err := d.db.CreateOrUpdateTag(&directories.DirectoriesTag); err != nil {
		return err
	}

	go d.watchFilesystemChanges()
	return nil
}

func (d *fswatchImpl) DirectoryChanged(path string) {
	d.changeChannel <- path
}

func handleError(err error) {
	logger.Errorf("Error processing %s", err)
}

func (d *fswatchImpl) watchFilesystemChanges() {
	for {
		select {
		case <-d.changeChannel:
			// d.syncDirectory(path)
			d.sync()
		case <-time.After(60 * time.Second):
			// d.periodicScan()
			d.sync()
		}
	}
}

func (d *fswatchImpl) GetBelongingItems(path string) (*[]model.Item, error) {
	return newFsDirectory(path).getItems(d.db, d.db)
}

func (d *fswatchImpl) GetBelongingItem(path string, filename string) (*model.Item, error) {
	return newFsDirectory(path).getItem(d.db, d.db, filename)
}

func (d *fswatchImpl) AddBelongingItem(item *model.Item) error {
	return newFsDirectory(item.Origin).addItem(d.db, d.db, item)
}

func (d *fswatchImpl) GetConcreteTags(path string) ([]*model.Tag, error) {
	directory, err := d.db.GetDirectory(path)
	if err != nil {
		return nil, err
	}

	return tags.GetOrCreateTags(d.db, directories.BuildDirectoryTags(directory))
}

func (g *fswatchImpl) filesFilter(path string) bool {
	return utils.IsVideo(g.trustFileExtenssion, path)
}

func (d *fswatchImpl) sync() {
	fsSync, err := newFsSync(relativasor.GetRootDirectory(), d.filesFilter, d.db, d)
	if err != nil {
		handleError(err)
		return
	}

	fsSync.debugPrint()
	errors := make([]error, 0)
	errors = append(errors, d.addMissingDirectoryTags()...)
	errors = append(errors, fsSync.removeStaleItems(d, d.db)...)
	errors = append(errors, fsSync.removeStaleDirs(d.db, d.db)...)
	errors = append(errors, fsSync.addMissingDirs(d.db)...)
	errors = append(errors, fsSync.removeDeletedDirs()...)
	errors = append(errors, fsSync.removeDeletedFiles()...)
	errors = append(errors, fsSync.renameDirs()...)
	errors = append(errors, fsSync.renameFiles()...)
	fsSync.addNewFiles(d.db, d, d)

	for _, err := range errors {
		handleError(err)
	}

	// syncConcreteTags()
}

func (d *fswatchImpl) addMissingDirectoryTags() []error {
	errors := make([]error, 0)
	allDirectories, err := d.db.GetAllDirectories()
	if err != nil {
		return append(errors, err)
	}

	for _, dir := range *allDirectories {
		title := directories.DirectoryNameToTag(dir.Path)
		_, err := tags.GetOrCreateChildTag(d.db, directories.DIRECTORIES_TAG_ID, title)
		if err != nil {
			errors = append(errors, err)
		}
	}

	return errors
}

func (d *fswatchImpl) periodicScan() {
	// TODO: Get all directories path instead of full directories
	allDirectories, err := directories.GetAllDirectoriesWithCache(d.db)
	if err != nil {
		logger.Errorf("Error getting all directories %t", err)
		return
	}

	for _, dir := range *allDirectories {
		millisSinceScanned := time.Now().UnixMilli() - dir.LastSynced

		if millisSinceScanned < 1000*60*5 {
			return
		}

		d.syncDirectory(dir.Path)
	}
}

func (d *fswatchImpl) syncDirectory(path string) {
	directory, err := d.db.GetDirectory(path)
	if err != nil {
		handleError(err)
		return
	}

	if err := directories.StartDirectoryProcessing(d.db, directory); err != nil {
		logger.Errorf("Error updating directory %s %t", directory.Path, err)
	}

	handleError(d.processDirectory(directory))

	if err := directories.FinishDirectoryProcessing(d.db, directory); err != nil {
		logger.Errorf("Error updating directory %s %t", directory.Path, err)
	}
}

func (d *fswatchImpl) processDirectory(directory *model.Directory) error {
	var processor directoryProcessor
	var err error

	if directories.IsExcluded(directory) {
		processor, err = newDirectoryExcluder(directory, d.db)
	} else {
		processor, err = newDirectoryIncluder(d.trustFileExtenssion, directory, d.db, d.db, d.db, d.processor)
	}

	if err != nil {
		return err
	}

	return processor.process()
}

// func (d *fswatchImpl) handleExcludedDirectory(directory *model.Directory) {
// 	tag, err := tags.GetChildTag(d.db, DIRECTORIES_TAG_ID, directories.DirectoryNameToTag(directory.Path))
// 	if err != nil {
// 		logger.Errorf("Unable to find directory of %s - %s", directory.Path, err)
// 		return
// 	}

// 	for _, item := range tag.Items {
// 		if !items.HasSingleTag(item, tag) {
// 			continue
// 		}

// 		items.RemoveItemAndItsAssociations(d.db, item)
// 	}

// 	if err := d.db.RemoveTag(tag.Id); err != nil {
// 		logger.Errorf("Unable to remove tag %d - %s", tag.Id, err)
// 	}
// }

// func (d *fswatchImpl) removeExcludedSubDirectories(directoryPath string) {
// 	allDirectories, err := d.db.GetAllDirectories()
// 	if err != nil {
// 		logger.Errorf("Error getting all directories %t", err)
// 		return
// 	}

// 	for _, dir := range *allDirectories {
// 		if dir.Excluded == nil || !*dir.Excluded {
// 			continue
// 		}

// 		if strings.HasPrefix(dir.Path, directoryPath) {
// 			if err := d.db.RemoveDirectory(dir.Path); err != nil {
// 				logger.Errorf("Error removing directory %s", dir.Path)
// 			}
// 		}
// 	}
// }
