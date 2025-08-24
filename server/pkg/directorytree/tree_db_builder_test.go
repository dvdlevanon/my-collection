package directorytree

import (
	"my-collection/server/pkg/model"
	"testing"

	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
)

func buildTestDirectoryReader(ctrl *gomock.Controller, additionalDirs ...model.Directory) model.DirectoryReader {
	dr := model.NewMockDirectoryReader(ctrl)
	dirs := append([]model.Directory{
		{Path: "1/2/3/4"},
		{Path: "1"},
		{Path: "1/2/3.1"},
		{Path: "1/2/3.2"},
		{Path: "1/2"},
	}, additionalDirs...)
	dr.EXPECT().GetAllDirectories().Return(&dirs, nil)
	return dr
}

func buildTestDirectoryItemsGetter(ctrl *gomock.Controller) *model.MockDirectoryItemsGetter {
	dig := model.NewMockDirectoryItemsGetter(ctrl)
	dig.EXPECT().GetBelongingItems("1").Return(&[]model.Item{}, nil)
	dig.EXPECT().GetBelongingItems("1/2").Return(&[]model.Item{
		{Title: "file2-1"},
	}, nil)
	dig.EXPECT().GetBelongingItems("1/2/3.1").Return(&[]model.Item{}, nil)
	dig.EXPECT().GetBelongingItems("1/2/3.2").Return(&[]model.Item{}, nil)
	dig.EXPECT().GetBelongingItems("1/2/3/4").Return(&[]model.Item{
		{Title: "file4-1"},
		{Title: "file4-2"},
		{Title: "file4-3"},
	}, nil)

	return dig
}

func TestBuildFromDb(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	root, err := BuildFromDb(buildTestDirectoryReader(ctrl), buildTestDirectoryItemsGetter(ctrl))
	assert.NoError(t, err)
	assert.NotNil(t, root)
	assert.Equal(t, 3, len(root.getOrCreateChild("1/2/3/4").Files))
}

func TestGetOrCreateChild(t *testing.T) {
	root := createDirectoryNode(nil, "1")
	deep := root.getOrCreateChild("2/3/4/5")
	assert.Equal(t, "5", deep.Title)
	assert.Equal(t, "4", deep.Parent.Title)
	assert.Equal(t, "3", deep.Parent.Parent.Title)
	assert.Equal(t, "2", deep.Parent.Parent.Parent.Title)
	assert.Equal(t, "1", deep.Parent.Parent.Parent.Parent.Title)
	assert.Equal(t, root, deep.Parent.Parent.Parent.Parent)
	assert.Nil(t, root.Parent)
}
