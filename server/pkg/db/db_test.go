package db

import (
	"errors"
	"fmt"
	"my-collection/server/pkg/model"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"k8s.io/utils/pointer"
)

func setupNewDb(t *testing.T, filename string) (Database, error) {
	assert.NoError(t, os.MkdirAll(".tests", 0750))
	dbpath := fmt.Sprintf(".tests/%s", filename)
	_, err := os.Create(dbpath)
	assert.NoError(t, err)
	assert.NoError(t, os.Remove(dbpath))
	return New(filepath.Join("", dbpath), false)
}

func TestItem(t *testing.T) {
	db, err := setupNewDb(t, "test-item.sqlite")
	assert.NoError(t, err)

	for j := 0; j < 2; j++ {
		for i := 1; i < 5; i++ {
			expectedId := uint64((j * 4) + i)
			origin := fmt.Sprintf("origin-%d", j)
			title := fmt.Sprintf("title-%d", i)
			emptyItem := model.Item{}
			assert.Error(t, db.CreateOrUpdateItem(&emptyItem))
			onlyTitleItem := model.Item{Title: title}
			assert.Error(t, db.CreateOrUpdateItem(&onlyTitleItem))
			onlyOriginItem := model.Item{Origin: origin}
			assert.Error(t, db.CreateOrUpdateItem(&onlyOriginItem))
			validItem := model.Item{Title: title, Origin: origin}
			assert.NoError(t, db.CreateOrUpdateItem(&validItem))
			itemFromDb, err := db.GetItem(expectedId)
			assert.NoError(t, err)
			assert.Equal(t, expectedId, itemFromDb.Id)
			assert.Equal(t, origin, itemFromDb.Origin)
			assert.Equal(t, title, itemFromDb.Title)
		}
	}

	assert.NoError(t, db.CreateOrUpdateTag(&model.Tag{Title: "Tag1"}))
	assert.NoError(t, db.CreateOrUpdateTag(&model.Tag{Title: "Tag2"}))

	for j := 0; j < 2; j++ {
		for i := 1; i < 5; i++ {
			expectedId := uint64((j * 4) + i)
			origin := fmt.Sprintf("origin-%d", j)
			title := fmt.Sprintf("title-%d", i)
			onlyTitleItem := model.Item{Title: title}
			assert.Error(t, db.CreateOrUpdateItem(&onlyTitleItem))
			onlyOriginItem := model.Item{Origin: origin}
			assert.Error(t, db.CreateOrUpdateItem(&onlyOriginItem))
			validItem := model.Item{Title: title, Origin: origin, Url: fmt.Sprintf("url-%d", i)}
			assert.NoError(t, db.CreateOrUpdateItem(&validItem))
			itemFromDb, err := db.GetItem(model.Item{Title: title, Origin: origin})
			assert.NoError(t, err)
			assert.Equal(t, expectedId, itemFromDb.Id)
			assert.Equal(t, origin, itemFromDb.Origin)
			assert.Equal(t, title, itemFromDb.Title)
			assert.Equal(t, fmt.Sprintf("url-%d", i), itemFromDb.Url)
			onlyIdItem := model.Item{Id: expectedId, PreviewUrl: fmt.Sprintf("preview-url-%d", i)}
			assert.NoError(t, db.CreateOrUpdateItem(&onlyIdItem))
			itemFromDb2, err := db.GetItem(model.Item{Id: expectedId})
			assert.NoError(t, err)
			assert.Equal(t, expectedId, itemFromDb.Id)
			assert.Equal(t, fmt.Sprintf("preview-url-%d", i), itemFromDb2.PreviewUrl)

			assert.NoError(t, db.CreateOrUpdateItem(&model.Item{Id: expectedId, Tags: []*model.Tag{{Id: 1}, {Id: 2}}}))
			// todo - add more checks
		}
	}
}

