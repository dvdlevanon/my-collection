package db

import (
	"errors"
	"fmt"
	"my-collection/server/pkg/model"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"k8s.io/utils/pointer"
)

func setupNewDb(t *testing.T, filename string) (*Database, error) {
	assert.NoError(t, os.MkdirAll(".tests", 0750))
	dbpath := fmt.Sprintf(".tests/%s", filename)
	_, err := os.Create(dbpath)
	assert.NoError(t, err)
	assert.NoError(t, os.Remove(dbpath))
	return New("", dbpath)
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

	directory := model.Directory{
		Path:       "path/to/file",
		FilesCount: 3,
		LastSynced: 1234567,
		Tags: []*model.Tag{
			{
				Title: "tag1",
			},
			{
				Title: "tag2",
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
	assert.Equal(t, directory.Tags[0].Id, directoryFromDB.Tags[0].Id)
	assert.Equal(t, directory.Tags[1].Id, directoryFromDB.Tags[1].Id)
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
	assert.Equal(t, directory.Tags[0].Id, directoryFromDB.Tags[0].Id)
	assert.Equal(t, directory.Tags[1].Id, directoryFromDB.Tags[1].Id)
	assert.Empty(t, directoryFromDB.Tags[0].Title)
	assert.Empty(t, directoryFromDB.Tags[1].Title)
	assert.Equal(t, excludedDirectory.Path, excludedDirectoryFromDB.Path)
	assert.NotNil(t, excludedDirectoryFromDB.Excluded)
	assert.True(t, *excludedDirectoryFromDB.Excluded)

	directories, err := db.GetAllDirectories()
	assert.NoError(t, err)
	assert.Len(t, *directories, 2)
}
