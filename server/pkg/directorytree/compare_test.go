package directorytree

import (
	"my-collection/server/pkg/model"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
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
	assert.Equal(t, 0, diff.ChangesTotal())
}

func TestCompareDirectoryAdded(t *testing.T) {
	diff := skeletonCompare(t, func(rootDir string) {
		assert.NoError(t, os.MkdirAll(filepath.Join(rootDir, "new-dir"), 0755))
	}, func(dig *model.MockDirectoryItemsGetter) {})
	assert.Equal(t, 1, diff.ChangesTotal())
	assert.Equal(t, ChangeType(DIRECTORY_ADDED), diff.AddedDirectories[0].ChangeType)
	assert.Equal(t, "new-dir", filepath.Base(diff.AddedDirectories[0].Path1))
	assert.Equal(t, "", diff.AddedDirectories[0].Path2)
}

func TestCompareDirectoryRemoved(t *testing.T) {
	diff := skeletonCompare(t, func(rootDir string) {
		assert.NoError(t, os.Remove(filepath.Join(rootDir, "1", "2", "3.2")))
	}, func(dig *model.MockDirectoryItemsGetter) {})
	assert.Equal(t, 1, diff.ChangesTotal())
	assert.Equal(t, ChangeType(DIRECTORY_REMOVED), diff.RemovedDirectories[0].ChangeType)
	assert.Equal(t, "3.2", filepath.Base(diff.RemovedDirectories[0].Path1))
	assert.Equal(t, "", diff.RemovedDirectories[0].Path2)
}

func TestCompareFileAdded(t *testing.T) {
	diff := skeletonCompare(t, func(rootDir string) {
		_, err := os.Create(filepath.Join(rootDir, "1", "2", "new-file"))
		assert.NoError(t, err)
	}, func(dig *model.MockDirectoryItemsGetter) {})
	assert.Equal(t, 1, diff.ChangesTotal())
	assert.Equal(t, ChangeType(FILE_ADDED), diff.AddedFiles[0].ChangeType)
	assert.Equal(t, "1/2/new-file", diff.AddedFiles[0].Path1)
	assert.Equal(t, "", diff.AddedFiles[0].Path2)
}

func TestCompareFileAddedToInnerDirectory(t *testing.T) {
	diff := skeletonCompare(t, func(rootDir string) {
		assert.NoError(t, os.MkdirAll(filepath.Join(rootDir, "1", "2", "completely", "new", "dir"), 0755))
		_, err := os.Create(filepath.Join(rootDir, "1", "2", "completely", "new", "dir", "new-file"))
		assert.NoError(t, err)
	}, func(dig *model.MockDirectoryItemsGetter) {})
	assert.Equal(t, 1, diff.ChangesTotal())
}

func TestCompareFileRemoved(t *testing.T) {
	diff := skeletonCompare(t, func(rootDir string) {
		assert.NoError(t, os.Remove(filepath.Join(rootDir, "1", "2", "3", "4", "file4-3")))
	}, func(dig *model.MockDirectoryItemsGetter) {})
	assert.Equal(t, 1, diff.ChangesTotal())
	assert.Equal(t, ChangeType(FILE_REMOVED), diff.RemovedFiles[0].ChangeType)
	assert.Equal(t, "1/2/3/4/file4-3", diff.RemovedFiles[0].Path1)
	assert.Equal(t, "", diff.RemovedFiles[0].Path2)
}

func TestCompareFileMoved(t *testing.T) {
	diff := skeletonCompare(t, func(rootDir string) {
		_, err := os.Create(filepath.Join(rootDir, "1", "2", "3", "file4-3"))
		assert.NoError(t, err)
		assert.NoError(t, os.Remove(filepath.Join(rootDir, "1", "2", "3", "4", "file4-3")))
	}, func(dig *model.MockDirectoryItemsGetter) {})
	assert.Equal(t, 1, diff.ChangesTotal())
	assert.Equal(t, ChangeType(FILE_MOVED), diff.MovedFiles[0].ChangeType)
	assert.Equal(t, "1/2/3/4/file4-3", diff.MovedFiles[0].Path1)
	assert.Equal(t, "1/2/3/file4-3", diff.MovedFiles[0].Path2)
}

func TestCompareDirectoryMoved(t *testing.T) {
	diff := skeletonCompare(t, func(rootDir string) {
		assert.NoError(t, os.Remove(filepath.Join(rootDir, "1", "2", "3.2")))
		assert.NoError(t, os.MkdirAll(filepath.Join(rootDir, "1", "2", "3", "3.2"), 0755))
	}, func(dig *model.MockDirectoryItemsGetter) {})
	assert.Equal(t, 1, diff.ChangesTotal())
	assert.Equal(t, ChangeType(DIRECTORY_MOVED), diff.MovedDirectories[0].ChangeType)
	assert.Equal(t, "1/2/3.2", diff.MovedDirectories[0].Path1)
	assert.Equal(t, "1/2/3/3.2", diff.MovedDirectories[0].Path2)
}

func TestCompareDirectoryMovedToInnerDirectory(t *testing.T) {
	diff := skeletonCompare(t, func(rootDir string) {
		assert.NoError(t, os.Remove(filepath.Join(rootDir, "1", "2", "3.2")))
		assert.NoError(t, os.MkdirAll(filepath.Join(rootDir, "1", "2", "completely", "new", "path", "3.2"), 0755))
	}, func(dig *model.MockDirectoryItemsGetter) {})
	assert.Equal(t, 2, diff.ChangesTotal())
	assert.Equal(t, ChangeType(DIRECTORY_MOVED), diff.MovedDirectories[0].ChangeType)
	assert.Equal(t, "1/2/3.2", diff.MovedDirectories[0].Path1)
	assert.Equal(t, "1/2/completely/new/path/3.2", diff.MovedDirectories[0].Path2)
}

func TestCompareDirectoryMovedToOuterDirectory(t *testing.T) {
	diff := skeletonCompare(t, func(rootDir string) {
		assert.NoError(t, os.Rename(filepath.Join(rootDir, "1", "2", "3", "4"), filepath.Join(rootDir, "1", "2", "4")))
		assert.NoError(t, os.Remove(filepath.Join(rootDir, "1", "2", "3")))
		assert.NoError(t, os.Remove(filepath.Join(rootDir, "1", "2", "3.2")))
		assert.NoError(t, os.Remove(filepath.Join(rootDir, "1", "2", "3.1")))
	}, func(dig *model.MockDirectoryItemsGetter) {})
	assert.Equal(t, 7, diff.ChangesTotal())
	assert.Equal(t, ChangeType(DIRECTORY_MOVED), diff.MovedDirectories[0].ChangeType)
	assert.Equal(t, "1/2/3/4", diff.MovedDirectories[0].Path1)
	assert.Equal(t, "1/2/4", diff.MovedDirectories[0].Path2)
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

	assert.Equal(t, 0, diff.ChangesTotal())
}

func TestCompareRootNotExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dr := model.NewMockDirectoryReader(ctrl)
	dr.EXPECT().GetAllDirectories().Return(&[]model.Directory{}, nil)
	dbRoot, err := BuildFromDb(dr, model.NewMockDirectoryItemsGetter(ctrl))
	assert.NoError(t, err)
	fsRoot, _ := skeleton(t, func(rootDir string) {}, func(dig *model.MockDirectoryItemsGetter) {})
	diff := Compare(fsRoot, dbRoot)
	assert.Equal(t, 1, diff.ChangesTotal())
}
