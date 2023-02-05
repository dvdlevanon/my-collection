package gallery

import (
	"my-collection/server/pkg/model"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var importExportTestJson = "{\"items\":[{\"id\":10,\"title\":\"item10\",\"origin\":\"origin\",\"tags\":[{\"id\":10},{\"id\":20}]},{\"id\":20,\"title\":\"item20\",\"origin\":\"origin\",\"tags\":[{\"id\":10},{\"id\":30}]},{\"id\":30,\"title\":\"item30\",\"origin\":\"origin\"}],\"tags\":[{\"id\":10,\"title\":\"tag10\",\"items\":[{\"id\":10},{\"id\":20}]},{\"id\":20,\"title\":\"tag20\",\"items\":[{\"id\":10}]},{\"id\":30,\"title\":\"tag30\",\"items\":[{\"id\":20}]}]}"

func TestExport(t *testing.T) {
	gallery := setupNewGallery(t, "test-export.sqlite")

	item10 := model.Item{
		Id:     10,
		Origin: "origin",
		Title:  "item10",
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
		Id:     20,
		Origin: "origin",
		Title:  "item20",
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
		Id:     30,
		Origin: "origin",
		Title:  "item30",
	}

	assert.NoError(t, gallery.CreateOrUpdateItem(&item10))
	assert.NoError(t, gallery.CreateOrUpdateItem(&item20))
	assert.NoError(t, gallery.CreateOrUpdateItem(&item30))

	out := strings.Builder{}
	assert.NoError(t, gallery.Export(&out))
	assert.Equal(t, importExportTestJson, out.String())
}
