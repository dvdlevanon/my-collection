package backup

import (
	"fmt"
	"my-collection/server/pkg/db"
	"my-collection/server/pkg/model"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var importExportTestJson = "{\"items\":[{\"id\":10,\"title\":\"item10\",\"origin\":\"origin\",\"tags\":[{\"id\":10,\"title\":\"tag10\"},{\"id\":20,\"title\":\"tag20\"}]},{\"id\":20,\"title\":\"item20\",\"origin\":\"origin\",\"tags\":[{\"id\":10,\"title\":\"tag10\"},{\"id\":30,\"title\":\"tag30\"}]},{\"id\":30,\"title\":\"item30\",\"origin\":\"origin\"}],\"tags\":[{\"id\":10,\"title\":\"tag10\",\"items\":[{\"id\":10},{\"id\":20}]},{\"id\":20,\"title\":\"tag20\",\"items\":[{\"id\":10}]},{\"id\":30,\"title\":\"tag30\",\"items\":[{\"id\":20}]}]}"

func setupNewDb(t *testing.T, filename string) db.Database {
	assert.NoError(t, os.MkdirAll(".tests", 0750))
	dbpath := fmt.Sprintf(".tests/%s", filename)
	_, err := os.Create(dbpath)
	assert.NoError(t, err)
	assert.NoError(t, os.Remove(dbpath))
	db, err := db.New(filepath.Join("", dbpath))
	assert.NoError(t, err)
	return db
}

func TestExport(t *testing.T) {
	db := setupNewDb(t, "test-export.sqlite")

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

	assert.NoError(t, db.CreateOrUpdateItem(&item10))
	assert.NoError(t, db.CreateOrUpdateItem(&item20))
	assert.NoError(t, db.CreateOrUpdateItem(&item30))

	out := strings.Builder{}
	assert.NoError(t, Export(db, db, &out))
	assert.Equal(t, importExportTestJson, out.String())
}
