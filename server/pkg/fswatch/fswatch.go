package fswatch

import (
	"errors"
	"math"
	"my-collection/server/pkg/bl/directories"
	"my-collection/server/pkg/bl/tags"
	"my-collection/server/pkg/db"
	"my-collection/server/pkg/model"
	processor "my-collection/server/pkg/processor"
	"my-collection/server/pkg/relativasor"
	"my-collection/server/pkg/storage"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/h2non/filetype"
	"github.com/op/go-logging"
	"gorm.io/gorm"
	"k8s.io/utils/pointer"
)

const DIRECTORIES_TAG_ID = uint64(1) // tags-util.js

var logger = logging.MustGetLogger("fswatch")
var directoriesTag = model.Tag{
	Id:    DIRECTORIES_TAG_ID,
	Title: "Directories",
}

type Fswatch interface {
	Init() error
	DirectoryChanged(directory *model.Directory)
	DirectoryExcluded(path string)
}

type FswatchMock struct{}

func (d *FswatchMock) Init() error                                 { return nil }
func (d *FswatchMock) DirectoryChanged(directory *model.Directory) {}
func (d *FswatchMock) DirectoryExcluded(path string)               {}

type fswatchImpl struct {
	db                  *db.Database
	storage             *storage.Storage
	processor           processor.Processor
	changeChannel       chan model.Directory
	excludedChannel     chan string
	trustFileExtenssion bool
}

func New(db *db.Database, storage *storage.Storage, processor processor.Processor) Fswatch {
	logger.Infof("FS watch initialized")

	return &fswatchImpl{
		db:                  db,
		storage:             storage,
		processor:           processor,
		changeChannel:       make(chan model.Directory),
		excludedChannel:     make(chan string),
		trustFileExtenssion: true,
	}
}

func (d *fswatchImpl) Init() error {
	if err := d.db.CreateOrUpdateTag(&directoriesTag); err != nil {
		return err
	}

	go d.watchFilesystemChanges()
	return nil
}

func (d *fswatchImpl) DirectoryChanged(directory *model.Directory) {
	d.changeChannel <- *directory
}

func (d *fswatchImpl) DirectoryExcluded(directoryPath string) {
	d.excludedChannel <- directoryPath
}

func (d *fswatchImpl) watchFilesystemChanges() {
	for {
		select {
		case directory := <-d.changeChannel:
			logger.Infof("Directory changed %s", directory.Path)
			d.directoryChanged(&directory)
		case directoryPath := <-d.excludedChannel:
			logger.Infof("Directory excluded %s", directoryPath)
			d.directoryExcluded(directoryPath)
		case <-time.After(60 * time.Second):
			d.periodicScan()
		}
	}
}

func (d *fswatchImpl) periodicScan() {
	allDirectories, err := d.db.GetAllDirectories()
	if err != nil {
		logger.Errorf("Error getting all directories %t", err)
		return
	}

	for _, dir := range *allDirectories {
		millisSinceScanned := time.Now().UnixMilli() - dir.LastSynced

		if millisSinceScanned < 1000*60*5 {
			return
		}

		d.directoryChanged(&dir)
	}
}

func (d *fswatchImpl) directoryChanged(directory *model.Directory) {
	allDirectories, err := d.db.GetAllDirectories()
	if err != nil {
		logger.Errorf("Error getting all directories %t", err)
		return
	}

	if (*directory).Excluded != nil && *((*directory).Excluded) {
		d.handleExcludedDirectory(directory, *allDirectories)
		return
	}

	d.handleIncludedDirectory(directory, *allDirectories)
	directory.LastSynced = time.Now().UnixMilli()
	directory.ProcessingStart = pointer.Int64(0)
	if err := d.db.CreateOrUpdateDirectory(directory); err != nil {
		logger.Errorf("Error updating directory %s %t", directory.Path, err)
	}
}

func (d *fswatchImpl) handleExcludedDirectory(directory *model.Directory, allDirectories []model.Directory) {
	tag, err := d.db.GetTag(model.Tag{
		ParentID: pointer.Uint64(DIRECTORIES_TAG_ID),
		Title:    directories.DirectoryNameToTag(directory.Path),
	})

	if err != nil {
		logger.Errorf("Unable to find directory of %s - %s", directory.Path, err)
		return
	}

	for _, item := range tag.Items {
		if len(item.Tags) > 1 {
			continue
		}

		if err := d.db.RemoveItem(item.Id); err != nil {
			logger.Errorf("Unable to remove item %d - %s", item.Id, err)
		}
	}

	if err := d.db.RemoveTag(tag.Id); err != nil {
		logger.Errorf("Unable to remove tag %d - %s", tag.Id, err)
	}
}

