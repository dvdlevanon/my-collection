package directorytree

import (
	"my-collection/server/pkg/model"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/utils/pointer"
)

func TestStale(t *testing.T) {
	_, dbRoot := skeleton(t, func(rootDir string) {
		assert.NoError(t, os.MkdirAll(filepath.Join(rootDir, "1/2/ex/5/6"), 0755))
		_, err := os.Create(filepath.Join(rootDir, "1/2/ex/5/6/file"))
		assert.NoError(t, err)
		_, err = os.Create(filepath.Join(rootDir, "1/2/ex/5/6/file2"))
		assert.NoError(t, err)
	}, func(dig *model.MockDirectoryItemsGetter) {
		dig.EXPECT().GetBelongingItems("1/2/ex").Return(&[]model.Item{}, nil)
		dig.EXPECT().GetBelongingItems("1/2/ex/5").Return(&[]model.Item{}, nil)
		dig.EXPECT().GetBelongingItems("1/2/ex/5/6").Return(&[]model.Item{
			{Title: "file"},
		}, nil)
	},
		model.Directory{Path: "1/2/ex", Excluded: pointer.Bool(true)},
		model.Directory{Path: "1/2/ex/5"},
		model.Directory{Path: "1/2/ex/5/6"})

	stales := FindStales(dbRoot)
	assert.NotNil(t, stales)
	assert.Equal(t, 3, len(stales.Dirs))
	assert.Equal(t, 1, len(stales.Files))
}
