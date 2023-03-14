package directorytree

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func buildTestFs(t *testing.T, rootDir string) {
	var err error
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
	_, err = os.Create(filepath.Join(rootDir, "1", "2", "3", "4", "file4-excluded"))
	assert.NoError(t, err)
}

func TestBuildFromPath(t *testing.T) {
	rootDir, err := os.MkdirTemp("", "mc-build-from-fs-*")
	assert.NoError(t, err)
	buildTestFs(t, rootDir)

	root, err := BuildFromPath(rootDir, func(path string) bool {
		return !strings.HasSuffix(path, "file4-excluded")
	})
	assert.NoError(t, err)
	assert.NotNil(t, root)
	assert.Equal(t, 3, len(root.getOrCreateChild("1/2/3/4").Files))
}
