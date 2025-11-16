package tags

import (
	"context"
	"errors"
	"my-collection/server/pkg/model"

	"github.com/op/go-logging"
	"gorm.io/gorm"
)

var logger = logging.MustGetLogger("tags")

func GetItems(ctx context.Context, ir model.ItemReader, tag *model.Tag) (*[]model.Item, error) {
	itemIds := make([]uint64, 0)
	for _, item := range tag.Items {
		itemIds = append(itemIds, item.Id)
	}

	if len(itemIds) == 0 {
		result := make([]model.Item, 0)
		return &result, nil
	}

	items, err := ir.GetItems(ctx, itemIds)
	if err != nil {
		logger.Errorf("Error getting items of tag %t", err)
		return nil, err
	}

	return items, nil
}

func GetItemByTitle(ctx context.Context, ir model.ItemReader, tag *model.Tag, title string) (*model.Item, error) {
	items, err := GetItems(ctx, ir, tag)
	if err != nil {
		return nil, err
	}

	for _, item := range *items {
		if item.Title == title {
			return &item, nil
		}
	}

	return nil, nil
}

func GetOrCreateTag(ctx context.Context, trw model.TagReaderWriter, tag *model.Tag) (*model.Tag, error) {
	existing, err := trw.GetTag(ctx, tag)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Errorf("Error getting tag %t", err)
		return nil, err
	}

	if existing != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return existing, nil
	}

	if err := trw.CreateOrUpdateTag(ctx, tag); err != nil {
		logger.Errorf("Error creating tag %v - %t", tag, err)
		return nil, err
	}

	return tag, nil
}

func GetOrCreateChildTag(ctx context.Context, trw model.TagReaderWriter, parentId uint64, title string) (*model.Tag, error) {
	tag := model.Tag{
		ParentID: &parentId,
		Title:    title,
	}

	return GetOrCreateTag(ctx, trw, &tag)
}

func GetChildTag(ctx context.Context, tr model.TagReader, parentId uint64, title string) (*model.Tag, error) {
	tag := model.Tag{
		ParentID: &parentId,
		Title:    title,
	}

	return tr.GetTag(ctx, tag)
}

func GetOrCreateTags(ctx context.Context, trw model.TagReaderWriter, tags []*model.Tag) ([]*model.Tag, error) {
	result := make([]*model.Tag, 0)

	for _, tag := range tags {
		tag, err := GetOrCreateTag(ctx, trw, tag)

		if err != nil {
			return nil, err
		}

		result = append(result, tag)
	}

	return result, nil
}

func RemoveTagAndItsAssociations(ctx context.Context, tw model.TagWriter, tag *model.Tag) []error {
	errors := make([]error, 0)
	if err := tw.RemoveTag(ctx, tag.Id); err != nil {
		errors = append(errors, err)
	}

	return errors
}

func GetFullTags(ctx context.Context, tr model.TagReader, tagIds []*model.Tag) (*[]model.Tag, error) {
	ids := make([]uint64, len(tagIds))
	for i, tag := range tagIds {
		ids[i] = tag.Id
	}

	if len(ids) == 0 {
		result := make([]model.Tag, 0)
		return &result, nil
	}

	return tr.GetTags(ctx, ids)
}

func GetCategories(ctx context.Context, tr model.TagReader) (*[]model.Tag, error) {
	return tr.GetTags(ctx, "parent_id is NULL")
}

func IsBelongToCategory(tag *model.Tag, category *model.Tag) bool {
	return tag.ParentID != nil && *tag.ParentID == category.Id
}

func RemoveTagImages(ctx context.Context, trw model.TagReaderWriter, tagId uint64, titId uint64) error {
	tag, err := trw.GetTag(ctx, tagId)
	if err != nil {
		return err
	}

	for _, image := range tag.Images {
		if image.ImageTypeId != titId {
			continue
		}

		if err := trw.RemoveTagImageFromTag(ctx, tagId, image.Id); err != nil {
			return err
		}
	}

	return nil
}
