package directories

import (
	"context"
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

var directoriesTag = &model.Tag{
	Title:    "Directories", // tags-utils.js
	ParentID: nil,
}

var logger = logging.MustGetLogger("directories")

func Init(ctx context.Context, trw model.TagReaderWriter) error {
	d, err := trw.GetTag(ctx, directoriesTag)
	if err != nil {
		if err := trw.CreateOrUpdateTag(ctx, directoriesTag); err != nil {
			return err
		}
	} else {
		directoriesTag = d
	}

	return nil
}

func ExcludeDirectory(ctx context.Context, drw model.DirectoryReaderWriter, path string) error {
	directory, err := GetDirectory(ctx, drw, NormalizeDirectoryPath(path))
	if err != nil {
		return err
	}

	if *directory.Excluded {
		return nil
	}

	directory.Excluded = pointer.Bool(true)
	directory.AutoIncludeChildren = pointer.Bool(false)
	directory.AutoIncludeHierarchy = pointer.Bool(false)
	return drw.CreateOrUpdateDirectory(ctx, directory)
}

func IncludeOrCreateDirectory(ctx context.Context, drw model.DirectoryReaderWriter, path string) error {
	if err := IncludeDirectory(ctx, drw, path); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := CreateOrUpdateDirectory(ctx, drw, &model.Directory{Path: path, Excluded: pointer.Bool(false)}); err != nil {
				return err
			}
		} else {
			return err
		}
	}

	return nil
}

func AutoIncludeHierarchy(ctx context.Context, drw model.DirectoryReaderWriter, path string) error {
	directory, err := GetDirectory(ctx, drw, NormalizeDirectoryPath(path))
	if err != nil {
		return err
	}

	if directory.AutoIncludeHierarchy != nil && *directory.AutoIncludeHierarchy {
		return nil
	}

	directory.AutoIncludeHierarchy = pointer.Bool(true)
	return drw.CreateOrUpdateDirectory(ctx, directory)
}

func AutoIncludeChildren(ctx context.Context, drw model.DirectoryReaderWriter, path string) error {
	directory, err := GetDirectory(ctx, drw, NormalizeDirectoryPath(path))
	if err != nil {
		return err
	}

	if directory.AutoIncludeChildren != nil && *directory.AutoIncludeChildren {
		return nil
	}

	directory.AutoIncludeChildren = pointer.Bool(true)
	return drw.CreateOrUpdateDirectory(ctx, directory)
}