func TestTag(t *testing.T) {
	db, err := setupNewDb(t, "test-tag.sqlite")
	assert.NoError(t, err)

	assert.NoError(t, db.CreateOrUpdateItem(&model.Item{Title: "Item1", Origin: "origin"}))
	assert.NoError(t, db.CreateOrUpdateItem(&model.Item{Title: "Item2", Origin: "origin"}))

	for i := 1; i < 5; i++ {
		expectedId := uint64(i)
		title := fmt.Sprintf("parent-%d", i)
		emptyTag := model.Tag{}
		assert.Error(t, db.CreateOrUpdateTag(&emptyTag))
		parentTag := model.Tag{Title: title}
		assert.NoError(t, db.CreateOrUpdateTag(&parentTag))
		parentFromDb, err := db.GetTag(expectedId)
		assert.NoError(t, err)
		assert.Equal(t, title, parentFromDb.Title)
		updatedFromDb, err := db.GetTag(model.Tag{Id: expectedId})
		assert.NoError(t, err)
		assert.Equal(t, expectedId, updatedFromDb.Id)
		assert.Equal(t, title, updatedFromDb.Title)
		assert.NoError(t, db.CreateOrUpdateTag(&model.Tag{Id: expectedId, Items: []*model.Item{{Id: 1}, {Id: 2}}}))
		withItemsFromDb, err := db.GetTag("title = ?", title)
		assert.NoError(t, err)
		assert.Equal(t, expectedId, withItemsFromDb.Id)
		assert.Equal(t, 2, len(withItemsFromDb.Items))
		assert.Equal(t, uint64(1), withItemsFromDb.Items[0].Id)
		assert.Equal(t, uint64(2), withItemsFromDb.Items[1].Id)

		itemIds := make([]uint64, 0)
		for _, item := range withItemsFromDb.Items {
			itemIds = append(itemIds, item.Id)
		}

		items, err := db.GetItems(itemIds)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(*items))
		assert.Equal(t, uint64(1), (*items)[0].Id)
		assert.Equal(t, uint64(2), (*items)[1].Id)
	}

	tags, err := db.GetTags("parent_id is NULL")
	assert.NoError(t, err)
	assert.Equal(t, 4, len(*tags))

	for j, parent := range *tags {
		for i := 1; i < 3; i++ {
			expectedId := uint64(4 + ((j * 2) + i))
			title := fmt.Sprintf("child-%d", i)
			childTag := model.Tag{Title: title, ParentID: &parent.Id}
			assert.NoError(t, db.CreateOrUpdateTag(&childTag))
			childFromDb, err := db.GetTag("title = ? and parent_id = ?", title, parent.Id)
			assert.NoError(t, err)
			assert.Equal(t, expectedId, childFromDb.Id)
		}
	}
}

func TestManyToMany(t *testing.T) {
	db, err := setupNewDb(t, "many-to-many.sqlite")
	assert.NoError(t, err)
	item := &model.Item{Title: "item1", Origin: "origin"}
	tag := &model.Tag{Title: "tag1"}
	assert.NoError(t, db.CreateOrUpdateItem(item))
	assert.NoError(t, db.CreateOrUpdateTag(tag))
	item.Tags = append(item.Tags, tag)
	tag.Items = append(tag.Items, item)
	assert.NoError(t, db.UpdateItem(item))
	assert.NoError(t, db.UpdateTag(tag))
	itemFromDB, err := db.GetItem(1)
	assert.NoError(t, err)
	tagFromDB, err := db.GetTag(1)
	assert.NoError(t, err)
	assert.Equal(t, itemFromDB.Title, item.Title)
	assert.Equal(t, tagFromDB.Title, tag.Title)
	assert.Equal(t, len(itemFromDB.Tags), len(item.Tags))
	assert.Equal(t, len(tagFromDB.Items), len(tag.Items))
	assert.Equal(t, itemFromDB.Tags[0].Id, item.Tags[0].Id)
	assert.Equal(t, itemFromDB.Tags[0].Title, item.Tags[0].Title)
	assert.Equal(t, tagFromDB.Items[0].Id, tag.Items[0].Id)
	assert.Empty(t, tagFromDB.Items[0].Title)
}

