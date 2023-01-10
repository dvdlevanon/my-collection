package gallery

import (
	"fmt"
	"my-collection/server/pkg/db"
	"my-collection/server/pkg/model"
	"my-collection/server/pkg/storage"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupNewGallery(t *testing.T, filename string) *Gallery {
	assert.NoError(t, os.MkdirAll(".tests", 0750))
	dbpath := fmt.Sprintf(".tests/%s", filename)
	_, err := os.Create(dbpath)
	assert.NoError(t, err)
	assert.NoError(t, os.Remove(dbpath))
	db, err := db.New("", dbpath)
	assert.NoError(t, err)
	storage, err := storage.New("/tmp/root-directory/.storage")
	assert.NoError(t, err)
	return New(db, storage, "/tmp/root-directory")
}

func TestNormalizeUrl(t *testing.T) {
	gallery := setupNewGallery(t, "test-normalize-url.sqlite")
	gallery.CreateOrUpdateItem(&model.Item{Title: "title1", Url: "/tmp/root-directory/some-path/inner-path/file.ext"})
	item, err := gallery.GetItem(1)
	assert.NoError(t, err)
	assert.Equal(t, item.Url, "some-path/inner-path/file.ext")
}
