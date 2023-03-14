package directorytree

import (
	"my-collection/server/pkg/model"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"k8s.io/utils/pointer"
)

type extendsFs func(rootDir string)
type extendsDig func(dig *model.MockDirectoryItemsGetter)

func skeleton(t *testing.T, extendsFs extendsFs, extendsDig extendsDig,
	additionalDbDirs ...model.Directory) (*DirectoryNode, *DirectoryNode) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	rootDir, err := os.MkdirTemp("", "mc-build-from-fs-*")
	assert.NoError(t, err)
	buildTestFs(t, rootDir)
	extendsFs(rootDir)

	dig := buildTestDirectoryItemsGetter(ctrl)
	extendsDig(dig)
	dbRoot, err := BuildFromDb(buildTestDirectoryReader(ctrl, additionalDbDirs...), dig)
	assert.NoError(t, err)

	fsRoot, err := BuildFromPath(rootDir, func(path string) bool {
		return !strings.HasSuffix(path, "file4-excluded")
	})
	assert.NoError(t, err)
	return fsRoot, dbRoot
}

func skeletonCompare(t *testing.T, extendsFs extendsFs, extendsDig extendsDig, additionalDbDirs ...model.Directory) *Diff {
	fsRoot, dbRoot := skeleton(t, extendsFs, extendsDig, additionalDbDirs...)
	diff := Compare(fsRoot, dbRoot)
	assert.NotNil(t, diff)
	return diff
}

func TestCompareEqual(t *testing.T) {
	diff := skeletonCompare(t, func(rootDir string) {}, func(dig *model.MockDirectoryItemsGetter) {})
	assert.Equal(t, 0, len(diff.Changes))
}

func TestCompareDirectoryAdded(t *testing.T) {
	diff := skeletonCompare(t, func(rootDir string) {
		assert.NoError(t, os.MkdirAll(filepath.Join(rootDir, "new-dir"), 0755))
	}, func(dig *model.MockDirectoryItemsGetter) {})
	assert.Equal(t, 1, len(diff.Changes))
	assert.Equal(t, ChangeType(DIRECTORY_ADDED), diff.Changes[0].ChangeType)
	assert.Equal(t, "new-dir", filepath.Base(diff.Changes[0].Path1))
	assert.Equal(t, "", diff.Changes[0].Path2)
}

func TestCompareDirectoryRemoved(t *testing.T) {
	diff := skeletonCompare(t, func(rootDir string) {
		assert.NoError(t, os.Remove(filepath.Join(rootDir, "1", "2", "3.2")))
	}, func(dig *model.MockDirectoryItemsGetter) {})
	assert.Equal(t, 1, len(diff.Changes))
	assert.Equal(t, ChangeType(DIRECTORY_REMOVED), diff.Changes[0].ChangeType)
	assert.Equal(t, "3.2", filepath.Base(diff.Changes[0].Path1))
	assert.Equal(t, "", diff.Changes[0].Path2)
}

func TestCompareFileAdded(t *testing.T) {
	diff := skeletonCompare(t, func(rootDir string) {
		_, err := os.Create(filepath.Join(rootDir, "1", "2", "new-file"))
		assert.NoError(t, err)
	}, func(dig *model.MockDirectoryItemsGetter) {})
	assert.Equal(t, 1, len(diff.Changes))
	assert.Equal(t, ChangeType(FILE_ADDED), diff.Changes[0].ChangeType)
	assert.Equal(t, "1/2/new-file", diff.Changes[0].Path1)
	assert.Equal(t, "", diff.Changes[0].Path2)
}

func TestCompareFileAddedToInnerDirectory(t *testing.T) {
	diff := skeletonCompare(t, func(rootDir string) {
		assert.NoError(t, os.MkdirAll(filepath.Join(rootDir, "1", "2", "completely", "new", "dir"), 0755))
		_, err := os.Create(filepath.Join(rootDir, "1", "2", "completely", "new", "dir", "new-file"))
		assert.NoError(t, err)
	}, func(dig *model.MockDirectoryItemsGetter) {})
	assert.Equal(t, 4, len(diff.Changes))
}

func TestCompareFileRemoved(t *testing.T) {
	diff := skeletonCompare(t, func(rootDir string) {
		assert.NoError(t, os.Remove(filepath.Join(rootDir, "1", "2", "3", "4", "file4-3")))
	}, func(dig *model.MockDirectoryItemsGetter) {})
	assert.Equal(t, 1, len(diff.Changes))
	assert.Equal(t, ChangeType(FILE_REMOVED), diff.Changes[0].ChangeType)
	assert.Equal(t, "1/2/3/4/file4-3", diff.Changes[0].Path1)
	assert.Equal(t, "", diff.Changes[0].Path2)
}

func TestCompareFileMoved(t *testing.T) {
	diff := skeletonCompare(t, func(rootDir string) {
		_, err := os.Create(filepath.Join(rootDir, "1", "2", "3", "file4-3"))
		assert.NoError(t, err)
		assert.NoError(t, os.Remove(filepath.Join(rootDir, "1", "2", "3", "4", "file4-3")))
	}, func(dig *model.MockDirectoryItemsGetter) {})
	assert.Equal(t, 1, len(diff.Changes))
	assert.Equal(t, ChangeType(FILE_MOVED), diff.Changes[0].ChangeType)
	assert.Equal(t, "1/2/3/4/file4-3", diff.Changes[0].Path1)
	assert.Equal(t, "1/2/3/file4-3", diff.Changes[0].Path2)
}

func TestCompareDirectoryMoved(t *testing.T) {
	diff := skeletonCompare(t, func(rootDir string) {
		assert.NoError(t, os.Remove(filepath.Join(rootDir, "1", "2", "3.2")))
		assert.NoError(t, os.MkdirAll(filepath.Join(rootDir, "1", "2", "3", "3.2"), 0755))
	}, func(dig *model.MockDirectoryItemsGetter) {})
	assert.Equal(t, 1, len(diff.Changes))
	assert.Equal(t, ChangeType(DIRECTORY_MOVED), diff.Changes[0].ChangeType)
	assert.Equal(t, "1/2/3.2", diff.Changes[0].Path1)
	assert.Equal(t, "1/2/3/3.2", diff.Changes[0].Path2)
}

func TestCompareDirectoryMovedToInnerDirectory(t *testing.T) {
	diff := skeletonCompare(t, func(rootDir string) {
		assert.NoError(t, os.Remove(filepath.Join(rootDir, "1", "2", "3.2")))
		assert.NoError(t, os.MkdirAll(filepath.Join(rootDir, "1", "2", "completely", "new", "path", "3.2"), 0755))
	}, func(dig *model.MockDirectoryItemsGetter) {})
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
	diff := skeletonCompare(t, func(rootDir string) {
		assert.NoError(t, os.Rename(filepath.Join(rootDir, "1", "2", "3", "4"), filepath.Join(rootDir, "1", "2", "4")))
		assert.NoError(t, os.Remove(filepath.Join(rootDir, "1", "2", "3")))
		assert.NoError(t, os.Remove(filepath.Join(rootDir, "1", "2", "3.2")))
		assert.NoError(t, os.Remove(filepath.Join(rootDir, "1", "2", "3.1")))
	}, func(dig *model.MockDirectoryItemsGetter) {})
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

func TestCompareExcluded(t *testing.T) {
	diff := skeletonCompare(t, func(rootDir string) {
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

	assert.Equal(t, 0, len(diff.Changes))
}