func TestOneToManyParent(t *testing.T) {
	db, err := setupNewDb(t, "one-to-many.sqlite")
	assert.NoError(t, err)
	parent := &model.Tag{Title: "parent"}
	child1 := &model.Tag{Title: "child1"}
	child2 := &model.Tag{Title: "child2"}
	assert.NoError(t, db.CreateOrUpdateTag(parent))
	assert.NoError(t, db.CreateOrUpdateTag(child1))
	assert.NoError(t, db.CreateOrUpdateTag(child2))
	parent.Children = append(parent.Children, child1, child2)
	assert.NoError(t, db.UpdateTag(parent))
	parentFromDB, err := db.GetTag(1)
	assert.NoError(t, err)
	child1FromDB, err := db.GetTag(2)
	assert.NoError(t, err)
	child2FromDB, err := db.GetTag(3)
	assert.NoError(t, err)
	assert.Equal(t, parentFromDB.Title, parent.Title)
	assert.Equal(t, child1FromDB.Title, child1.Title)
	assert.Equal(t, child2FromDB.Title, child2.Title)
	assert.Equal(t, len(parentFromDB.Children), len(parent.Children))
	assert.Equal(t, parentFromDB.Children[0].Id, parent.Children[0].Id)
	assert.Equal(t, parentFromDB.Children[1].Id, parent.Children[1].Id)
	assert.Equal(t, *child1FromDB.ParentID, parent.Id)
	assert.Equal(t, *child2FromDB.ParentID, parent.Id)
}

func TestOneToMany(t *testing.T) {
	db, err := setupNewDb(t, "one-to-many-test.sqlite")
	assert.NoError(t, err)
	item1 := &model.Item{Title: "title1", Origin: "origin"}
	assert.NoError(t, db.CreateOrUpdateItem(item1))
	itemFromDB, err := db.GetItem(1)
	assert.NoError(t, err)
	itemFromDB.Covers = []model.Cover{
		{
			Url: "cover1",
		},
		{
			Url: "cover2",
		},
	}
	err = db.UpdateItem(itemFromDB)
	assert.NoError(t, err)
	updatedFromDB, err := db.GetItem(1)
	assert.NoError(t, err)
	assert.Equal(t, len(itemFromDB.Covers), len(updatedFromDB.Covers))
	assert.Equal(t, updatedFromDB.Covers[0].ItemId, uint64(1))
	assert.Equal(t, updatedFromDB.Covers[1].ItemId, uint64(1))
	assert.Equal(t, updatedFromDB.Covers[0].Id, uint64(1))
	assert.Equal(t, updatedFromDB.Covers[1].Id, uint64(2))
	assert.Equal(t, updatedFromDB.Covers[0].Url, itemFromDB.Covers[0].Url)
	assert.Equal(t, updatedFromDB.Covers[1].Url, itemFromDB.Covers[1].Url)
}

func TestGetMissingItem(t *testing.T) {
	db, err := setupNewDb(t, "missing-item-test.sqlite")
	assert.NoError(t, err)
	_, err = db.GetItem(666)
	assert.Error(t, err)
}

func TestGetMissingTag(t *testing.T) {
	db, err := setupNewDb(t, "missing-tag-test.sqlite")
	assert.NoError(t, err)
	_, err = db.GetTag(666)
	assert.Error(t, err)
}

