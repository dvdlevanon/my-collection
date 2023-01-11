package db

import (
	"fmt"
	"my-collection/server/pkg/model"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupNewDb(t *testing.T, filename string) (*Database, error) {
	assert.NoError(t, os.MkdirAll(".tests", 0750))
	dbpath := fmt.Sprintf(".tests/%s", filename)
	_, err := os.Create(dbpath)
	assert.NoError(t, err)
	assert.NoError(t, os.Remove(dbpath))
	return New("", dbpath)
}

func TestCreate(t *testing.T) {
	db, err := setupNewDb(t, "create-test.sqlite")
	assert.NoError(t, err)
	item1 := &model.Item{Title: "title1"}
	assert.Equal(t, item1.Id, uint64(0))
	assert.NoError(t, db.CreateItem(item1))
	assert.Equal(t, item1.Id, uint64(1))
	itemFromDB, err := db.GetItem(1)
	assert.NoError(t, err)
	assert.Equal(t, itemFromDB.Title, item1.Title)
	item1.Title = "update-title"
	assert.Error(t, db.CreateItem(item1))
}

func TestGetBy(t *testing.T) {
	db, err := setupNewDb(t, "create-test-get-by.sqlite")
	assert.NoError(t, err)
	item1 := &model.Item{Title: "title1"}
	assert.Equal(t, item1.Id, uint64(0))
	assert.NoError(t, db.CreateItem(item1))
	assert.Equal(t, item1.Id, uint64(1))
	itemFromDB, err := db.GetItem("title = ?", item1.Title)
	assert.NoError(t, err)
	assert.Equal(t, itemFromDB.Title, item1.Title)
	item1.Title = "updated-title"
	assert.Error(t, db.CreateItem(item1))
}

func TestCreateOrUpdate(t *testing.T) {
	db, err := setupNewDb(t, "create-or-update-test.sqlite")
	assert.NoError(t, err)
	item1 := &model.Item{Title: "title1"}
	assert.Equal(t, item1.Id, uint64(0))
	assert.NoError(t, db.CreateItem(item1))
	assert.Equal(t, item1.Id, uint64(1))
	itemFromDB, err := db.GetItem(1)
	assert.NoError(t, err)
	assert.Equal(t, itemFromDB.Title, item1.Title)
	item1.Url = "updated-url"
	item1.Id = 0
	assert.NoError(t, db.CreateOrUpdateItem(item1))
	item1.Id = 0
	item1.Title = "title5"
	assert.NoError(t, db.CreateOrUpdateItem(item1))
	assert.Equal(t, item1.Id, uint64(2))
	item5FromDB, err := db.GetItem(2)
	assert.NoError(t, err)
	assert.Equal(t, item5FromDB.Title, item1.Title)
}

func TestUpdate(t *testing.T) {
	db, err := setupNewDb(t, "update-test.sqlite")
	assert.NoError(t, err)
	item1 := &model.Item{Title: "title1"}
	assert.NoError(t, db.CreateItem(item1))
	item1.Title = "update-title"
	assert.NoError(t, db.UpdateItem(item1))
	itemFromDB, err := db.GetItem(1)
	assert.NoError(t, err)
	assert.Equal(t, itemFromDB.Title, item1.Title)
}

func TestManyToMany(t *testing.T) {
	db, err := setupNewDb(t, "many-to-many.sqlite")
	assert.NoError(t, err)
	item := &model.Item{Title: "item1"}
	tag := &model.Tag{Title: "tag1"}
	assert.NoError(t, db.CreateItem(item))
	assert.NoError(t, db.CreateTag(tag))
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
	assert.NoError(t, db.CreateTag(parent))
	assert.NoError(t, db.CreateTag(child1))
	assert.NoError(t, db.CreateTag(child2))
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
	item1 := &model.Item{Title: "title1"}
	assert.NoError(t, db.CreateItem(item1))
	itemFromDB, err := db.GetItem(1)
	assert.NoError(t, err)
	itemFromDB.Previews = []model.Preview{
		{
			Url: "preview1",
		},
		{
			Url: "preview2",
		},
	}
	err = db.UpdateItem(itemFromDB)
	assert.NoError(t, err)
	updatedFromDB, err := db.GetItem(1)
	assert.NoError(t, err)
	assert.Equal(t, len(itemFromDB.Previews), len(updatedFromDB.Previews))
	assert.Equal(t, updatedFromDB.Previews[0].ItemId, uint64(1))
	assert.Equal(t, updatedFromDB.Previews[1].ItemId, uint64(1))
	assert.Equal(t, updatedFromDB.Previews[0].Id, uint64(1))
	assert.Equal(t, updatedFromDB.Previews[1].Id, uint64(2))
	assert.Equal(t, updatedFromDB.Previews[0].Url, itemFromDB.Previews[0].Url)
	assert.Equal(t, updatedFromDB.Previews[1].Url, itemFromDB.Previews[1].Url)
}
