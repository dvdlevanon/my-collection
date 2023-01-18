package model

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
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
	}

	parent.Children = append(parent.Children, &child)

	bytes, err := json.Marshal(parent)
	assert.NoError(t, err)
	assert.Equal(t, string(bytes), `{"id":1,"title":"parent","children":[{"id":2,"title":"child","parentId":1}]}`)
}

func TestTagUnmarshal(t *testing.T) {
	jsonTag := `{"id":1,"title":"parent","children":[{"id":2,"title":"child","parentId":1}]}`
	var tag Tag
	assert.NoError(t, json.Unmarshal([]byte(jsonTag), &tag))
	assert.Equal(t, tag.Id, uint64(1))
	assert.Equal(t, tag.Children[0].Id, uint64(2))
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
	assert.Equal(t, string(bytes), `{"id":1,"title":"item","url":"url","covers":[{"id":20,"url":"cover1","itemId":1},{"id":21,"url":"cover2","itemId":1}],"tags":[{"id":10},{"id":11},{"id":12}]}`)
}