func TestRemoveTag(t *testing.T) {
	db, err := setupNewDb(t, "remove-tag-test.sqlite")
	assert.NoError(t, err)
	item := &model.Item{Title: "item1", Origin: "origin"}
	tag := &model.Tag{Title: "tag1"}
	assert.NoError(t, db.CreateOrUpdateItem(item))
	assert.NoError(t, db.CreateOrUpdateTag(tag))
	item.Tags = append(item.Tags, tag)
	tag.Items = append(tag.Items, item)
	assert.NoError(t, db.UpdateItem(item))
	assert.NoError(t, db.UpdateTag(tag))
	itemFromDB, err := db.GetItem(1)
	assert.NoError(t, err)
	assert.Equal(t, len(itemFromDB.Tags), len(item.Tags))
	assert.NoError(t, db.RemoveTagFromItem(item.Id, item.Tags[0].Id))
	itemFromDBAfterRemoval, err := db.GetItem(1)
	assert.NoError(t, err)
	assert.NotEqual(t, len(itemFromDB.Tags), len(itemFromDBAfterRemoval.Tags))
	assert.NotEqual(t, uint64(0), len(itemFromDBAfterRemoval.Tags))
	tagFromDB, err := db.GetTag(1)
	assert.NoError(t, err)
	assert.Equal(t, tagFromDB.Title, tag.Title)
}

func TestTagAnnotations(t *testing.T) {
	db, err := setupNewDb(t, "tag-annotations-test.sqlite")
	assert.NoError(t, err)

	child1 := model.Tag{
		Title: "child1",
	}

	child2 := model.Tag{
		Title: "child2",
	}

	assert.NoError(t, db.CreateOrUpdateTag(&child1))
	assert.NoError(t, db.CreateOrUpdateTag(&child2))

	root := model.Tag{
		Title:    "root",
		Children: []*model.Tag{&child1, &child2},
	}

	assert.NoError(t, db.CreateTagAnnotation(&model.TagAnnotation{Title: "annotation1"}))
	annotation1, err := db.GetTagAnnotation(1)
	assert.NoError(t, err)
	assert.NoError(t, db.CreateTagAnnotation(&model.TagAnnotation{Title: "annotation2"}))
	annotation2, err := db.GetTagAnnotation(2)
	assert.NoError(t, err)
	assert.NoError(t, db.CreateTagAnnotation(&model.TagAnnotation{Title: "annotation3"}))
	annotation3, err := db.GetTagAnnotation(3)
	assert.NoError(t, err)

	assert.NoError(t, db.CreateOrUpdateTag(&root))
	child1.Annotations = append(child1.Annotations, annotation1, annotation2)
	assert.NoError(t, db.CreateOrUpdateTag(&child1))
	child2.Annotations = append(child2.Annotations, annotation1, annotation3)
	assert.NoError(t, db.CreateOrUpdateTag(&child2))

	rootFromDB, err := db.GetTag(root.Id)
	assert.NoError(t, err)
	assert.Equal(t, root.Title, rootFromDB.Title)
	assert.Equal(t, 0, len(rootFromDB.Annotations))
	child1FromDB, err := db.GetTag(child1.Id)
	assert.NoError(t, err)
	assert.Equal(t, child1.Title, child1FromDB.Title)
	assert.Equal(t, 2, len(child1FromDB.Annotations))
	assert.Empty(t, child1FromDB.Annotations[0].Title)
	assert.Empty(t, child1FromDB.Annotations[1].Title)
	child2FromDB, err := db.GetTag(child2.Id)
	assert.NoError(t, err)
	assert.Equal(t, child2.Title, child2FromDB.Title)
	assert.Equal(t, 2, len(child2FromDB.Annotations))
	assert.Empty(t, child2FromDB.Annotations[0].Title)
	assert.Empty(t, child2FromDB.Annotations[1].Title)

	assert.Equal(t, annotation1.Id, child1.Annotations[0].Id)
	assert.Equal(t, annotation2.Id, child1.Annotations[1].Id)
	assert.Equal(t, annotation1.Id, child2.Annotations[0].Id)
	assert.Equal(t, annotation3.Id, child2.Annotations[1].Id)
}