func (d *fswatchImpl) handleIncludedDirectory(directory *model.Directory, allDirectories []model.Directory) {
	tag, err := d.handleDirectoryTag(directory)
	if err != nil {
		return
	}

	directory.FilesCount = pointer.Int(d.syncDirectory(directory, tag, allDirectories))
}

func (d *fswatchImpl) syncDirectory(directory *model.Directory, tag *model.Tag, allDirectories []model.Directory) int {
	path := relativasor.GetAbsoluteFile(directory.Path)
	files, err := os.ReadDir(path)
	if err != nil {
		logger.Errorf("Error getting files of %s %t", path, err)
	}

	filesCount := 0
	for _, file := range files {
		path := filepath.Join(path, file.Name())
		if file.IsDir() {
			d.addDirectoryIfMissing(path, allDirectories)
		} else {
			added, _ := d.addFileIfMissing(directory, tag, path)
			filesCount++

			if added {
				tag, err = d.db.GetTag(tag.Id)
				if err != nil {
					logger.Errorf("Error refetching tag from DB %t", err)
				}
			}
		}
	}

	items, err := tags.GetItems(d.db, tag)
	if err != nil {
		logger.Errorf("Error getting files of tag %t", err)
	}

	for _, item := range *items {
		if d.fileExists(directory, item) {
			continue
		}

		if err := d.db.RemoveItem(item.Id); err != nil {
			logger.Errorf("Error removing item %d - %t", item.Id, err)
		}
	}

	return filesCount
}

func (d *fswatchImpl) fileExists(directory *model.Directory, item model.Item) bool {
	if item.Origin != directory.Path {
		return true
	}

	path := relativasor.GetAbsoluteFile(filepath.Join(item.Origin, item.Title))
	_, err := os.Stat(path)
	return err == nil
}

func (d *fswatchImpl) isVideo(path string) bool {
	if d.trustFileExtenssion {
		return strings.HasSuffix(path, ".avi") ||
			strings.HasSuffix(path, ".mkv") ||
			strings.HasSuffix(path, ".mpg") ||
			strings.HasSuffix(path, ".mpeg") ||
			strings.HasSuffix(path, ".wmv") ||
			strings.HasSuffix(path, ".mp4")
	}

	file, err := os.Open(path)
	if err != nil {
		logger.Errorf("Error opening file for reading %s - %t", file, err)
		return false
	}

	stat, err := file.Stat()
	if err != nil {
		logger.Errorf("Error getting stats of file %s - %t", path, err)
		return false
	}

	header := make([]byte, int(math.Max(float64(stat.Size())-1, 1024)))
	_, err = file.Read(header)
	if err != nil {
		logger.Errorf("Error reading from file %s - %t", path, err)
		return false
	}

	return filetype.IsVideo(header)
}

func (d *fswatchImpl) addDirectoryIfMissing(path string, allDirectories []model.Directory) {
	if d.directoryExists(path, allDirectories) {
		return
	}

	if err := d.addExcludedDirectory(path); err != nil {
		logger.Errorf("Error saving directory %s %t", path, err)
	}
}

func (d *fswatchImpl) directoryExists(path string, allDirectories []model.Directory) bool {
	relativePath := relativasor.GetRelativePath(path)
	for _, dir := range allDirectories {
		if dir.Path == relativePath {
			return true
		}
	}

	return false
}

func (d *fswatchImpl) addExcludedDirectory(path string) error {
	newDirectory := &model.Directory{
		Path:     path,
		Excluded: pointer.Bool(true),
	}

	return d.db.CreateOrUpdateDirectory(newDirectory)
}

func (d *fswatchImpl) getConcreteTagOfDirectory(directory *model.Directory, directoryTag *model.Tag) (*model.Tag, error) {
	tag := model.Tag{
		ParentID: &directoryTag.Id,
		Title:    directories.DirectoryNameToTag(directory.Path),
	}

	return d.getOrCreateTag(&tag)
}

func (d *fswatchImpl) getConcreteTagsOfDirectory(directory *model.Directory) ([]*model.Tag, error) {
	result := make([]*model.Tag, 0)

	for _, directoryTag := range directory.Tags {
		concreteTag, err := d.getConcreteTagOfDirectory(directory, directoryTag)

		if err != nil {
			return nil, err
		}

		result = append(result, concreteTag)
	}

	return result, nil
}

