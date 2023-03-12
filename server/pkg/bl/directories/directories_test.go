package directories

import (
	"my-collection/server/pkg/model"
	"testing"

	"github.com/go-errors/errors"
	"github.com/stretchr/testify/assert"
)

type testDirectoriesReaderWriter struct {
	directories []model.Directory
	errorGet    bool
}

func (d *testDirectoriesReaderWriter) GetDirectory(conds ...interface{}) (*model.Directory, error) {
	if d.errorGet {
		return nil, errors.Errorf("test error")
	}
	return &(d.directories[0]), nil
}

func (d *testDirectoriesReaderWriter) GetDirectories(conds ...interface{}) (*[]model.Directory, error) {
	return d.GetAllDirectories()
}

func (d *testDirectoriesReaderWriter) GetAllDirectories() (*[]model.Directory, error) {
	if d.errorGet {
		return nil, errors.Errorf("test error")
	}
	cloned := d.directories
	return &cloned, nil
}

func TestGetAllDirectoriesWithCache(t *testing.T) {
	dr := testDirectoriesReaderWriter{directories: []model.Directory{{Path: "dir1"}, {Path: "dir2"}}}

	// error get
	drErrorGet := testDirectoriesReaderWriter{errorGet: true}
	_, err := GetAllDirectoriesWithCache(&drErrorGet)
	assert.Error(t, err)

	// clean cache
	dirs1, err := GetAllDirectoriesWithCache(&dr)
	assert.NoError(t, err)
	assert.Equal(t, len(dr.directories), len(*dirs1))

	// change orig
	dr.directories = append(dr.directories, model.Directory{Path: "dir3"})

	// already cached
	dirs2, err := GetAllDirectoriesWithCache(&dr)
	assert.NoError(t, err)
	assert.Equal(t, len(*dirs1), len(*dirs2))

	// clean cache again
	directoriesCache.Flush()
	dirs3, err := GetAllDirectoriesWithCache(&dr)
	assert.NoError(t, err)
	assert.Equal(t, len(dr.directories), len(*dirs3))
}