func TestGetTagAnnotation(t *testing.T) {
	db, err := setupNewDb(t, "tag-get-tag-annotation.sqlite")
	assert.NoError(t, err)
	assert.NoError(t, db.CreateTagAnnotation(&model.TagAnnotation{Title: "annotation"}))

	annotation, err := db.GetTagAnnotation(&model.TagAnnotation{Title: "annotation"})
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), annotation.Id)
	assert.Equal(t, "annotation", annotation.Title)

	annotationUsingId, err := db.GetTagAnnotation(1)
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), annotationUsingId.Id)
	assert.Equal(t, "annotation", annotationUsingId.Title)

	_, err = db.GetTagAnnotation(&model.TagAnnotation{Title: "not-exists"})
	assert.Error(t, err)
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
}

func TestRemoveTagAnnotation(t *testing.T) {
	db, err := setupNewDb(t, "remove-tag-annotation.sqlite")
	assert.NoError(t, err)

	assert.NoError(t, db.CreateOrUpdateTag(&model.Tag{
		Title: "tag1",
		Annotations: []*model.TagAnnotation{
			{
				Title: "annotation1",
			},
		},
	}))

	tag, err := db.GetTag(1)
	assert.NoError(t, err)
	assert.Equal(t, "tag1", tag.Title)
	assert.Equal(t, 1, len(tag.Annotations))
	assert.Equal(t, uint64(1), tag.Annotations[0].Id)
	assert.NoError(t, db.RemoveTagAnnotationFromTag(1, 1))
	tagAfterRemove, err := db.GetTag(1)
	assert.NoError(t, err)
	assert.Equal(t, "tag1", tagAfterRemove.Title)
	assert.Equal(t, 0, len(tagAfterRemove.Annotations))
}

func TestDirectories(t *testing.T) {
	db, err := setupNewDb(t, "create-directories.sqlite")
	assert.NoError(t, err)

	assert.NoError(t, db.CreateOrUpdateTag(&model.Tag{Title: "tag1"}))
	assert.NoError(t, db.CreateOrUpdateTag(&model.Tag{Title: "tag2"}))

	directory := model.Directory{
		Path:       "path/to/file",
		FilesCount: pointer.Int(3),
		LastSynced: 1234567,
		Tags: []*model.Tag{
			{
				ParentID: pointer.Uint64(0),
				Id:       2,
			},
			{
				ParentID: pointer.Uint64(0),
				Id:       1,
			},
		},
	}

	assert.NoError(t, db.CreateOrUpdateDirectory(&directory))
	directoryFromDB, err := db.GetDirectory("path = ?", directory.Path)
	assert.NoError(t, err)
	assert.Equal(t, directory.Path, directoryFromDB.Path)
	assert.Equal(t, directory.FilesCount, directoryFromDB.FilesCount)
	assert.Equal(t, directory.LastSynced, directoryFromDB.LastSynced)
	assert.Nil(t, directoryFromDB.Excluded)
	assert.NotNil(t, directoryFromDB.Tags)
	assert.Equal(t, len(directory.Tags), len(directoryFromDB.Tags))
	assert.Equal(t, directory.Tags[0].Id, directoryFromDB.Tags[1].Id)
	assert.Equal(t, directory.Tags[1].Id, directoryFromDB.Tags[0].Id)
	assert.Empty(t, directoryFromDB.Tags[0].Title)
	assert.Empty(t, directoryFromDB.Tags[1].Title)

	excludedDirectory := model.Directory{
		Path:     "path/to/excluded",
		Excluded: pointer.Bool(true),
	}

	duplicatedDirectory := model.Directory{
		Path: "path/to/file",
	}

	assert.NoError(t, db.CreateOrUpdateDirectory(&excludedDirectory))
	assert.NoError(t, db.CreateOrUpdateDirectory(&duplicatedDirectory))
	excludedDirectoryFromDB, err := db.GetDirectory("path = ?", excludedDirectory.Path)
	assert.NoError(t, err)
	duplicatedDirectoryFromDB, err := db.GetDirectory("path = ?", duplicatedDirectory.Path)
	assert.NoError(t, err)
	assert.Equal(t, directory.Path, duplicatedDirectoryFromDB.Path)
	assert.Equal(t, directory.FilesCount, duplicatedDirectoryFromDB.FilesCount)
	assert.Equal(t, directory.LastSynced, duplicatedDirectoryFromDB.LastSynced)
	assert.Nil(t, directoryFromDB.Excluded)
	assert.NotNil(t, directoryFromDB.Tags)
	assert.Equal(t, len(directory.Tags), len(directoryFromDB.Tags))
	assert.Equal(t, directory.Tags[0].Id, directoryFromDB.Tags[1].Id)
	assert.Equal(t, directory.Tags[1].Id, directoryFromDB.Tags[0].Id)
	assert.Empty(t, directoryFromDB.Tags[0].Title)
	assert.Empty(t, directoryFromDB.Tags[1].Title)
	assert.Equal(t, excludedDirectory.Path, excludedDirectoryFromDB.Path)
	assert.NotNil(t, excludedDirectoryFromDB.Excluded)
	assert.True(t, *excludedDirectoryFromDB.Excluded)

	directories, err := db.GetAllDirectories()
	assert.NoError(t, err)
	assert.Len(t, *directories, 2)
}

