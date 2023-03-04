package directories

import (
	"errors"
	"math"
	"my-collection/server/pkg/gallery"
	"my-collection/server/pkg/model"
	processor "my-collection/server/pkg/processor"
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

var logger = logging.MustGetLogger("directories")
var directoriesTag = model.Tag{
	Id:    DIRECTORIES_TAG_ID,
	Title: "Directories",
}

type Directories interface {
	Init() error
	DirectoryChanged(directory *model.Directory)
	DirectoryExcluded(path string)
}

type DirectoriesMock struct{}

func (d *DirectoriesMock) Init() error                                 { return nil }
func (d *DirectoriesMock) DirectoryChanged(directory *model.Directory) {}
func (d *DirectoriesMock) DirectoryExcluded(path string)               {}

type directoriesImpl struct {
	gallery         *gallery.Gallery
	storage         *storage.Storage
	processor       processor.Processor
	changeChannel   chan model.Directory
	excludedChannel chan string
}

func New(gallery *gallery.Gallery, storage *storage.Storage, processor processor.Processor) Directories {
	logger.Infof("Directories initialized")

	return &directoriesImpl{
		gallery:         gallery,
		storage:         storage,
		processor:       processor,
		changeChannel:   make(chan model.Directory),
		excludedChannel: make(chan string),
	}
}

func (d *directoriesImpl) Init() error {
	if err := d.gallery.CreateOrUpdateTag(&directoriesTag); err != nil {
		return err
	}

	go d.watchFilesystemChanges()
	return nil
}

func (d *directoriesImpl) DirectoryChanged(directory *model.Directory) {
	d.changeChannel <- *directory
}

func (d *directoriesImpl) DirectoryExcluded(directoryPath string) {
	d.excludedChannel <- directoryPath
}

func (d *directoriesImpl) watchFilesystemChanges() {
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

func (d *directoriesImpl) periodicScan() {
	allDirectories, err := d.gallery.GetAllDirectories()
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

func (d *directoriesImpl) directoryChanged(directory *model.Directory) {
	allDirectories, err := d.gallery.GetAllDirectories()
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
	if err := d.gallery.CreateOrUpdateDirectory(directory); err != nil {
		logger.Errorf("Error updating directory %s %t", directory.Path, err)
	}
}

func (d *directoriesImpl) handleExcludedDirectory(directory *model.Directory, allDirectories []model.Directory) {
	tag, err := d.gallery.GetTag(model.Tag{
		ParentID: pointer.Uint64(DIRECTORIES_TAG_ID),
		Title:    d.gallery.DirectoryNameToTag(directory.Path),
	})

	if err != nil {
		logger.Errorf("Unable to find directory of %s - %s", directory.Path, err)
		return
	}

	for _, item := range tag.Items {
		if len(item.Tags) > 1 {
			continue
		}

		if err := d.gallery.RemoveItem(item.Id); err != nil {
			logger.Errorf("Unable to remove item %d - %s", item.Id, err)
		}
	}

	if err := d.gallery.RemoveTag(tag.Id); err != nil {
		logger.Errorf("Unable to remove tag %d - %s", tag.Id, err)
	}
}

func (d *directoriesImpl) handleIncludedDirectory(directory *model.Directory, allDirectories []model.Directory) {
	tag, err := d.handleDirectoryTag(directory)
	if err != nil {
		return
	}

	directory.FilesCount = pointer.Int(d.syncDirectory(directory, tag, allDirectories))
}

func (d *directoriesImpl) syncDirectory(directory *model.Directory, tag *model.Tag, allDirectories []model.Directory) int {
	path := d.gallery.GetFile(directory.Path)
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
				tag, err = d.gallery.GetTag(tag.Id)
				if err != nil {
					logger.Errorf("Error refetching tag from DB %t", err)
				}
			}
		}
	}

	items, err := d.gallery.GetItemsOfTag(tag)
	if err != nil {
		logger.Errorf("Error getting files of tag %t", err)
	}

	for _, item := range *items {
		if d.fileExists(directory, item) {
			continue
		}

		if err := d.gallery.RemoveItem(item.Id); err != nil {
			logger.Errorf("Error removing item %d - %t", item.Id, err)
		}
	}

	return filesCount
}

func (d *directoriesImpl) fileExists(directory *model.Directory, item model.Item) bool {
	if item.Origin != directory.Path {
		return true
	}

	path := d.gallery.GetFile(filepath.Join(item.Origin, item.Title))
	_, err := os.Stat(path)
	return err == nil
}

