package tags

import (
	"fmt"
	"my-collection/server/pkg/bl/directories"
	"my-collection/server/pkg/model"
	"my-collection/server/pkg/storage"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/op/go-logging"
	cp "github.com/otiai10/copy"
)

var logger = logging.MustGetLogger("tags")

func GetItems(ir model.ItemReader, tag *model.Tag) (*[]model.Item, error) {
	itemIds := make([]uint64, 0)
	for _, item := range tag.Items {
		itemIds = append(itemIds, item.Id)
	}

	items, err := ir.GetItems(itemIds)
	if err != nil {
		logger.Errorf("Error getting files of tag %t", err)
		return nil, err
	}

	return items, nil
}

func AutoImageChildren(storage *storage.Storage, tw model.TagWriter, tag *model.Tag, directoryPath string) error {
	for _, childTag := range tag.Children {
		if childTag.Image != "" && childTag.Image != "none" {
			continue
		}

		if err := autoImageTag(storage, tw, childTag, directoryPath); err != nil {
			logger.Errorf("Error auto tagging %v from %s - %t", childTag, directoryPath, err)
		}
	}

	return nil
}

func autoImageTag(storage *storage.Storage, tw model.TagWriter, tag *model.Tag, directoryPath string) error {
	path, err := findExistingImage(tag.Title, directoryPath)
	if err != nil {
		return err
	}

	if path == "" {
		return nil
	}

	fileName := fmt.Sprintf("%s-%s", filepath.Base(path), uuid.NewString())
	relativeFile := filepath.Join("tags-image", fmt.Sprint(tag.Id), fileName)
	storageFile, err := storage.GetFileForWriting(relativeFile)
	if err != nil {
		return err
	}

	if err = cp.Copy(path, storageFile); err != nil {
		logger.Errorf("Error coping %s to %s - %t", path, storageFile, err)
		return nil
	}

	tag.Image = storage.GetStorageUrl(relativeFile)
	return tw.CreateOrUpdateTag(tag)
}

func findExistingImage(tagTitle string, directory string) (string, error) {
	possiblePaths := []string{
		filepath.Join(directory, tagTitle),
		filepath.Join(directory, directories.TagTitleToDirectory(tagTitle)),
	}

	possibleExtenssions := []string{"jpg", "png"}

	for _, pathWithoutExt := range possiblePaths {
		for _, ext := range possibleExtenssions {
			path := fmt.Sprintf("%s.%s", pathWithoutExt, ext)
			if _, err := os.Stat(path); err != nil {
				continue
			}

			return path, nil
		}
	}

	return "", nil
}
