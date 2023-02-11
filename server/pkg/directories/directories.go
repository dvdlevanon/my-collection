package directories

import (
	"errors"
	"math"
	"my-collection/server/pkg/gallery"
	itemprocessor "my-collection/server/pkg/item-processor"
	"my-collection/server/pkg/model"
	"my-collection/server/pkg/storage"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/h2non/filetype"
	"github.com/op/go-logging"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gorm.io/gorm"
	"k8s.io/utils/pointer"
)

const DIRECTORIES_TAG_ID = uint64(1000000000) // tags-util.js

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

type directoriesImpl struct {
	gallery         *gallery.Gallery
	storage         *storage.Storage
	processor       itemprocessor.ItemProcessor
	changeChannel   chan model.Directory
	excludedChannel chan string
}

func New(gallery *gallery.Gallery, storage *storage.Storage, processor itemprocessor.ItemProcessor) Directories {
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
			// millisSinceScanned := time.Now().UnixMilli() - directory.LastSynced

			// if millisSinceScanned < 10000 {
			// 	return
			// }

			logger.Infof("Periodic scan")
		}
	}
}

func (d *directoriesImpl) directoryChanged(directory *model.Directory) {
	allDirectories, err := d.gallery.GetAllDirectories()
	if err != nil {
		logger.Errorf("Error getting all directories %t", err)
		return
	}

	tag, err := d.handleDirectoryTag(directory)
	if err != nil {
		return
	}

	directory.FilesCount = pointer.Int(d.syncDirectory(directory, tag, *allDirectories))
	directory.LastSynced = time.Now().UnixMilli()
	directory.ProcessingStart = pointer.Int64(0)
	if err := d.gallery.CreateOrUpdateDirectory(directory); err != nil {
		logger.Errorf("Error updating directory %s %t", directory.Path, err)
	}
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
	for _, dir := range allDirectories {
		if dir.Path == path {
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

func (d *directoriesImpl) addFileIfMissing(directory *model.Directory, tag *model.Tag, path string) (bool, error) {
	exists, lastModified, err := d.itemExists(path, tag)

	if exists || err != nil {
		return false, err
	}

	if !d.isVideo(path) {
		return false, nil
	}

	title := filepath.Base(path)
	item := model.Item{
		Title:        title,
		Origin:       directory.Path,
		Url:          path,
		LastModified: lastModified,
		Tags: []*model.Tag{
			tag,
		},
	}

	logger.Debugf("Adding a new file %s to %v", path, item)

	if err := d.gallery.CreateOrUpdateItem(&item); err != nil {
		logger.Errorf("Error creating item %v - %t", item, err)
		return false, err
	}

	d.processor.EnqueueItemVideoMetadata(item.Id)
	d.processor.EnqueueItemCovers(item.Id)
	d.processor.EnqueueItemPreview(item.Id)

	return true, nil
}

func (d *directoriesImpl) itemExists(path string, tag *model.Tag) (bool, int64, error) {
	title := filepath.Base(path)
	file, err := os.Stat(path)
	if err != nil {
		logger.Errorf("Error getting file stat %s - %t", file, err)
		return false, 0, err
	}

	items, err := d.gallery.GetItemsOfTag(tag)
	if err != nil {
		logger.Errorf("Error getting files of tag %t", err)
	}

	for _, item := range *items {
		if item.Title == title && item.LastModified == file.ModTime().UnixMilli() {
			return true, file.ModTime().UnixMilli(), nil
		}
	}

	return false, file.ModTime().UnixMilli(), nil
}

func (d *directoriesImpl) directoryNameToTag(path string) string {
	caser := cases.Title(language.English)
	return caser.String(strings.ReplaceAll(strings.ReplaceAll(filepath.Base(path), "-", " "), "_", " "))
}

func (d *directoriesImpl) handleDirectoryTag(directory *model.Directory) (*model.Tag, error) {
	tag := model.Tag{
		ParentID: pointer.Uint64(DIRECTORIES_TAG_ID),
		Title:    d.directoryNameToTag(directory.Path),
	}

	existing, err := d.gallery.GetTag(tag)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Errorf("Error handling directory tag %t", err)
		return nil, err
	}

	if existing != nil {
		return existing, nil
	}

	if err := d.gallery.CreateOrUpdateTag(&tag); err != nil {
		logger.Errorf("Error creating tag %v - %t", tag, err)
		return nil, err
	}

	return &tag, nil
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