func (d *fswatchImpl) ensureConcreteTagsOnItem(item *model.Item, concreteTags []*model.Tag) error {
	tagsToAdd := make([]*model.Tag, 0)
	for _, concreteTag := range concreteTags {
		exists := false
		for _, tag := range item.Tags {
			if tag.Id == concreteTag.Id {
				exists = true
				break
			}
		}

		if !exists {
			tagsToAdd = append(tagsToAdd, concreteTag)
		}
	}

	if len(tagsToAdd) == 0 {
		return nil
	}

	item.Tags = append(item.Tags, tagsToAdd...)
	if err := d.db.CreateOrUpdateItem(item); err != nil {
		logger.Errorf("Error updating item %v - %t", item, err)
		return err
	}

	return nil
}

func (d *fswatchImpl) addFileIfMissing(directory *model.Directory, tag *model.Tag, path string) (bool, error) {
	existingItem, lastModified, err := d.itemExists(path, tag)

	if err != nil {
		return false, err
	}

	if existingItem == nil && !d.isVideo(path) {
		return false, nil
	}

	concreteTags, err := d.getConcreteTagsOfDirectory(directory)
	if err != nil {
		return false, err
	}

	if existingItem != nil {
		return false, d.ensureConcreteTagsOnItem(existingItem, concreteTags)
	}

	title := filepath.Base(path)
	item := model.Item{
		Title:        title,
		Origin:       directory.Path,
		Url:          path,
		LastModified: lastModified,
		Tags:         append(concreteTags, tag),
	}

	logger.Debugf("Adding a new file %s to %v", path, item)

	if err := d.db.CreateOrUpdateItem(&item); err != nil {
		logger.Errorf("Error creating item %v - %t", item, err)
		return false, err
	}

	if d.processor.IsAutomaticProcessing() {
		d.processor.EnqueueItemVideoMetadata(item.Id)
		d.processor.EnqueueItemCovers(item.Id)
		d.processor.EnqueueItemPreview(item.Id)
	}

	return true, nil
}

func (d *fswatchImpl) itemExists(path string, tag *model.Tag) (*model.Item, int64, error) {
	title := filepath.Base(path)
	file, err := os.Stat(path)
	if err != nil {
		logger.Errorf("Error getting file stat %s - %t", file, err)
		return nil, 0, err
	}

	items, err := tags.GetItems(d.db, tag)
	if err != nil {
		logger.Errorf("Error getting files of tag %t", err)
	}

	for _, item := range *items {
		if item.Title == title && item.LastModified == file.ModTime().UnixMilli() {
			return &item, file.ModTime().UnixMilli(), nil
		}
	}

	return nil, file.ModTime().UnixMilli(), nil
}

func (d *fswatchImpl) getOrCreateTag(tag *model.Tag) (*model.Tag, error) {
	existing, err := d.db.GetTag(tag)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Errorf("Error getting tag %t", err)
		return nil, err
	}

	if existing != nil {
		return existing, nil
	}

	if err := d.db.CreateOrUpdateTag(tag); err != nil {
		logger.Errorf("Error creating tag %v - %t", tag, err)
		return nil, err
	}

	return tag, nil
}

func (d *fswatchImpl) handleDirectoryTag(directory *model.Directory) (*model.Tag, error) {
	tag := model.Tag{
		ParentID: pointer.Uint64(DIRECTORIES_TAG_ID),
		Title:    directories.DirectoryNameToTag(directory.Path),
	}

	return d.getOrCreateTag(&tag)
}

func (d *fswatchImpl) directoryExcluded(directoryPath string) {
	d.removeExcludedSubDirectories(directoryPath)
	d.removeBelongingItems(directoryPath)
}

func (d *fswatchImpl) removeExcludedSubDirectories(directoryPath string) {
	allDirectories, err := d.db.GetAllDirectories()
	if err != nil {
		logger.Errorf("Error getting all directories %t", err)
		return
	}

	for _, dir := range *allDirectories {
		if dir.Excluded == nil || !*dir.Excluded {
			continue
		}

		if strings.HasPrefix(dir.Path, directoryPath) {
			if err := d.db.RemoveDirectory(dir.Path); err != nil {
				logger.Errorf("Error removing directory %s", dir.Path)
			}
		}
	}
}

func (d *fswatchImpl) removeBelongingItems(directoryPath string) {
	// TODO
}
