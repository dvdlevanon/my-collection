package tags

import (
	"errors"
	"my-collection/server/pkg/model"

	"github.com/op/go-logging"
	"gorm.io/gorm"
)

var logger = logging.MustGetLogger("tags")

func GetItems(ir model.ItemReader, tag *model.Tag) (*[]model.Item, error) {
	itemIds := make([]uint64, 0)
	for _, item := range tag.Items {
		itemIds = append(itemIds, item.Id)
	}

	if len(itemIds) == 0 {
		result := make([]model.Item, 0)
		return &result, nil
	}

	items, err := ir.GetItems(itemIds)
	if err != nil {
		logger.Errorf("Error getting items of tag %t", err)
		return nil, err
	}

	return items, nil
}

func GetItemByTitle(ir model.ItemReader, tag *model.Tag, title string) (*model.Item, error) {
	items, err := GetItems(ir, tag)
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

func GetOrCreateTag(trw model.TagReaderWriter, tag *model.Tag) (*model.Tag, error) {
	existing, err := trw.GetTag(tag)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Errorf("Error getting tag %t", err)
		return nil, err
	}

	if existing != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return existing, nil
	}

	if err := trw.CreateOrUpdateTag(tag); err != nil {
		logger.Errorf("Error creating tag %v - %t", tag, err)
		return nil, err
	}

	return tag, nil
}

func GetOrCreateChildTag(trw model.TagReaderWriter, parentId uint64, title string) (*model.Tag, error) {
	tag := model.Tag{
		ParentID: &parentId,
		Title:    title,
	}

	return GetOrCreateTag(trw, &tag)
}

func GetChildTag(tr model.TagReader, parentId uint64, title string) (*model.Tag, error) {
	tag := model.Tag{
		ParentID: &parentId,
		Title:    title,
	}

	return tr.GetTag(tag)
}

func GetOrCreateTags(trw model.TagReaderWriter, tags []*model.Tag) ([]*model.Tag, error) {
	result := make([]*model.Tag, 0)

	for _, tag := range tags {
		tag, err := GetOrCreateTag(trw, tag)

		if err != nil {
			return nil, err
		}

		result = append(result, tag)
	}

	return result, nil
}

func RemoveTagAndItsAssociations(tw model.TagWriter, tag *model.Tag) []error {
	errors := make([]error, 0)
	if err := tw.RemoveTag(tag.Id); err != nil {
		errors = append(errors, err)
	}

	return errors
}

func GetFullTags(tr model.TagReader, tagIds []*model.Tag) (*[]model.Tag, error) {
	ids := make([]uint64, len(tagIds))
	for i, tag := range tagIds {
		ids[i] = tag.Id
	}

	if len(ids) == 0 {
		result := make([]model.Tag, 0)
		return &result, nil
	}

	return tr.GetTags(ids)
}

func GetCategories(tr model.TagReader) (*[]model.Tag, error) {
	return tr.GetTags("parent_id is NULL")
}

func IsBelongToCategory(tag *model.Tag, category *model.Tag) bool {
	return tag.ParentID != nil && *tag.ParentID == category.Id
}

func RemoveTagImages(trw model.TagReaderWriter, tagId uint64, titId uint64) error {
	tag, err := trw.GetTag(tagId)
	if err != nil {
		return err
	}

	for _, image := range tag.Images {
		if image.ImageTypeId != titId {
			continue
		}

		if err := trw.RemoveTagImageFromTag(tagId, image.Id); err != nil {
			return err
		}
	}

	return nil
}
