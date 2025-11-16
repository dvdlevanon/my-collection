package tags

import (
	"context"
	"errors"
	"my-collection/server/pkg/model"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

func TestGetItems(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockItemReader := model.NewMockItemReader(ctrl)

	t.Run("tag with items", func(t *testing.T) {
		items := []*model.Item{
			{Id: 1, Title: "Item 1"},
			{Id: 2, Title: "Item 2"},
			{Id: 3, Title: "Item 3"},
		}
		tag := &model.Tag{
			Id:    100,
			Title: "Test Tag",
			Items: items,
		}

		expectedItemIds := []uint64{1, 2, 3}
		expectedItems := &[]model.Item{
			{Id: 1, Title: "Item 1"},
			{Id: 2, Title: "Item 2"},
			{Id: 3, Title: "Item 3"},
		}

		mockItemReader.EXPECT().
			GetItems(gomock.Any(), expectedItemIds).
			Return(expectedItems, nil)

		result, err := GetItems(context.Background(), mockItemReader, tag)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, *result, 3)
		assert.Equal(t, expectedItems, result)
	})

	t.Run("tag with no items", func(t *testing.T) {
		tag := &model.Tag{
			Id:    100,
			Title: "Empty Tag",
			Items: []*model.Item{},
		}

		result, err := GetItems(context.Background(), mockItemReader, tag)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, *result, 0)
	})

	t.Run("error getting items", func(t *testing.T) {
		items := []*model.Item{
			{Id: 1, Title: "Item 1"},
		}
		tag := &model.Tag{
			Id:    100,
			Title: "Test Tag",
			Items: items,
		}

		expectedError := errors.New("database error")
		mockItemReader.EXPECT().
			GetItems(gomock.Any(), []uint64{1}).
			Return(nil, expectedError)

		result, err := GetItems(context.Background(), mockItemReader, tag)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedError, err)
	})
}

func TestGetItemByTitle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockItemReader := model.NewMockItemReader(ctrl)

	t.Run("item found by title", func(t *testing.T) {
		items := []*model.Item{
			{Id: 1, Title: "First Item"},
			{Id: 2, Title: "Target Item"},
			{Id: 3, Title: "Third Item"},
		}
		tag := &model.Tag{
			Id:    100,
			Title: "Test Tag",
			Items: items,
		}

		expectedItems := &[]model.Item{
			{Id: 1, Title: "First Item"},
			{Id: 2, Title: "Target Item"},
			{Id: 3, Title: "Third Item"},
		}

		mockItemReader.EXPECT().
			GetItems(gomock.Any(), []uint64{1, 2, 3}).
			Return(expectedItems, nil)

		result, err := GetItemByTitle(context.Background(), mockItemReader, tag, "Target Item")

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, uint64(2), result.Id)
		assert.Equal(t, "Target Item", result.Title)
	})

	t.Run("item not found by title", func(t *testing.T) {
		items := []*model.Item{
			{Id: 1, Title: "First Item"},
			{Id: 2, Title: "Second Item"},
		}
		tag := &model.Tag{
			Id:    100,
			Title: "Test Tag",
			Items: items,
		}

		expectedItems := &[]model.Item{
			{Id: 1, Title: "First Item"},
			{Id: 2, Title: "Second Item"},
		}

		mockItemReader.EXPECT().
			GetItems(gomock.Any(), []uint64{1, 2}).
			Return(expectedItems, nil)

		result, err := GetItemByTitle(context.Background(), mockItemReader, tag, "Non-existent Item")

		assert.NoError(t, err)
		assert.Nil(t, result)
	})

	t.Run("error in GetItems", func(t *testing.T) {
		items := []*model.Item{
			{Id: 1, Title: "Item 1"},
		}
		tag := &model.Tag{
			Id:    100,
			Title: "Test Tag",
			Items: items,
		}

		expectedError := errors.New("database error")
		mockItemReader.EXPECT().
			GetItems(gomock.Any(), []uint64{1}).
			Return(nil, expectedError)

		result, err := GetItemByTitle(context.Background(), mockItemReader, tag, "Item 1")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedError, err)
	})
}

