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
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gorm.io/gorm"
	"k8s.io/utils/pointer"
)

const ROOT_DIRECTORY_PATH = "<root>"

var directoriesTag = &model.Tag{
	Title:    "Directories", // tags-utils.js
	ParentID: nil,
}

var logger = logging.MustGetLogger("directories")

func Init(trw model.TagReaderWriter) error {
	d, err := trw.GetTag(directoriesTag)
	if err != nil {
		if err := trw.CreateOrUpdateTag(directoriesTag); err != nil {
			return err
		}
	} else {
		directoriesTag = d
	}

	return nil
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
		if TagExists(tags, tag) {
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

func AddDirectory(dw model.DirectoryWriter, dir string, excluded bool) error {
	newDirectory := &model.Directory{
		Path:     NormalizeDirectoryPath(dir),
		Excluded: pointer.Bool(excluded),
	}

	return dw.CreateOrUpdateDirectory(newDirectory)
}

func AddDirectoryIfMissing(drw model.DirectoryReaderWriter, dir string, excluded bool) error {
	exists, err := DirectoryExists(drw, dir)
	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	return AddDirectory(drw, dir, excluded)
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
	directory.Excluded = pointer.Bool(false)
	directory.Path = NormalizeDirectoryPath(directory.Path)
	return dw.CreateOrUpdateDirectory(directory)
}

func UpdatePath(dw model.DirectoryWriter, directory *model.Directory, newpath string) error {
	oldPath := directory.Path
	directory.Path = NormalizeDirectoryPath(newpath)
	if err := dw.CreateOrUpdateDirectory(directory); err != nil {
		return err
	}

	return dw.RemoveDirectory(oldPath)
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

func ValidateReadyDirectory(drw model.DirectoryReaderWriter, path string) (*model.Directory, error) {
	dirpath := NormalizeDirectoryPath(path)
	if err := AddDirectoryIfMissing(drw, dirpath, false); err != nil {
		return nil, err
	}

	dir, err := GetDirectory(drw, dirpath)
	if err != nil {
		return nil, err
	}

	if IsExcluded(dir) {
		dir.Excluded = pointer.Bool(false)
		return dir, drw.CreateOrUpdateDirectory(dir)
	}

	return dir, nil
}

func GetDirectoriesTagId() uint64 {
	return directoriesTag.Id
}
