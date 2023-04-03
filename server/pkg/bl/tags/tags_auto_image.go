package tags

import (
	"fmt"
	"my-collection/server/pkg/bl/directories"
	"my-collection/server/pkg/model"
	"my-collection/server/pkg/storage"
	"os"
	"path/filepath"

	"github.com/go-errors/errors"
	"github.com/google/uuid"
	cp "github.com/otiai10/copy"
	"gorm.io/gorm"
)

func AutoImageChildren(storage *storage.Storage, tw model.TagWriter,
	titrw model.TagImageTypeReaderWriter, tag *model.Tag, directoryPath string) error {
	dirs, err := os.ReadDir(directoryPath)
	if err != nil {
		return errors.Wrap(err, 0)
	}

	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}

		for _, childTag := range tag.Children {
			imageTypePath := filepath.Join(directoryPath, dir.Name())
			if err := autoImageTagType(storage, tw, titrw, childTag, imageTypePath, dir.Name()); err != nil {
				logger.Errorf("Error auto tagging %v from %s - %t", tag, imageTypePath, err)
			}
		}
	}

	return nil
}

func autoImageTagType(storage *storage.Storage, tw model.TagWriter,
	titrw model.TagImageTypeReaderWriter, tag *model.Tag, directoryPath string, nickname string) error {
	tit, err := getOrCreateTagImageType(titrw, nickname)
	if err != nil {
		logger.Errorf("Unable to get tag image type for %s - %s", nickname, err)
		return err
	}

	if imageExists(tag, tit) {
		return nil
	}

	return autoImageTag(storage, tw, tag, directoryPath, tit)
}

func imageExists(tag *model.Tag, tit *model.TagImageType) bool {
	for _, image := range tag.Images {
		if image.ImageTypeId == tit.Id {
			return image.Url != "" && image.Url != "none"
		}
	}

	return false
}

func getOrCreateTagImageType(titrw model.TagImageTypeReaderWriter, nickname string) (*model.TagImageType, error) {
	tit, err := titrw.GetTagImageType("nickname = ?", nickname)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if err == nil {
		return tit, nil
	}

	tit = &model.TagImageType{
		Nickname: nickname,
	}

	if err := titrw.CreateOrUpdateTagImageType(tit); err != nil {
		return nil, err
	}

	return tit, nil
}

func autoImageTag(storage *storage.Storage, tw model.TagWriter, tag *model.Tag,
	directoryPath string, tit *model.TagImageType) error {
	path, err := findExistingImage(tag.Title, directoryPath)
	if err != nil {
		return err
	}

	if path == "" {
		return nil
	}

	fileName := fmt.Sprintf("%s-%s", filepath.Base(path), uuid.NewString())
	relativeFile := filepath.Join("tags-image-types", fmt.Sprint(tag.Id), fmt.Sprint(tit.Id), fileName)
	storageFile, err := storage.GetFileForWriting(relativeFile)
	if err != nil {
		return err
	}

	if err = cp.Copy(path, storageFile); err != nil {
		logger.Errorf("Error coping %s to %s - %t", path, storageFile, err)
		return nil
	}

	tag.Images = append(tag.Images, &model.TagImage{
		TagId:       tag.Id,
		Url:         storage.GetStorageUrl(relativeFile),
		ImageTypeId: tit.Id,
	})

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