func TestGetOrCreateTag(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTagReaderWriter := model.NewMockTagReaderWriter(ctrl)

	t.Run("tag already exists", func(t *testing.T) {
		inputTag := &model.Tag{
			Title: "Existing Tag",
		}
		existingTag := &model.Tag{
			Id:    123,
			Title: "Existing Tag",
		}

		mockTagReaderWriter.EXPECT().
			GetTag(gomock.Any(), inputTag).
			Return(existingTag, nil)

		result, err := GetOrCreateTag(context.Background(), mockTagReaderWriter, inputTag)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, existingTag, result)
	})

	t.Run("tag does not exist - create new", func(t *testing.T) {
		inputTag := &model.Tag{
			Title: "New Tag",
		}

		mockTagReaderWriter.EXPECT().
			GetTag(gomock.Any(), inputTag).
			Return(nil, gorm.ErrRecordNotFound)
		mockTagReaderWriter.EXPECT().
			CreateOrUpdateTag(gomock.Any(), inputTag).
			Return(nil)

		result, err := GetOrCreateTag(context.Background(), mockTagReaderWriter, inputTag)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, inputTag, result)
	})

	t.Run("error getting tag (not record not found)", func(t *testing.T) {
		inputTag := &model.Tag{
			Title: "Test Tag",
		}
		expectedError := errors.New("database connection error")

		mockTagReaderWriter.EXPECT().
			GetTag(gomock.Any(), inputTag).
			Return(nil, expectedError)

		result, err := GetOrCreateTag(context.Background(), mockTagReaderWriter, inputTag)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedError, err)
	})

	t.Run("error creating tag", func(t *testing.T) {
		inputTag := &model.Tag{
			Title: "New Tag",
		}
		expectedError := errors.New("create error")

		mockTagReaderWriter.EXPECT().
			GetTag(gomock.Any(), inputTag).
			Return(nil, gorm.ErrRecordNotFound)
		mockTagReaderWriter.EXPECT().
			CreateOrUpdateTag(gomock.Any(), inputTag).
			Return(expectedError)

		result, err := GetOrCreateTag(context.Background(), mockTagReaderWriter, inputTag)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedError, err)
	})
}

func TestGetOrCreateChildTag(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTagReaderWriter := model.NewMockTagReaderWriter(ctrl)

	t.Run("create child tag", func(t *testing.T) {
		parentId := uint64(456)
		title := "Child Tag"

		expectedTag := &model.Tag{
			ParentID: &parentId,
			Title:    title,
		}

		mockTagReaderWriter.EXPECT().
			GetTag(gomock.Any(), expectedTag).
			Return(nil, gorm.ErrRecordNotFound)
		mockTagReaderWriter.EXPECT().
			CreateOrUpdateTag(gomock.Any(), expectedTag).
			Return(nil)

		result, err := GetOrCreateChildTag(context.Background(), mockTagReaderWriter, parentId, title)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, &parentId, result.ParentID)
		assert.Equal(t, title, result.Title)
	})

	t.Run("child tag already exists", func(t *testing.T) {
		parentId := uint64(456)
		title := "Existing Child"

		queryTag := &model.Tag{
			ParentID: &parentId,
			Title:    title,
		}
		existingTag := &model.Tag{
			Id:       789,
			ParentID: &parentId,
			Title:    title,
		}

		mockTagReaderWriter.EXPECT().
			GetTag(gomock.Any(), queryTag).
			Return(existingTag, nil)

		result, err := GetOrCreateChildTag(context.Background(), mockTagReaderWriter, parentId, title)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, existingTag, result)
	})
}