func IncludeDirectory(ctx context.Context, drw model.DirectoryReaderWriter, path string) error {
	directory, err := GetDirectory(ctx, drw, NormalizeDirectoryPath(path))
	if err != nil {
		return err
	}

	if !(*directory.Excluded) {
		return nil
	}

	directory.Excluded = pointer.Bool(false)
	return drw.CreateOrUpdateDirectory(ctx, directory)
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

func RemoveMissingTags(ctx context.Context, drw model.DirectoryReaderWriter, directory *model.Directory, tags []*model.Tag) {
	for _, tag := range directory.Tags {
		if TagExists(tags, tag) {
			continue
		}

		if err := drw.RemoveTagFromDirectory(ctx, directory.Path, tag.Id); err != nil {
			logger.Warningf("Unable to remove tag %d from directory %s - %t",
				directory.Path, tag.Id, err)
		}
	}
}

func UpdateDirectoryTags(ctx context.Context, drw model.DirectoryReaderWriter, directory *model.Directory) error {
	existingDirectory, err := drw.GetDirectory(ctx, "path = ?", directory.Path)
	if err != nil {
		logger.Errorf("Error getting exising directory %s %t", directory.Path, err)
		return err
	}

	RemoveMissingTags(ctx, drw, existingDirectory, directory.Tags)

	return drw.CreateOrUpdateDirectory(ctx, directory)
}

func GetDirectory(ctx context.Context, dr model.DirectoryReader, path string) (*model.Directory, error) {
	return dr.GetDirectory(ctx, "path = ?", NormalizeDirectoryPath(path))
}

func DirectoryExists(ctx context.Context, dr model.DirectoryReader, path string) (bool, error) {
	_, err := GetDirectory(ctx, dr, path)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func AddDirectory(ctx context.Context, dw model.DirectoryWriter, dir string, excluded bool) error {
	newDirectory := &model.Directory{
		Path:     NormalizeDirectoryPath(dir),
		Excluded: pointer.Bool(excluded),
	}

	return dw.CreateOrUpdateDirectory(ctx, newDirectory)
}

func ShouldInclude(ctx context.Context, dr model.DirectoryReader, path string) bool {
	dir, err := GetParent(ctx, dr, path)
	if err != nil {
		return false
	}

	if dir.AutoIncludeChildren != nil && *dir.AutoIncludeChildren {
		return true
	}

	for {
		if dir.AutoIncludeHierarchy != nil && *dir.AutoIncludeHierarchy {
			return true
		}

		dir, err = GetParent(ctx, dr, dir.Path)
		if err != nil {
			return false
		}
		if dir == nil {
			return false
		}
	}
}

func GetParent(ctx context.Context, dr model.DirectoryReader, path string) (*model.Directory, error) {
	if path == model.ROOT_DIRECTORY_PATH {
		return nil, nil
	}
	parentPath := filepath.Dir(path)
	return GetDirectory(ctx, dr, parentPath)
}

func AddDirectoryIfMissing(ctx context.Context, drw model.DirectoryReaderWriter, dir string, excluded bool) error {
	exists, err := DirectoryExists(ctx, drw, dir)
	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	return AddDirectory(ctx, drw, dir, excluded)
}

func BuildDirectoryTags(directory *model.Directory) []*model.Tag {
	result := make([]*model.Tag, 0)
	title := DirectoryNameToTag(directory.Path)
	for _, directoryTag := range directory.Tags {
		result = append(result, &model.Tag{ParentID: &directoryTag.Id, Title: title})
	}

	return result
}

func AddRootDirectory(ctx context.Context, drw model.DirectoryReaderWriter) error {
	return AddDirectoryIfMissing(ctx, drw, model.ROOT_DIRECTORY_PATH, false)
}

func NormalizeDirectoryPath(path string) string {
	normalizedPath := relativasor.GetRelativePath(path)

	if normalizedPath == "" {
		return model.ROOT_DIRECTORY_PATH
	}

	return normalizedPath
}

func CreateOrUpdateDirectory(ctx context.Context, dw model.DirectoryWriter, directory *model.Directory) error {
	directory.Excluded = pointer.Bool(false)
	directory.Path = NormalizeDirectoryPath(directory.Path)
	return dw.CreateOrUpdateDirectory(ctx, directory)
}

func UpdatePath(ctx context.Context, dw model.DirectoryWriter, directory *model.Directory, newpath string) error {
	oldPath := directory.Path
	directory.Path = NormalizeDirectoryPath(newpath)
	if err := dw.CreateOrUpdateDirectory(ctx, directory); err != nil {
		return err
	}

	return dw.RemoveDirectory(ctx, oldPath)
}

func StartDirectoryProcessing(ctx context.Context, dw model.DirectoryWriter, directory *model.Directory) error {
	if directory.ProcessingStart != nil && *directory.ProcessingStart != 0 {
		return nil
	}

	directory.ProcessingStart = pointer.Int64(time.Now().UnixMilli())
	return dw.CreateOrUpdateDirectory(ctx, directory)
}

func FinishDirectoryProcessing(ctx context.Context, dw model.DirectoryWriter, directory *model.Directory) error {
	directory.LastSynced = time.Now().UnixMilli()
	directory.ProcessingStart = pointer.Int64(0)
	return dw.CreateOrUpdateDirectory(ctx, directory)
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

func ValidateReadyDirectory(ctx context.Context, drw model.DirectoryReaderWriter, path string) (*model.Directory, error) {
	dirpath := NormalizeDirectoryPath(path)
	if err := AddDirectoryIfMissing(ctx, drw, dirpath, false); err != nil {
		return nil, err
	}

	dir, err := GetDirectory(ctx, drw, dirpath)
	if err != nil {
		return nil, err
	}

	if IsExcluded(dir) {
		dir.Excluded = pointer.Bool(false)
		return dir, drw.CreateOrUpdateDirectory(ctx, dir)
	}

	return dir, nil
}

func GetDirectoriesTagId() uint64 {
	return directoriesTag.Id
}

func EnrichFsNode(ctx context.Context, dr model.DirectoryReader, node *model.FsNode) (*model.FsNode, error) {
	dir, err := GetDirectory(ctx, dr, node.Path)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return node, nil
		}
		return nil, err
	}

	node.DirInfo = dir

	for i := 0; i < len(node.Children); i++ {
		child, err := EnrichFsNode(ctx, dr, node.Children[i])
		if err != nil {
			return nil, err
		}
		node.Children[i] = child
	}

	return node, nil
}
