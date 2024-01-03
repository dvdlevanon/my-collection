package tag_annotations

import (
	"fmt"
	"my-collection/server/pkg/db"
	"my-collection/server/pkg/model"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupNewDb(t *testing.T, filename string) *db.Database {
	assert.NoError(t, os.MkdirAll(".tests", 0750))
	dbpath := fmt.Sprintf(".tests/%s", filename)
	_, err := os.Create(dbpath)
	assert.NoError(t, err)
	assert.NoError(t, os.Remove(dbpath))
	db, err := db.New(filepath.Join("", dbpath))
	assert.NoError(t, err)
	return db
}

func annotaionExists(annotations []model.TagAnnotation, title string) bool {
	for _, annotation := range annotations {
		if annotation.Title == title {
			return true
		}
	}

	return false
}

func TestAvailableTagAnnoations(t *testing.T) {
	db := setupNewDb(t, "get-available-tags.sqlite")
	var trw model.TagReaderWriter
	var tarw model.TagAnnotationReaderWriter
	trw = db
	tarw = db

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

	assert.NoError(t, db.CreateOrUpdateTag(&root))
	_, err := AddAnnotationToTag(trw, tarw, child1.Id, model.TagAnnotation{Title: "annotation1"})
	assert.NoError(t, err)
	_, err = AddAnnotationToTag(trw, tarw, child1.Id, model.TagAnnotation{Title: "annotation2"})
	assert.NoError(t, err)
	_, err = AddAnnotationToTag(trw, tarw, child2.Id, model.TagAnnotation{Title: "annotation1"})
	assert.NoError(t, err)
	_, err = AddAnnotationToTag(trw, tarw, child2.Id, model.TagAnnotation{Title: "annotation3"})
	assert.NoError(t, err)

	annotations, err := GetTagAvailableAnnotations(trw, tarw, root.Id)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(annotations))
	assert.True(t, annotaionExists(annotations, "annotation1"))
	assert.True(t, annotaionExists(annotations, "annotation2"))
	assert.True(t, annotaionExists(annotations, "annotation3"))

	_, err = AddAnnotationToTag(trw, tarw, child2.Id, model.TagAnnotation{Title: "annotation3"})
	assert.NoError(t, err)
}