func TestGetChildTag(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTagReader := model.NewMockTagReader(ctrl)

	t.Run("get existing child tag", func(t *testing.T) {
		parentId := uint64(456)
		title := "Child Tag"

		expectedTag := &model.Tag{
			ParentID: &parentId,
			Title:    title,
		}
		returnedTag := &model.Tag{
			Id:       789,
			ParentID: &parentId,
			Title:    title,
		}

		mockTagReader.EXPECT().
			GetTag(gomock.Any(), *expectedTag).
			Return(returnedTag, nil)

		result, err := GetChildTag(context.Background(), mockTagReader, parentId, title)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, returnedTag, result)
	})

	t.Run("child tag not found", func(t *testing.T) {
		parentId := uint64(456)
		title := "Non-existent Child"

		expectedTag := model.Tag{
			ParentID: &parentId,
			Title:    title,
		}

		mockTagReader.EXPECT().
			GetTag(gomock.Any(), expectedTag).
			Return(nil, gorm.ErrRecordNotFound)

		result, err := GetChildTag(context.Background(), mockTagReader, parentId, title)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})
}

func TestGetOrCreateTags(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTagReaderWriter := model.NewMockTagReaderWriter(ctrl)

	t.Run("get or create multiple tags", func(t *testing.T) {
		inputTags := []*model.Tag{
			{Title: "Tag 1"},
			{Title: "Tag 2"},
			{Title: "Tag 3"},
		}

		// First tag exists
		existingTag1 := &model.Tag{Id: 1, Title: "Tag 1"}
		mockTagReaderWriter.EXPECT().
			GetTag(gomock.Any(), inputTags[0]).
			Return(existingTag1, nil)

		// Second tag doesn't exist, needs to be created
		mockTagReaderWriter.EXPECT().
			GetTag(gomock.Any(), inputTags[1]).
			Return(nil, gorm.ErrRecordNotFound)
		mockTagReaderWriter.EXPECT().
			CreateOrUpdateTag(gomock.Any(), inputTags[1]).
			Return(nil)

		// Third tag exists
		existingTag3 := &model.Tag{Id: 3, Title: "Tag 3"}
		mockTagReaderWriter.EXPECT().
			GetTag(gomock.Any(), inputTags[2]).
			Return(existingTag3, nil)

		result, err := GetOrCreateTags(context.Background(), mockTagReaderWriter, inputTags)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 3)
		assert.Equal(t, existingTag1, result[0])
		assert.Equal(t, inputTags[1], result[1])
		assert.Equal(t, existingTag3, result[2])
	})

	t.Run("error creating one tag", func(t *testing.T) {
		inputTags := []*model.Tag{
			{Title: "Tag 1"},
			{Title: "Tag 2"},
		}

		// First tag succeeds
		existingTag1 := &model.Tag{Id: 1, Title: "Tag 1"}
		mockTagReaderWriter.EXPECT().
			GetTag(gomock.Any(), inputTags[0]).
			Return(existingTag1, nil)

		// Second tag fails to create
		expectedError := errors.New("create error")
		mockTagReaderWriter.EXPECT().
			GetTag(gomock.Any(), inputTags[1]).
			Return(nil, gorm.ErrRecordNotFound)
		mockTagReaderWriter.EXPECT().
			CreateOrUpdateTag(gomock.Any(), inputTags[1]).
			Return(expectedError)

		result, err := GetOrCreateTags(context.Background(), mockTagReaderWriter, inputTags)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedError, err)
	})

	t.Run("empty tags slice", func(t *testing.T) {
		inputTags := []*model.Tag{}

		result, err := GetOrCreateTags(context.Background(), mockTagReaderWriter, inputTags)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 0)
	})
}