func TestRemoveItem(t *testing.T) {
	db, err := setupNewDb(t, "remove-item.sqlite")
	assert.NoError(t, err)
	item := &model.Item{Title: "item1", Origin: "origin"}
	tag := &model.Tag{Title: "tag1"}
	assert.NoError(t, db.CreateOrUpdateTag(tag))
	item.Tags = append(item.Tags, tag)
	assert.NoError(t, db.CreateOrUpdateItem(item))
	assert.NoError(t, db.RemoveItem(item.Id))
	newItem := &model.Item{Title: "new-item", Origin: "origin"}
	assert.NoError(t, db.CreateOrUpdateItem(newItem))
	loaded, err := db.GetItem(newItem.Id)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(loaded.Tags))
}

func TestSubItems(t *testing.T) {
	db, err := setupNewDb(t, "sub-items.sqlite")
	assert.NoError(t, err)
	item := &model.Item{Title: "main", Origin: "origin", DurationSeconds: 100}
	sub1 := &model.Item{Title: "sub1", Origin: "origin", StartPosition: 0, EndPosition: 50}
	sub2 := &model.Item{Title: "sub2", Origin: "origin", StartPosition: 50, EndPosition: 100}
	assert.NoError(t, db.CreateOrUpdateItem(item))
	assert.NoError(t, db.CreateOrUpdateItem(sub1))
	assert.NoError(t, db.CreateOrUpdateItem(sub2))
	item.SubItems = append(item.SubItems, sub1, sub2)
	assert.NoError(t, db.UpdateItem(item))
	sub1.Covers = append(sub1.Covers, model.Cover{Url: "test"})
	sub2.Covers = append(sub1.Covers, model.Cover{Url: "test"})
	assert.NoError(t, db.UpdateItem(sub1))
	assert.NoError(t, db.UpdateItem(sub2))
	mainItem, err := db.GetItem(item.Id)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(mainItem.SubItems))
	assert.Equal(t, mainItem.Id, *sub1.MainItemId)
	assert.Equal(t, mainItem.Id, *sub2.MainItemId)
	assert.Equal(t, 1, len(mainItem.SubItems[0].Covers))
	assert.Equal(t, "test", mainItem.SubItems[0].Covers[0].Url)
}

func TestHighlights(t *testing.T) {
	db, err := setupNewDb(t, "highlights.sqlite")
	assert.NoError(t, err)
	item := &model.Item{Title: "main", Origin: "origin", DurationSeconds: 100}
	hl1 := &model.Item{Title: "hl1", Origin: "origin", StartPosition: 0, EndPosition: 50}
	hl2 := &model.Item{Title: "hl2", Origin: "origin", StartPosition: 50, EndPosition: 100}
	assert.NoError(t, db.CreateOrUpdateItem(item))
	assert.NoError(t, db.CreateOrUpdateItem(hl1))
	assert.NoError(t, db.CreateOrUpdateItem(hl2))
	item.Highlights = append(item.Highlights, hl1, hl2)
	assert.NoError(t, db.UpdateItem(item))
	hl1.Covers = append(hl1.Covers, model.Cover{Url: "test"})
	hl2.Covers = append(hl1.Covers, model.Cover{Url: "test"})
	assert.NoError(t, db.UpdateItem(hl1))
	assert.NoError(t, db.UpdateItem(hl2))
	mainItem, err := db.GetItem(item.Id)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(mainItem.Highlights))
	assert.Equal(t, mainItem.Id, *hl1.HighlightParentItemId)
	assert.Equal(t, mainItem.Id, *hl2.HighlightParentItemId)
	assert.Equal(t, 1, len(mainItem.Highlights[0].Covers))
	assert.Equal(t, "test", mainItem.Highlights[0].Covers[0].Url)
}

