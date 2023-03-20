package directories

import (
	"errors"
	"my-collection/server/pkg/model"
	"my-collection/server/pkg/relativasor"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/op/go-logging"
	"github.com/patrickmn/go-cache"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gorm.io/gorm"
	"k8s.io/utils/pointer"
)

const DIRECTORIES_TAG_ID = uint64(1) // tags-util.js
const ROOT_DIRECTORY_PATH = "<root>"

var DirectoriesTag = model.Tag{
	Id:    DIRECTORIES_TAG_ID,
	Title: "Directories",
}

var logger = logging.MustGetLogger("directories")
var directoriesCache = cache.New(time.Second*3, time.Second)
var directoriesCacheKey = "all-directories"

func GetAllDirectoriesWithCache(dr model.DirectoryReader) (*[]model.Directory, error) {
	allDirectoriesIfc, found := directoriesCache.Get(directoriesCacheKey)
	allDirectories, ok := allDirectoriesIfc.(*[]model.Directory)

	if found && ok {
		return allDirectories, nil
	}

	allDirectories, err := dr.GetAllDirectories()
	if err != nil {
		return nil, err
	}

	directoriesCache.Add(directoriesCacheKey, allDirectories, cache.DefaultExpiration)
	return allDirectories, nil
}

func ExcludeDirectory(drw model.DirectoryReaderWriter, path string) error {
	directory, err := GetDirectory(drw, path)
	if err != nil {
		return err
	}

	if *directory.Excluded {
		return nil
	}

	directory.Excluded = pointer.Bool(true)
	return drw.CreateOrUpdateDirectory(directory)
}

func IncludeDirectory(drw model.DirectoryReaderWriter, path string) error {
	directory, err := GetDirectory(drw, path)
	if err != nil {
		return err
	}

	if !(*directory.Excluded) {
		return nil
	}

	directory.Excluded = pointer.Bool(false)
	return drw.CreateOrUpdateDirectory(directory)
}

func DirectoryNameToTag(path string) string {
	caser := cases.Title(language.English)
	return caser.String(strings.ReplaceAll(strings.ReplaceAll(filepath.Base(path), "-", " "), "_", " "))
}

func TagTitleToDirectory(title string) string {
	return strings.ToLower(strings.ReplaceAll(title, " ", "-"))
}

func TagExists(tags []*model.Tag, tag *model.Tag) bool {
	for _, t := range tags {
		if tag.Id == t.Id {
			return true
		}
	}

	return false
}

func RemoveMissingTags(drw model.DirectoryReaderWriter, directory *model.Directory, tags []*model.Tag) {
	for _, tag := range directory.Tags {
		if TagExists(directory.Tags, tag) {
			continue
		}

		if err := drw.RemoveTagFromDirectory(directory.Path, tag.Id); err != nil {
			logger.Warningf("Unable to remove tag %d from directory %s - %t",
				directory.Path, tag.Id, err)
		}
	}
}

func UpdateDirectoryTags(drw model.DirectoryReaderWriter, directory *model.Directory) error {
	existingDirectory, err := drw.GetDirectory("path = ?", directory.Path)
	if err != nil {
		logger.Errorf("Error getting exising directory %s %t", directory.Path, err)
		return err
	}

	RemoveMissingTags(drw, existingDirectory, directory.Tags)

	return drw.CreateOrUpdateDirectory(directory)
}

func GetDirectory(dr model.DirectoryReader, path string) (*model.Directory, error) {
	return dr.GetDirectory("path = ?", NormalizeDirectoryPath(path))
}

func DirectoryExists(dr model.DirectoryReader, path string) (bool, error) {
	_, err := GetDirectory(dr, path)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func AddExcludedDirectory(dw model.DirectoryWriter, dir string) error {
	newDirectory := &model.Directory{
		Path:     NormalizeDirectoryPath(dir),
		Excluded: pointer.Bool(true),
	}

	return dw.CreateOrUpdateDirectory(newDirectory)
}

func AddExcludedDirectoryIfMissing(drw model.DirectoryReaderWriter, dir string) error {
	exists, err := DirectoryExists(drw, dir)
	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	return AddExcludedDirectory(drw, dir)
}

func BuildDirectoryTags(directory *model.Directory) []*model.Tag {
	result := make([]*model.Tag, 0)
	title := DirectoryNameToTag(directory.Path)
	for _, directoryTag := range directory.Tags {
		result = append(result, &model.Tag{ParentID: &directoryTag.Id, Title: title})
	}

	return result
}

func NormalizeDirectoryPath(path string) string {
	normalizedPath := relativasor.GetRelativePath(path)

	if normalizedPath == "" {
		return ROOT_DIRECTORY_PATH
	}

	return normalizedPath
}

func CreateOrUpdateDirectory(dw model.DirectoryWriter, directory *model.Directory) error {
	directory.ProcessingStart = pointer.Int64(time.Now().UnixMilli())
	directory.Path = NormalizeDirectoryPath(directory.Path)
	return dw.CreateOrUpdateDirectory(directory)
}

func StartDirectoryProcessing(dw model.DirectoryWriter, directory *model.Directory) error {
	if directory.ProcessingStart != nil && *directory.ProcessingStart != 0 {
		return nil
	}

	directory.ProcessingStart = pointer.Int64(time.Now().UnixMilli())
	return dw.CreateOrUpdateDirectory(directory)
}

func FinishDirectoryProcessing(dw model.DirectoryWriter, directory *model.Directory) error {
	directory.LastSynced = time.Now().UnixMilli()
	directory.ProcessingStart = pointer.Int64(0)
	return dw.CreateOrUpdateDirectory(directory)
}

func GetDirectoryFiles(directory *model.Directory) ([]os.DirEntry, error) {
	path := relativasor.GetAbsoluteFile(directory.Path)
	files, err := os.ReadDir(path)
	if err != nil {
		logger.Errorf("Error getting directory files %s %s", path, err)
	}

	return files, err
}

func GetDirectoryFile(directory *model.Directory, filename string) string {
	path := relativasor.GetAbsoluteFile(directory.Path)
	return filepath.Join(path, filename)
}

func IsExcluded(directory *model.Directory) bool {
	if directory.Excluded == nil {
		return false
	}

	return *directory.Excluded
}
