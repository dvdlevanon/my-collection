package gallery

import (
	"my-collection/server/pkg/model"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var importExportTestJson = "{\"items\":[{\"id\":10,\"title\":\"item10\",\"tags\":[{\"id\":10},{\"id\":20}]},{\"id\":20,\"title\":\"item20\",\"tags\":[{\"id\":10},{\"id\":30}]},{\"id\":30,\"title\":\"item30\"}],\"tags\":[{\"id\":10,\"title\":\"tag10\",\"items\":[{\"id\":10},{\"id\":20}]},{\"id\":20,\"title\":\"tag20\",\"items\":[{\"id\":10}]},{\"id\":30,\"title\":\"tag30\",\"items\":[{\"id\":20}]}]}"

func TestExport(t *testing.T) {
	gallery := setupNewGallery(t, "test-export.sqlite")

	item10 := model.Item{
		Id:    10,
		Title: "item10",
		Tags: []*model.Tag{
			{
				Id:    10,
				Title: "tag10",
			},
			{
				Id:    20,
				Title: "tag20",
			},
		},
	}

	item20 := model.Item{
		Id:    20,
		Title: "item20",
		Tags: []*model.Tag{
			{
				Id: 10,
			},
			{
				Id:    30,
				Title: "tag30",
			},
		},
	}

	item30 := model.Item{
		Id:    30,
		Title: "item30",
	}

	assert.NoError(t, gallery.CreateOrUpdateItem(&item10))
	assert.NoError(t, gallery.CreateOrUpdateItem(&item20))
	assert.NoError(t, gallery.CreateOrUpdateItem(&item30))

	out := strings.Builder{}
	assert.NoError(t, gallery.Export(&out))
	assert.Equal(t, importExportTestJson, out.String())
}

func TestImport(t *testing.T) {
	gallery := setupNewGallery(t, "test-import.sqlite")
	assert.NoError(t, gallery.Import([]byte(importExportTestJson)))

	item10, err := gallery.GetItem(10)
	assert.NoError(t, err)
	assert.Equal(t, item10.Title, "item10")
	assert.Equal(t, len(item10.Tags), 2)
	assert.Equal(t, item10.Tags[0].Id, uint64(10))
	assert.Equal(t, item10.Tags[1].Id, uint64(20))

	item20, err := gallery.GetItem(20)
	assert.NoError(t, err)
	assert.Equal(t, item20.Title, "item20")
	assert.Equal(t, len(item20.Tags), 2)
	assert.Equal(t, item20.Tags[0].Id, uint64(10))
	assert.Equal(t, item20.Tags[1].Id, uint64(30))

	item30, err := gallery.GetItem(30)
	assert.NoError(t, err)
	assert.Equal(t, item30.Title, "item30")
	assert.Equal(t, len(item30.Tags), 0)

	tag10, err := gallery.GetTag(10)
	assert.NoError(t, err)
	assert.Equal(t, tag10.Title, "tag10")
	assert.Equal(t, len(tag10.Items), 2)
	assert.Equal(t, tag10.Items[0].Id, uint64(10))
	assert.Equal(t, tag10.Items[1].Id, uint64(20))

	tag20, err := gallery.GetTag(20)
	assert.NoError(t, err)
	assert.Equal(t, tag20.Title, "tag20")
	assert.Equal(t, 1, len(tag20.Items))
	assert.Equal(t, tag20.Items[0].Id, uint64(10))

	tag30, err := gallery.GetTag(30)
	assert.NoError(t, err)
	assert.Equal(t, tag30.Title, "tag30")
	assert.Equal(t, 1, len(tag30.Items))
	assert.Equal(t, tag30.Items[0].Id, uint64(20))
}