func TestRemoveTagAndItsAssociations(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTagWriter := model.NewMockTagWriter(ctrl)

	t.Run("successful removal", func(t *testing.T) {
		tag := &model.Tag{
			Id:    123,
			Title: "Tag to Remove",
		}

		mockTagWriter.EXPECT().
			RemoveTag(gomock.Any(), tag.Id).
			Return(nil)

		errors := RemoveTagAndItsAssociations(context.Background(), mockTagWriter, tag)

		assert.Empty(t, errors)
	})

	t.Run("removal error", func(t *testing.T) {
		tag := &model.Tag{
			Id:    123,
			Title: "Tag to Remove",
		}

		expectedError := errors.New("removal error")
		mockTagWriter.EXPECT().
			RemoveTag(gomock.Any(), tag.Id).
			Return(expectedError)

		errorList := RemoveTagAndItsAssociations(context.Background(), mockTagWriter, tag)

		assert.Len(t, errorList, 1)
		assert.Equal(t, expectedError, errorList[0])
	})
}

func TestGetFullTags(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTagReader := model.NewMockTagReader(ctrl)

	t.Run("get full tags with items", func(t *testing.T) {
		tagIds := []*model.Tag{
			{Id: 1},
			{Id: 2},
			{Id: 3},
		}

		expectedIds := []uint64{1, 2, 3}
		expectedTags := &[]model.Tag{
			{Id: 1, Title: "Tag 1"},
			{Id: 2, Title: "Tag 2"},
			{Id: 3, Title: "Tag 3"},
		}

		mockTagReader.EXPECT().
			GetTags(gomock.Any(), expectedIds).
			Return(expectedTags, nil)

		result, err := GetFullTags(context.Background(), mockTagReader, tagIds)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedTags, result)
	})

	t.Run("empty tags slice", func(t *testing.T) {
		tagIds := []*model.Tag{}

		result, err := GetFullTags(context.Background(), mockTagReader, tagIds)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, *result, 0)
	})

	t.Run("error getting tags", func(t *testing.T) {
		tagIds := []*model.Tag{
			{Id: 1},
		}

		expectedIds := []uint64{1}
		expectedError := errors.New("database error")

		mockTagReader.EXPECT().
			GetTags(gomock.Any(), expectedIds).
			Return(nil, expectedError)

		result, err := GetFullTags(context.Background(), mockTagReader, tagIds)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedError, err)
	})
}

func TestGetCategories(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTagReader := model.NewMockTagReader(ctrl)

	t.Run("get categories successfully", func(t *testing.T) {
		expectedTags := &[]model.Tag{
			{Id: 1, Title: "Category 1", ParentID: nil},
			{Id: 2, Title: "Category 2", ParentID: nil},
		}

		mockTagReader.EXPECT().
			GetTags(gomock.Any(), "parent_id is NULL").
			Return(expectedTags, nil)

		result, err := GetCategories(context.Background(), mockTagReader)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedTags, result)
	})

	t.Run("error getting categories", func(t *testing.T) {
		expectedError := errors.New("database error")

		mockTagReader.EXPECT().
			GetTags(gomock.Any(), "parent_id is NULL").
			Return(nil, expectedError)

		result, err := GetCategories(context.Background(), mockTagReader)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedError, err)
	})
}

func TestIsBelongToCategory(t *testing.T) {
	t.Run("tag belongs to category", func(t *testing.T) {
		categoryId := uint64(100)
		tag := &model.Tag{
			Id:       1,
			Title:    "Child Tag",
			ParentID: &categoryId,
		}
		category := &model.Tag{
			Id:    categoryId,
			Title: "Parent Category",
		}

		result := IsBelongToCategory(tag, category)
		assert.True(t, result)
	})

	t.Run("tag does not belong to category", func(t *testing.T) {
		differentCategoryId := uint64(200)
		tag := &model.Tag{
			Id:       1,
			Title:    "Child Tag",
			ParentID: &differentCategoryId,
		}
		category := &model.Tag{
			Id:    100,
			Title: "Parent Category",
		}

		result := IsBelongToCategory(tag, category)
		assert.False(t, result)
	})

	t.Run("tag has no parent (top-level tag)", func(t *testing.T) {
		tag := &model.Tag{
			Id:       1,
			Title:    "Top Level Tag",
			ParentID: nil,
		}
		category := &model.Tag{
			Id:    100,
			Title: "Parent Category",
		}

		result := IsBelongToCategory(tag, category)
		assert.False(t, result)
	})
}

