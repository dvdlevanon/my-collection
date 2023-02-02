package directories

import (
	"io/ioutil"
	"my-collection/server/pkg/gallery"
	"my-collection/server/pkg/model"
	"my-collection/server/pkg/storage"
	"path/filepath"
	"strings"
	"time"

	"github.com/op/go-logging"
	"k8s.io/utils/pointer"
)

var logger = logging.MustGetLogger("directories")

type Directories interface {
	DirectoryChanged(directory *model.Directory)
	DirectoryRemoved(path string)
}

type DirectoriesImpl struct {
	gallery        *gallery.Gallery
	storage        *storage.Storage
	changeChannel  chan model.Directory
	removedChannel chan string
}

func New(gallery *gallery.Gallery, storage *storage.Storage) Directories {
	logger.Infof("Directories initialized")

	result := &DirectoriesImpl{
		gallery:        gallery,
		storage:        storage,
		changeChannel:  make(chan model.Directory),
		removedChannel: make(chan string),
	}

	go result.watchFilesystemChanges()

	return result
}

func (d *DirectoriesImpl) DirectoryChanged(directory *model.Directory) {
	d.changeChannel <- *directory
}

func (d *DirectoriesImpl) DirectoryRemoved(directoryPath string) {
	d.removedChannel <- directoryPath
}

func (d *DirectoriesImpl) watchFilesystemChanges() {
	for {
		select {
		case directory := <-d.changeChannel:
			logger.Infof("Directory changed %s", directory.Path)

			directories, err := d.gallery.GetAllDirectories()
			if err != nil {
				logger.Errorf("Error getting all directories %t", err)
			}

			d.directoryChanged(&directory, *directories)
		case directoryPath := <-d.removedChannel:
			logger.Infof("Directory removed %s", directoryPath)

			directories, err := d.gallery.GetAllDirectories()
			if err != nil {
				logger.Errorf("Error getting all directories %t", err)
			}

			d.directoryRemoved(directoryPath, *directories)
		case <-time.After(60 * time.Second):
			logger.Infof("Periodic scan")
		}
	}
}

func (d *DirectoriesImpl) directoryRemoved(directoryPath string, allDirectories []model.Directory) {
	for _, dir := range allDirectories {
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

func (d *DirectoriesImpl) directoryChanged(directory *model.Directory, allDirectories []model.Directory) {
	millisSinceScanned := time.Now().UnixMilli() - directory.LastSynced

	if millisSinceScanned < 10000 {
		return
	}

	absolutePath := d.gallery.GetFile(directory.Path)
	files, err := ioutil.ReadDir(absolutePath)
	if err != nil {
		logger.Errorf("Error getting files of %s %t", directory.Path, err)
	}

	filesCount := 0

	for _, file := range files {
		if file.IsDir() {
			path := filepath.Join(directory.Path, file.Name())
			if !directoryExists(path, allDirectories) {
				newDirectory := &model.Directory{
					Path:     path,
					Excluded: pointer.Bool(true),
				}

				if err := d.gallery.CreateOrUpdateDirectory(newDirectory); err != nil {
					logger.Errorf("Error saving directory %s %t", path, err)
					continue
				}
			}
		} else {
			filesCount++
		}
	}

	directory.FilesCount = filesCount
	directory.LastSynced = time.Now().UnixMilli()
	if err := d.gallery.CreateOrUpdateDirectory(directory); err != nil {
		logger.Errorf("Error updating directory %s %t", directory.Path, err)
	}
}

func directoryExists(path string, allDirectories []model.Directory) bool {
	for _, dir := range allDirectories {
		if dir.Path == path {
			return true
		}
	}

	return false
}
