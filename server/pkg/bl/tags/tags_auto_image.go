package tags

import (
	"context"
	"fmt"
	"my-collection/server/pkg/bl/directories"
	"my-collection/server/pkg/model"
	"os"
	"path/filepath"
	"time"

	"github.com/go-errors/errors"
	cp "github.com/otiai10/copy"
	"gorm.io/gorm"
)

func AutoImageChildren(ctx context.Context, storage model.StorageUploader, trw model.TagReaderWriter,
	titrw model.TagImageTypeReaderWriter, tag *model.Tag, directoryPath string) error {
	dirs, err := os.ReadDir(directoryPath)
	if err != nil {
		return errors.Wrap(err, 0)
	}

	for _, slimChildTag := range tag.Children {
		childTag, err := trw.GetTag(ctx, slimChildTag.Id)
		if err != nil {
			logger.Errorf("Error getting tag %v - %t", slimChildTag, err)
		}

		for _, dir := range dirs {
			if !dir.IsDir() {
				continue
			}

			imageTypePath := filepath.Join(directoryPath, dir.Name())
			if err := autoImageTagType(ctx, storage, trw, titrw, childTag, imageTypePath, dir.Name()); err != nil {
				logger.Errorf("Error auto tagging %v from %s - %t", tag, imageTypePath, err)
			}
		}
	}

	return nil
}

func autoImageTagType(ctx context.Context, storage model.StorageUploader, tw model.TagWriter,
	titrw model.TagImageTypeReaderWriter, tag *model.Tag, directoryPath string, nickname string) error {
	tit, err := getOrCreateTagImageType(ctx, titrw, nickname)
	if err != nil {
		logger.Errorf("Unable to get tag image type for %s - %s", nickname, err)
		return err
	}

	if err := updateTagImageTypeIcon(ctx, storage, titrw, tit, directoryPath); err != nil {
		return err
	}

	return autoImageTag(ctx, storage, tw, tag, directoryPath, tit)
}

func updateTagImageTypeIcon(ctx context.Context, storage model.StorageUploader, titrw model.TagImageTypeReaderWriter, tit *model.TagImageType, directoryPath string) error {
	if tit.IconUrl != "" {
		return nil
	}

	path, err := findExistingImage("icon", directoryPath)
	if err != nil {
		return err
	}
	if path == "" {
		return nil
	}

	relativeFile := filepath.Join("tit-icon", fmt.Sprint(tit.Id), filepath.Base(path))
	storageFile, err := storage.GetFileForWriting(relativeFile)
	if err != nil {
		return err
	}

	if err = cp.Copy(path, storageFile); err != nil {
		logger.Errorf("Error coping tit icon %s to %s - %s", path, storageFile, err)
		return nil
	}

	tit.IconUrl = storage.GetStorageUrl(relativeFile)
	return titrw.CreateOrUpdateTagImageType(ctx, tit)
}

func imageExists(tag *model.Tag, tit *model.TagImageType) bool {
	for _, image := range tag.Images {
		if image.ImageTypeId == tit.Id {
			return image.Url != "" && image.Url != "none"
		}
	}

	return false
}

func getOrCreateTagImageType(ctx context.Context, titrw model.TagImageTypeReaderWriter, nickname string) (*model.TagImageType, error) {
	tit, err := titrw.GetTagImageType(ctx, "nickname = ?", nickname)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if err == nil {
		return tit, nil
	}

	tit = &model.TagImageType{
		Nickname: nickname,
	}

	if err := titrw.CreateOrUpdateTagImageType(ctx, tit); err != nil {
		return nil, err
	}

	return tit, nil
}

func autoImageTag(ctx context.Context, storage model.StorageUploader, tw model.TagWriter, tag *model.Tag,
	directoryPath string, tit *model.TagImageType) error {
	path, err := findExistingImage(tag.Title, directoryPath)
	if err != nil {
		return err
	}
	if path == "" {
		return nil
	}

	relativeFile := filepath.Join("tags-image-types", fmt.Sprint(tag.Id), fmt.Sprint(tit.Id), filepath.Base(path))
	storageFile, err := storage.GetFileForWriting(relativeFile)
	if err != nil {
		return err
	}

	if err = cp.Copy(path, storageFile); err != nil {
		logger.Errorf("Error coping %s to %s - %t", path, storageFile, err)
		return nil
	}

	if !imageExists(tag, tit) {
		tag.Images = append(tag.Images, &model.TagImage{
			TagId:       tag.Id,
			Url:         storage.GetStorageUrl(relativeFile),
			ImageNonce:  time.Now().UnixNano(),
			ImageTypeId: tit.Id,
		})
	}

	return tw.CreateOrUpdateTag(ctx, tag)
}

func findExistingImage(tagTitle string, directory string) (string, error) {
	possiblePaths := []string{
		filepath.Join(directory, tagTitle),
		filepath.Join(directory, directories.TagTitleToDirectory(tagTitle)),
	}

	possibleExtenssions := []string{"jpg", "png", "svg"}

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