func TestRemoveTagImages(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTagReaderWriter := model.NewMockTagReaderWriter(ctrl)

	t.Run("remove tag images successfully", func(t *testing.T) {
		tagId := uint64(123)
		titId := uint64(456)

		tag := &model.Tag{
			Id:    tagId,
			Title: "Test Tag",
			Images: []*model.TagImage{
				{Id: 1, ImageTypeId: 456, Url: "image1.jpg"}, // Should be removed
				{Id: 2, ImageTypeId: 789, Url: "image2.jpg"}, // Different type, should not be removed
				{Id: 3, ImageTypeId: 456, Url: "image3.jpg"}, // Should be removed
			},
		}

		mockTagReaderWriter.EXPECT().
			GetTag(gomock.Any(), tagId).
			Return(tag, nil)
		mockTagReaderWriter.EXPECT().
			RemoveTagImageFromTag(gomock.Any(), tagId, uint64(1)).
			Return(nil)
		mockTagReaderWriter.EXPECT().
			RemoveTagImageFromTag(gomock.Any(), tagId, uint64(3)).
			Return(nil)

		err := RemoveTagImages(context.Background(), mockTagReaderWriter, tagId, titId)

		assert.NoError(t, err)
	})

	t.Run("error getting tag", func(t *testing.T) {
		tagId := uint64(123)
		titId := uint64(456)
		expectedError := errors.New("tag not found")

		mockTagReaderWriter.EXPECT().
			GetTag(gomock.Any(), tagId).
			Return(nil, expectedError)

		err := RemoveTagImages(context.Background(), mockTagReaderWriter, tagId, titId)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})

	t.Run("error removing tag image", func(t *testing.T) {
		tagId := uint64(123)
		titId := uint64(456)
		expectedError := errors.New("remove error")

		tag := &model.Tag{
			Id:    tagId,
			Title: "Test Tag",
			Images: []*model.TagImage{
				{Id: 1, ImageTypeId: 456, Url: "image1.jpg"},
			},
		}

		mockTagReaderWriter.EXPECT().
			GetTag(gomock.Any(), tagId).
			Return(tag, nil)
		mockTagReaderWriter.EXPECT().
			RemoveTagImageFromTag(gomock.Any(), tagId, uint64(1)).
			Return(expectedError)

		err := RemoveTagImages(context.Background(), mockTagReaderWriter, tagId, titId)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})

	t.Run("no matching images to remove", func(t *testing.T) {
		tagId := uint64(123)
		titId := uint64(456)

		tag := &model.Tag{
			Id:    tagId,
			Title: "Test Tag",
			Images: []*model.TagImage{
				{Id: 1, ImageTypeId: 789, Url: "image1.jpg"}, // Different type
				{Id: 2, ImageTypeId: 999, Url: "image2.jpg"}, // Different type
			},
		}

		mockTagReaderWriter.EXPECT().
			GetTag(gomock.Any(), tagId).
			Return(tag, nil)

		err := RemoveTagImages(context.Background(), mockTagReaderWriter, tagId, titId)

		assert.NoError(t, err)
	})

	t.Run("tag with no images", func(t *testing.T) {
		tagId := uint64(123)
		titId := uint64(456)

		tag := &model.Tag{
			Id:     tagId,
			Title:  "Test Tag",
			Images: []*model.TagImage{},
		}

		mockTagReaderWriter.EXPECT().
			GetTag(gomock.Any(), tagId).
			Return(tag, nil)

		err := RemoveTagImages(context.Background(), mockTagReaderWriter, tagId, titId)

		assert.NoError(t, err)
	})
}
