package directorytree

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildFromPath(t *testing.T) {
	rootDir, err := os.MkdirTemp("", "mc-build-from-fs-*")
	assert.NoError(t, err)
	assert.NoError(t, os.MkdirAll(filepath.Join(rootDir, "1", "2", "3", "4"), 0755))
	assert.NoError(t, os.MkdirAll(filepath.Join(rootDir, "1", "2", "3.1"), 0755))
	assert.NoError(t, os.MkdirAll(filepath.Join(rootDir, "1", "2", "3.2"), 0755))
	_, err = os.Create(filepath.Join(rootDir, "1", "2", "file2-1"))
	assert.NoError(t, err)
	_, err = os.Create(filepath.Join(rootDir, "1", "2", "3", "4", "file4-1"))
	assert.NoError(t, err)
	_, err = os.Create(filepath.Join(rootDir, "1", "2", "3", "4", "file4-2"))
	assert.NoError(t, err)
	_, err = os.Create(filepath.Join(rootDir, "1", "2", "3", "4", "file4-3"))
	assert.NoError(t, err)

	tree, err := BuildFromPath(rootDir)
	assert.NoError(t, err)
	assert.NotNil(t, tree.Root)
	assert.Equal(t, 3, len(tree.Root.getOrCreateChild("1/2/3/4").Files))
}