func (d *directoriesImpl) isVideo(path string) bool {
	if d.gallery.TrustFileExtenssion {
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

func (d *directoriesImpl) addDirectoryIfMissing(path string, allDirectories []model.Directory) {
	if d.directoryExists(path, allDirectories) {
		return
	}

	if err := d.addExcludedDirectory(path); err != nil {
		logger.Errorf("Error saving directory %s %t", path, err)
	}
}

func (d *directoriesImpl) directoryExists(path string, allDirectories []model.Directory) bool {
	relativePath := d.gallery.GetRelativePath(path)
	for _, dir := range allDirectories {
		if dir.Path == relativePath {
			return true
		}
	}

	return false
}

func (d *directoriesImpl) addExcludedDirectory(path string) error {
	newDirectory := &model.Directory{
		Path:     path,
		Excluded: pointer.Bool(true),
	}

	return d.gallery.CreateOrUpdateDirectory(newDirectory)
}

func (d *directoriesImpl) getConcreteTagOfDirectory(directory *model.Directory, directoryTag *model.Tag) (*model.Tag, error) {
	tag := model.Tag{
		ParentID: &directoryTag.Id,
		Title:    d.gallery.DirectoryNameToTag(directory.Path),
	}

	return d.getOrCreateTag(&tag)
}

func (d *directoriesImpl) getConcreteTagsOfDirectory(directory *model.Directory) ([]*model.Tag, error) {
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

func (d *directoriesImpl) ensureConcreteTagsOnItem(item *model.Item, concreteTags []*model.Tag) error {
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
	if err := d.gallery.CreateOrUpdateItem(item); err != nil {
		logger.Errorf("Error updating item %v - %t", item, err)
		return err
	}

	return nil
}

func (d *directoriesImpl) addFileIfMissing(directory *model.Directory, tag *model.Tag, path string) (bool, error) {
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

	if err := d.gallery.CreateOrUpdateItem(&item); err != nil {
		logger.Errorf("Error creating item %v - %t", item, err)
		return false, err
	}

	if d.gallery.AutomaticProcessing {
		d.processor.EnqueueItemVideoMetadata(item.Id)
		d.processor.EnqueueItemCovers(item.Id)
		d.processor.EnqueueItemPreview(item.Id)
	}

	return true, nil
}

func (d *directoriesImpl) itemExists(path string, tag *model.Tag) (*model.Item, int64, error) {
	title := filepath.Base(path)
	file, err := os.Stat(path)
	if err != nil {
		logger.Errorf("Error getting file stat %s - %t", file, err)
		return nil, 0, err
	}

	items, err := d.gallery.GetItemsOfTag(tag)
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

func (d *directoriesImpl) getOrCreateTag(tag *model.Tag) (*model.Tag, error) {
	existing, err := d.gallery.GetTag(tag)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Errorf("Error getting tag %t", err)
		return nil, err
	}

	if existing != nil {
		return existing, nil
	}

	if err := d.gallery.CreateOrUpdateTag(tag); err != nil {
		logger.Errorf("Error creating tag %v - %t", tag, err)
		return nil, err
	}

	return tag, nil
}

func (d *directoriesImpl) handleDirectoryTag(directory *model.Directory) (*model.Tag, error) {
	tag := model.Tag{
		ParentID: pointer.Uint64(DIRECTORIES_TAG_ID),
		Title:    d.gallery.DirectoryNameToTag(directory.Path),
	}

	return d.getOrCreateTag(&tag)
}

func (d *directoriesImpl) directoryExcluded(directoryPath string) {
	d.removeExcludedSubDirectories(directoryPath)
	d.removeBelongingItems(directoryPath)
}

func (d *directoriesImpl) removeExcludedSubDirectories(directoryPath string) {
	allDirectories, err := d.gallery.GetAllDirectories()
	if err != nil {
		logger.Errorf("Error getting all directories %t", err)
		return
	}

	for _, dir := range *allDirectories {
		if dir.Excluded == nil || !*dir.Excluded {
			continue
		}

		if strings.HasPrefix(dir.Path, directoryPath) {
			if err := d.gallery.RemoveDirectory(dir.Path); err != nil {
				logger.Errorf("Error removing directory %s", dir.Path)
			}
		}
	}
}

func (d *directoriesImpl) removeBelongingItems(directoryPath string) {
	// TODO
}
