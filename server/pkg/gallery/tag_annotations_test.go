package gallery

import (
	"my-collection/server/pkg/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func annotaionExists(annotations []model.TagAnnotation, title string) bool {
	for _, annotation := range annotations {
		if annotation.Title == title {
			return true
		}
	}

	return false
}

func TestAvailableTagAnnoations(t *testing.T) {
	gallery := setupNewGallery(t, "get-available-tags.sqlite")
	child1 := model.Tag{
		Title: "child1",
	}

	child2 := model.Tag{
		Title: "child2",
	}

	assert.NoError(t, gallery.CreateTag(&child1))
	assert.NoError(t, gallery.CreateTag(&child2))

	root := model.Tag{
		Title:    "root",
		Children: []*model.Tag{&child1, &child2},
	}

	assert.NoError(t, gallery.CreateTag(&root))
	_, err := gallery.AddAnnotationToTag(child1.Id, model.TagAnnotation{Title: "annotation1"})
	assert.NoError(t, err)
	_, err = gallery.AddAnnotationToTag(child1.Id, model.TagAnnotation{Title: "annotation2"})
	assert.NoError(t, err)
	_, err = gallery.AddAnnotationToTag(child2.Id, model.TagAnnotation{Title: "annotation1"})
	assert.NoError(t, err)
	_, err = gallery.AddAnnotationToTag(child2.Id, model.TagAnnotation{Title: "annotation3"})
	assert.NoError(t, err)

	annotations, err := gallery.GetTagAvailableAnnotations(root.Id)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(annotations))
	assert.True(t, annotaionExists(annotations, "annotation1"))
	assert.True(t, annotaionExists(annotations, "annotation2"))
	assert.True(t, annotaionExists(annotations, "annotation3"))

	_, err = gallery.AddAnnotationToTag(child2.Id, model.TagAnnotation{Title: "annotation3"})
	assert.NoError(t, err)
}
