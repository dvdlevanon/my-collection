package directories

import (
	"my-collection/server/pkg/model"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestGetAllDirectoriesWithCache(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// error get
	drErrorGet := model.NewMockDirectoryReader(ctrl)
	drErrorGet.EXPECT().GetAllDirectories().Return(nil, errors.Errorf("test error"))
	_, err := GetAllDirectoriesWithCache(drErrorGet)
	assert.Error(t, err)

	// setup mock
	dirs := []model.Directory{{Path: "dir1"}, {Path: "dir2"}}
	dr := model.NewMockDirectoryReader(ctrl)
	dr.EXPECT().GetAllDirectories().AnyTimes().Return(&dirs, nil)

	// clean cache
	dirs1, err := GetAllDirectoriesWithCache(dr)
	assert.NoError(t, err)
	assert.Equal(t, len(dirs), len(*dirs1))

	// change orig
	dirs = append(dirs, model.Directory{Path: "dir3"})

	// already cached
	dirs2, err := GetAllDirectoriesWithCache(dr)
	assert.NoError(t, err)
	assert.Equal(t, len(*dirs1), len(*dirs2))

	// clean cache again
	directoriesCache.Flush()
	dirs3, err := GetAllDirectoriesWithCache(dr)
	assert.NoError(t, err)
	assert.Equal(t, len(dirs), len(*dirs3))
}
