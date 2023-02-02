package model

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/utils/pointer"
)

func TestTagMarshal(t *testing.T) {
	parent := Tag{
		Id:    1,
		Title: "parent",
	}

	child := Tag{
		Id:       2,
		Title:    "child",
		ParentID: &parent.Id,
		Annotations: []*TagAnnotation{
			{
				Id:    1,
				Title: "annotation1",
			},
		},
	}

	parent.Children = append(parent.Children, &child)

	bytes, err := json.Marshal(parent)
	assert.NoError(t, err)
	assert.Equal(t, `{"id":1,"title":"parent","children":[{"id":2,"title":"child","parentId":1,"tags_annotations":[{"id":1,"title":"annotation1"}]}]}`, string(bytes))
}

func TestTagUnmarshal(t *testing.T) {
	jsonTag := `{"id":1,"title":"parent","children":[{"id":2,"title":"child","parentId":1,"tags_annotations":[{"id":1,"title":"annotation1"}]}]}`
	var tag Tag
	assert.NoError(t, json.Unmarshal([]byte(jsonTag), &tag))
	assert.Equal(t, uint64(1), tag.Id)
	assert.Equal(t, uint64(2), tag.Children[0].Id)
	assert.Equal(t, 1, len(tag.Children[0].Annotations))
	assert.Equal(t, "annotation1", tag.Children[0].Annotations[0].Title)
}

func TestItemMarshal(t *testing.T) {
	item := Item{
		Id:    1,
		Title: "item",
		Url:   "url",
		Covers: []Cover{
			{
				Id:     20,
				ItemId: 1,
				Url:    "cover1",
			},
			{
				Id:     21,
				ItemId: 1,
				Url:    "cover2",
			},
		},
		DurationSeconds: 34324,
		Width:           1920,
		Height:          1080,
		Tags: []*Tag{
			{
				Id: 10,
			},
			{
				Id: 11,
			},
			{
				Id: 12,
			},
		},
	}

	bytes, err := json.Marshal(item)
	assert.NoError(t, err)
	assert.Equal(t, `{"id":1,"title":"item","duration_seconds":34324,"width":1920,"height":1080,"url":"url","covers":[{"id":20,"url":"cover1","itemId":1},{"id":21,"url":"cover2","itemId":1}],"tags":[{"id":10},{"id":11},{"id":12}]}`, string(bytes))
}

func TestDirectoryMarshal(t *testing.T) {
	dir := Directory{
		Path:       "path/to/dir",
		FilesCount: 1,
		LastSynced: 1675352131869,
		Excluded:   pointer.Bool(false),
		Tags:       []*Tag{{Id: 1}, {Id: 2}},
	}

	bytes, err := json.Marshal(dir)
	assert.NoError(t, err)
	assert.Equal(t, `{"path":"path/to/dir","excluded":false,"tags":[{"id":1},{"id":2}],"filesCount":1,"lastSynced":1675352131869}`, string(bytes))
}
