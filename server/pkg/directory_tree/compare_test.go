package directorytree

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type extendsFs func(rootDir string)

func skeleton(t *testing.T, extendsFs extendsFs) *Diff {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	rootDir, err := os.MkdirTemp("", "mc-build-from-fs-*")
	assert.NoError(t, err)
	buildTestFs(t, rootDir)
	extendsFs(rootDir)

	dbRoot, err := BuildFromDb(buildTestDirectoryReader(ctrl), buildTestDirectoryItemsGetter(ctrl))
	assert.NoError(t, err)

	fsRoot, err := BuildFromPath(rootDir, func(path string) bool {
		return !strings.HasSuffix(path, "file4-excluded")
	})
	assert.NoError(t, err)

	diff := Compare(fsRoot, dbRoot)
	assert.NotNil(t, diff)
	return diff

}

func TestCompareEqual(t *testing.T) {
	diff := skeleton(t, func(rootDir string) {})
	assert.Equal(t, 0, len(diff.Changes))
}

func TestCompareDirectoryAdded(t *testing.T) {
	diff := skeleton(t, func(rootDir string) {
		assert.NoError(t, os.MkdirAll(filepath.Join(rootDir, "new-dir"), 0755))
	})
	assert.Equal(t, 1, len(diff.Changes))
	assert.Equal(t, ChangeType(DIRECTORY_ADDED), diff.Changes[0].ChangeType)
	assert.Equal(t, "new-dir", filepath.Base(diff.Changes[0].Path1))
	assert.Equal(t, "", diff.Changes[0].Path2)
}

func TestCompareDirectoryRemoved(t *testing.T) {
	diff := skeleton(t, func(rootDir string) {
		assert.NoError(t, os.Remove(filepath.Join(rootDir, "1", "2", "3.2")))
	})
	assert.Equal(t, 1, len(diff.Changes))
	assert.Equal(t, ChangeType(DIRECTORY_REMOVED), diff.Changes[0].ChangeType)
	assert.Equal(t, "3.2", filepath.Base(diff.Changes[0].Path1))
	assert.Equal(t, "", diff.Changes[0].Path2)
}

func TestCompareFileAdded(t *testing.T) {
	diff := skeleton(t, func(rootDir string) {
		_, err := os.Create(filepath.Join(rootDir, "1", "2", "new-file"))
		assert.NoError(t, err)
	})
	assert.Equal(t, 1, len(diff.Changes))
	assert.Equal(t, ChangeType(FILE_ADDED), diff.Changes[0].ChangeType)
	assert.Equal(t, "1/2/new-file", diff.Changes[0].Path1)
	assert.Equal(t, "", diff.Changes[0].Path2)
}

func TestCompareFileAddedToInnerDirectory(t *testing.T) {
	diff := skeleton(t, func(rootDir string) {
		assert.NoError(t, os.MkdirAll(filepath.Join(rootDir, "1", "2", "completely", "new", "dir"), 0755))
		_, err := os.Create(filepath.Join(rootDir, "1", "2", "completely", "new", "dir", "new-file"))
		assert.NoError(t, err)
	})
	assert.Equal(t, 4, len(diff.Changes))
}

func TestCompareFileRemoved(t *testing.T) {
	diff := skeleton(t, func(rootDir string) {
		assert.NoError(t, os.Remove(filepath.Join(rootDir, "1", "2", "3", "4", "file4-3")))
	})
	assert.Equal(t, 1, len(diff.Changes))
	assert.Equal(t, ChangeType(FILE_REMOVED), diff.Changes[0].ChangeType)
	assert.Equal(t, "1/2/3/4/file4-3", diff.Changes[0].Path1)
	assert.Equal(t, "", diff.Changes[0].Path2)
}

func TestCompareFileMoved(t *testing.T) {
	diff := skeleton(t, func(rootDir string) {
		_, err := os.Create(filepath.Join(rootDir, "1", "2", "3", "file4-3"))
		assert.NoError(t, err)
		assert.NoError(t, os.Remove(filepath.Join(rootDir, "1", "2", "3", "4", "file4-3")))
	})
	assert.Equal(t, 1, len(diff.Changes))
	assert.Equal(t, ChangeType(FILE_MOVED), diff.Changes[0].ChangeType)
	assert.Equal(t, "1/2/3/4/file4-3", diff.Changes[0].Path1)
	assert.Equal(t, "1/2/3/file4-3", diff.Changes[0].Path2)
}

func TestCompareDirectoryMoved(t *testing.T) {
	diff := skeleton(t, func(rootDir string) {
		assert.NoError(t, os.Remove(filepath.Join(rootDir, "1", "2", "3.2")))
		assert.NoError(t, os.MkdirAll(filepath.Join(rootDir, "1", "2", "3", "3.2"), 0755))
	})
	assert.Equal(t, 1, len(diff.Changes))
	assert.Equal(t, ChangeType(DIRECTORY_MOVED), diff.Changes[0].ChangeType)
	assert.Equal(t, "1/2/3.2", diff.Changes[0].Path1)
	assert.Equal(t, "1/2/3/3.2", diff.Changes[0].Path2)
}

func TestCompareDirectoryMovedToInnerDirectory(t *testing.T) {
	diff := skeleton(t, func(rootDir string) {
		assert.NoError(t, os.Remove(filepath.Join(rootDir, "1", "2", "3.2")))
		assert.NoError(t, os.MkdirAll(filepath.Join(rootDir, "1", "2", "completely", "new", "path", "3.2"), 0755))
	})
	assert.Equal(t, 4, len(diff.Changes))
	found := false
	for _, change := range diff.Changes {
		if change.ChangeType != DIRECTORY_MOVED {
			continue
		}
		found = true
		assert.Equal(t, ChangeType(DIRECTORY_MOVED), change.ChangeType)
		assert.Equal(t, "1/2/3.2", change.Path1)
		assert.Equal(t, "1/2/completely/new/path/3.2", change.Path2)
	}
	assert.True(t, found)
}

func TestCompareDirectoryMovedToOuterDirectory(t *testing.T) {
	diff := skeleton(t, func(rootDir string) {
		assert.NoError(t, os.Rename(filepath.Join(rootDir, "1", "2", "3", "4"), filepath.Join(rootDir, "1", "2", "4")))
		assert.NoError(t, os.Remove(filepath.Join(rootDir, "1", "2", "3")))
		assert.NoError(t, os.Remove(filepath.Join(rootDir, "1", "2", "3.2")))
		assert.NoError(t, os.Remove(filepath.Join(rootDir, "1", "2", "3.1")))
	})
	assert.Equal(t, 7, len(diff.Changes))
	found := false
	for _, change := range diff.Changes {
		if change.ChangeType != DIRECTORY_MOVED {
			continue
		}
		found = true
		assert.Equal(t, ChangeType(DIRECTORY_MOVED), change.ChangeType)
		assert.Equal(t, "1/2/3/4", change.Path1)
		assert.Equal(t, "1/2/4", change.Path2)
	}
	assert.True(t, found)
}