func TestTagImageThumbnail(t *testing.T) {
	db, err := setupNewDb(t, "tag-image-thumbnail.sqlite")
	assert.NoError(t, err)

	tag := model.Tag{
		Title: "test",
		Images: []*model.TagImage{{
			Url: "some/url",
			ThumbnailRect: model.Rect{
				X: 100,
				Y: 101,
				W: 102,
				H: 103,
			},
		}},
	}

	assert.NoError(t, db.CreateOrUpdateTag(&tag))
	fromDb, err := db.GetTag("title = ?", tag.Title)
	assert.NoError(t, err)
	assert.Equal(t, tag.Title, fromDb.Title)
	assert.Equal(t, len(tag.Images), len(fromDb.Images))
	assert.Equal(t, 1, len(fromDb.Images))
	assert.Equal(t, tag.Images[0].Url, fromDb.Images[0].Url)
	assert.Equal(t, 100, fromDb.Images[0].ThumbnailRect.X)
	assert.Equal(t, 101, fromDb.Images[0].ThumbnailRect.Y)
	assert.Equal(t, 102, fromDb.Images[0].ThumbnailRect.W)
	assert.Equal(t, 103, fromDb.Images[0].ThumbnailRect.H)

	fromDb.Images[0].ThumbnailRect.H = 104
	assert.NoError(t, db.UpdateTagImage(fromDb.Images[0]))

	fromDb2, err := db.GetTag("title = ?", tag.Title)
	assert.NoError(t, err)
	assert.Equal(t, 104, fromDb2.Images[0].ThumbnailRect.H)
}

func TestTagCustomCommands(t *testing.T) {
	db, err := setupNewDb(t, "tags-custom-commands.sqlite")
	assert.NoError(t, err)

	command1 := model.TagCustomCommand{
		Title:   "command1",
		Type:    "command1-type",
		Arg:     "command1-arg",
		Tooltip: "command1-tooltip",
		Icon:    "command1-icon",
		TagId:   1234,
	}

	command2 := model.TagCustomCommand{
		Title:   "command2",
		Type:    "command2-type",
		Arg:     "command2-arg",
		Tooltip: "command2-tooltip",
		Icon:    "command2-icon",
		TagId:   2345,
	}

	assert.NoError(t, db.CreateOrUpdateTagCustomCommand(&command1))
	assert.NoError(t, db.CreateOrUpdateTagCustomCommand(&command2))

	commands, err := db.GetAllTagCustomCommands()
	assert.NoError(t, err)
	assert.Equal(t, 2, len(*commands))
}

// Fail because of https://github.com/go-gorm/sqlite/issues/134
//
// func TestReuseId(t *testing.T) {
// 	db, err := setupNewDb(t, "test-reuse-id.sqlite")
// 	assert.NoError(t, err)

// 	firstItem := test{Somethingelse: "origin"}
// 	assert.NoError(t, db.create(&firstItem))
// 	assert.Equal(t, uint64(1), firstItem.Dude)
// 	assert.NoError(t, db.delete(firstItem))
// 	newItem := test{Somethingelse: "origin"}
// 	assert.NoError(t, db.create(&newItem))
// 	assert.NotEqual(t, uint64(1), newItem.Dude)
// }
