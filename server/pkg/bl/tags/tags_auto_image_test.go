package tags

import (
	"errors"
	"io/fs"
	"my-collection/server/pkg/bl/directories"
	"my-collection/server/pkg/model"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

// Mock fs.DirEntry for testing
type mockDirEntry struct {
	name  string
	isDir bool
}

func (m *mockDirEntry) Name() string               { return m.name }
func (m *mockDirEntry) IsDir() bool                { return m.isDir }
func (m *mockDirEntry) Type() fs.FileMode          { return 0 }
func (m *mockDirEntry) Info() (fs.FileInfo, error) { return nil, nil }

func setupTempDirWithImages(t *testing.T) (string, func()) {
	tempDir, err := os.MkdirTemp("", "tags-auto-image-test-*")
	assert.NoError(t, err)

	// Create test image directories and files
	bannersDir := filepath.Join(tempDir, "banners")
	iconsDir := filepath.Join(tempDir, "icons")

	assert.NoError(t, os.MkdirAll(bannersDir, 0755))
	assert.NoError(t, os.MkdirAll(iconsDir, 0755))

	// Create test image files
	testFiles := []string{
		filepath.Join(bannersDir, "test-tag.jpg"),
		filepath.Join(bannersDir, "another-tag.png"),
		filepath.Join(iconsDir, "icon.svg"),
	}

	for _, file := range testFiles {
		f, err := os.Create(file)
		assert.NoError(t, err)
		f.WriteString("test image content")
		f.Close()
	}

	cleanup := func() {
		os.RemoveAll(tempDir)
	}

	return tempDir, cleanup
}

func TestAutoImageChildren(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := model.NewMockStorageUploader(ctrl)
	mockTagReaderWriter := model.NewMockTagReaderWriter(ctrl)
	mockTagImageTypeReaderWriter := model.NewMockTagImageTypeReaderWriter(ctrl)

	t.Run("successful auto image processing", func(t *testing.T) {
		tempDir, cleanup := setupTempDirWithImages(t)
		defer cleanup()

		tag := &model.Tag{
			Id:    1,
			Title: "Parent Tag",
			Children: []*model.Tag{
				{Id: 2, Title: "Child Tag 1"},
			},
		}

		// Mock getting child tags
		childTag1 := &model.Tag{Id: 2, Title: "Child Tag 1", Images: []*model.TagImage{}}

		mockTagReaderWriter.EXPECT().
			GetTag(uint64(2)).
			Return(childTag1, nil)

		// The function calls autoImageTagType for each directory, and each call can potentially
		// call getOrCreateTagImageType, updateTagImageTypeIcon, and autoImageTag
		// Since this is complex integration, we'll use AnyTimes() for flexibility
		mockTagImageTypeReaderWriter.EXPECT().
			GetTagImageType(gomock.Any(), gomock.Any()).
			Return(nil, gorm.ErrRecordNotFound).
			AnyTimes()
		mockTagImageTypeReaderWriter.EXPECT().
			CreateOrUpdateTagImageType(gomock.Any()).
			Return(nil).
			AnyTimes()

		// Mock storage operations
		mockStorage.EXPECT().
			GetFileForWriting(gomock.Any()).
			Return(filepath.Join(tempDir, "storage-file"), nil).
			AnyTimes()
		mockStorage.EXPECT().
			GetStorageUrl(gomock.Any()).
			Return("http://storage/test-url").
			AnyTimes()

		// Mock tag updates
		mockTagReaderWriter.EXPECT().
			CreateOrUpdateTag(gomock.Any()).
			Return(nil).
			AnyTimes()

		err := AutoImageChildren(mockStorage, mockTagReaderWriter, mockTagImageTypeReaderWriter, tag, tempDir)

		assert.NoError(t, err)
	})

	t.Run("error reading directory", func(t *testing.T) {
		tag := &model.Tag{
			Id:       1,
			Title:    "Parent Tag",
			Children: []*model.Tag{},
		}

		err := AutoImageChildren(mockStorage, mockTagReaderWriter, mockTagImageTypeReaderWriter, tag, "/non-existent-directory")

		assert.Error(t, err)
	})

	t.Run("empty children list", func(t *testing.T) {
		tempDir, cleanup := setupTempDirWithImages(t)
		defer cleanup()

		tag := &model.Tag{
			Id:       1,
			Title:    "Parent Tag",
			Children: []*model.Tag{}, // No children
		}

		err := AutoImageChildren(mockStorage, mockTagReaderWriter, mockTagImageTypeReaderWriter, tag, tempDir)

		assert.NoError(t, err)
	})
}

func TestAutoImageTagType(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := model.NewMockStorageUploader(ctrl)
	mockTagWriter := model.NewMockTagWriter(ctrl)
	mockTagImageTypeReaderWriter := model.NewMockTagImageTypeReaderWriter(ctrl)

	t.Run("successful auto image tag type", func(t *testing.T) {
		tempDir, cleanup := setupTempDirWithImages(t)
		defer cleanup()

		// Create a specific test file for the tag
		tagImageFile := filepath.Join(tempDir, "banners", "test-tag.jpg")
		f, _ := os.Create(tagImageFile)
		f.WriteString("tag image content")
		f.Close()

		tag := &model.Tag{
			Id:     1,
			Title:  "test-tag",
			Images: []*model.TagImage{},
		}

		tit := &model.TagImageType{
			Id:       123,
			Nickname: "banners",
			IconUrl:  "",
		}

		mockTagImageTypeReaderWriter.EXPECT().
			GetTagImageType("nickname = ?", "banners").
			Return(tit, nil)

		// No icon update needed since no icon file, so no storage calls for icon
		// Only storage calls for tag image
		mockStorage.EXPECT().
			GetFileForWriting(gomock.Any()).
			Return(filepath.Join(tempDir, "storage-file"), nil)

		mockStorage.EXPECT().
			GetStorageUrl(gomock.Any()).
			Return("http://storage/test-url")

		mockTagWriter.EXPECT().
			CreateOrUpdateTag(tag).
			Return(nil)

		err := autoImageTagType(mockStorage, mockTagWriter, mockTagImageTypeReaderWriter, tag, filepath.Join(tempDir, "banners"), "banners")

		assert.NoError(t, err)
		assert.Len(t, tag.Images, 1)
		assert.Equal(t, "http://storage/test-url", tag.Images[0].Url)
	})

	t.Run("error getting tag image type", func(t *testing.T) {
		tempDir, cleanup := setupTempDirWithImages(t)
		defer cleanup()

		tag := &model.Tag{Id: 1, Title: "test-tag"}
		expectedError := errors.New("database error")

		mockTagImageTypeReaderWriter.EXPECT().
			GetTagImageType("nickname = ?", "banners").
			Return(nil, expectedError)

		err := autoImageTagType(mockStorage, mockTagWriter, mockTagImageTypeReaderWriter, tag, filepath.Join(tempDir, "banners"), "banners")

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})
}

func TestUpdateTagImageTypeIcon(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := model.NewMockStorageUploader(ctrl)
	mockTagImageTypeReaderWriter := model.NewMockTagImageTypeReaderWriter(ctrl)

	t.Run("icon already exists", func(t *testing.T) {
		tit := &model.TagImageType{
			Id:      123,
			IconUrl: "http://existing-icon.jpg",
		}

		err := updateTagImageTypeIcon(mockStorage, mockTagImageTypeReaderWriter, tit, "/some/directory")

		assert.NoError(t, err)
	})

	t.Run("successful icon update", func(t *testing.T) {
		tempDir, cleanup := setupTempDirWithImages(t)
		defer cleanup()

		// Create icon file
		iconFile := filepath.Join(tempDir, "icon.jpg")
		f, _ := os.Create(iconFile)
		f.WriteString("icon content")
		f.Close()

		tit := &model.TagImageType{
			Id:      123,
			IconUrl: "",
		}

		mockStorage.EXPECT().
			GetFileForWriting(gomock.Any()).
			Return(filepath.Join(tempDir, "storage-icon"), nil)

		mockStorage.EXPECT().
			GetStorageUrl(gomock.Any()).
			Return("http://storage/icon-url")

		mockTagImageTypeReaderWriter.EXPECT().
			CreateOrUpdateTagImageType(tit).
			Return(nil)

		err := updateTagImageTypeIcon(mockStorage, mockTagImageTypeReaderWriter, tit, tempDir)

		assert.NoError(t, err)
		assert.Equal(t, "http://storage/icon-url", tit.IconUrl)
	})

	t.Run("no icon file found", func(t *testing.T) {
		tempDir, cleanup := setupTempDirWithImages(t)
		defer cleanup()

		tit := &model.TagImageType{
			Id:      123,
			IconUrl: "",
		}

		err := updateTagImageTypeIcon(mockStorage, mockTagImageTypeReaderWriter, tit, tempDir)

		assert.NoError(t, err)
		assert.Equal(t, "", tit.IconUrl) // Should remain empty
	})

	t.Run("error getting storage file", func(t *testing.T) {
		tempDir, cleanup := setupTempDirWithImages(t)
		defer cleanup()

		// Create icon file
		iconFile := filepath.Join(tempDir, "icon.jpg")
		f, _ := os.Create(iconFile)
		f.WriteString("icon content")
		f.Close()

		tit := &model.TagImageType{
			Id:      123,
			IconUrl: "",
		}

		expectedError := errors.New("storage error")
		mockStorage.EXPECT().
			GetFileForWriting(gomock.Any()).
			Return("", expectedError)

		err := updateTagImageTypeIcon(mockStorage, mockTagImageTypeReaderWriter, tit, tempDir)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})
}

func TestImageExists(t *testing.T) {
	t.Run("image exists with valid URL", func(t *testing.T) {
		tag := &model.Tag{
			Images: []*model.TagImage{
				{Id: 1, ImageTypeId: 123, Url: "http://example.com/image.jpg"},
				{Id: 2, ImageTypeId: 456, Url: "http://example.com/image2.jpg"},
			},
		}
		tit := &model.TagImageType{Id: 123}

		result := imageExists(tag, tit)
		assert.True(t, result)
	})

	t.Run("image exists but URL is empty", func(t *testing.T) {
		tag := &model.Tag{
			Images: []*model.TagImage{
				{Id: 1, ImageTypeId: 123, Url: ""},
			},
		}
		tit := &model.TagImageType{Id: 123}

		result := imageExists(tag, tit)
		assert.False(t, result)
	})

	t.Run("image exists but URL is 'none'", func(t *testing.T) {
		tag := &model.Tag{
			Images: []*model.TagImage{
				{Id: 1, ImageTypeId: 123, Url: "none"},
			},
		}
		tit := &model.TagImageType{Id: 123}

		result := imageExists(tag, tit)
		assert.False(t, result)
	})

	t.Run("image does not exist for image type", func(t *testing.T) {
		tag := &model.Tag{
			Images: []*model.TagImage{
				{Id: 1, ImageTypeId: 456, Url: "http://example.com/image.jpg"},
			},
		}
		tit := &model.TagImageType{Id: 123}

		result := imageExists(tag, tit)
		assert.False(t, result)
	})

	t.Run("tag has no images", func(t *testing.T) {
		tag := &model.Tag{
			Images: []*model.TagImage{},
		}
		tit := &model.TagImageType{Id: 123}

		result := imageExists(tag, tit)
		assert.False(t, result)
	})
}

func TestGetOrCreateTagImageType(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTagImageTypeReaderWriter := model.NewMockTagImageTypeReaderWriter(ctrl)

	t.Run("tag image type exists", func(t *testing.T) {
		nickname := "banners"
		existingTit := &model.TagImageType{
			Id:       123,
			Nickname: nickname,
			IconUrl:  "http://example.com/icon.jpg",
		}

		mockTagImageTypeReaderWriter.EXPECT().
			GetTagImageType("nickname = ?", nickname).
			Return(existingTit, nil)

		result, err := getOrCreateTagImageType(mockTagImageTypeReaderWriter, nickname)

		assert.NoError(t, err)
		assert.Equal(t, existingTit, result)
	})

	t.Run("tag image type does not exist - create new", func(t *testing.T) {
		nickname := "thumbnails"

		mockTagImageTypeReaderWriter.EXPECT().
			GetTagImageType("nickname = ?", nickname).
			Return(nil, gorm.ErrRecordNotFound)
		mockTagImageTypeReaderWriter.EXPECT().
			CreateOrUpdateTagImageType(gomock.Any()).
			DoAndReturn(func(tit *model.TagImageType) error {
				tit.Id = 456 // Simulate DB assigning ID
				return nil
			})

		result, err := getOrCreateTagImageType(mockTagImageTypeReaderWriter, nickname)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, nickname, result.Nickname)
		assert.Equal(t, uint64(456), result.Id)
	})

	t.Run("error getting tag image type (not record not found)", func(t *testing.T) {
		nickname := "banners"
		expectedError := errors.New("database connection error")

		mockTagImageTypeReaderWriter.EXPECT().
			GetTagImageType("nickname = ?", nickname).
			Return(nil, expectedError)

		result, err := getOrCreateTagImageType(mockTagImageTypeReaderWriter, nickname)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedError, err)
	})

	t.Run("error creating tag image type", func(t *testing.T) {
		nickname := "thumbnails"
		expectedError := errors.New("create error")

		mockTagImageTypeReaderWriter.EXPECT().
			GetTagImageType("nickname = ?", nickname).
			Return(nil, gorm.ErrRecordNotFound)
		mockTagImageTypeReaderWriter.EXPECT().
			CreateOrUpdateTagImageType(gomock.Any()).
			Return(expectedError)

		result, err := getOrCreateTagImageType(mockTagImageTypeReaderWriter, nickname)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedError, err)
	})
}

func TestAutoImageTag(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := model.NewMockStorageUploader(ctrl)
	mockTagWriter := model.NewMockTagWriter(ctrl)

	t.Run("successful auto image tag", func(t *testing.T) {
		tempDir, cleanup := setupTempDirWithImages(t)
		defer cleanup()

		// Create a specific test file for the tag
		tagImageFile := filepath.Join(tempDir, "test-tag.jpg")
		f, _ := os.Create(tagImageFile)
		f.WriteString("tag image content")
		f.Close()

		tag := &model.Tag{
			Id:     1,
			Title:  "test-tag",
			Images: []*model.TagImage{},
		}

		tit := &model.TagImageType{
			Id:       123,
			Nickname: "banners",
		}

		mockStorage.EXPECT().
			GetFileForWriting(gomock.Any()).
			Return(filepath.Join(tempDir, "storage-file"), nil)

		mockStorage.EXPECT().
			GetStorageUrl(gomock.Any()).
			Return("http://storage/tag-image-url")

		mockTagWriter.EXPECT().
			CreateOrUpdateTag(tag).
			Return(nil)

		err := autoImageTag(mockStorage, mockTagWriter, tag, tempDir, tit)

		assert.NoError(t, err)
		assert.Len(t, tag.Images, 1)
		assert.Equal(t, "http://storage/tag-image-url", tag.Images[0].Url)
		assert.Equal(t, tag.Id, tag.Images[0].TagId)
		assert.Equal(t, tit.Id, tag.Images[0].ImageTypeId)
		assert.NotZero(t, tag.Images[0].ImageNonce)
	})

	t.Run("no image file found", func(t *testing.T) {
		tempDir, cleanup := setupTempDirWithImages(t)
		defer cleanup()

		tag := &model.Tag{
			Id:     1,
			Title:  "non-existent-tag",
			Images: []*model.TagImage{},
		}

		tit := &model.TagImageType{
			Id:       123,
			Nickname: "banners",
		}

		err := autoImageTag(mockStorage, mockTagWriter, tag, tempDir, tit)

		assert.NoError(t, err)
		assert.Len(t, tag.Images, 0) // No image should be added
	})

	t.Run("image already exists for tag and image type", func(t *testing.T) {
		tempDir, cleanup := setupTempDirWithImages(t)
		defer cleanup()

		// Create a specific test file for the tag
		tagImageFile := filepath.Join(tempDir, "test-tag.jpg")
		f, _ := os.Create(tagImageFile)
		f.WriteString("tag image content")
		f.Close()

		tag := &model.Tag{
			Id:    1,
			Title: "test-tag",
			Images: []*model.TagImage{
				{Id: 1, TagId: 1, ImageTypeId: 123, Url: "http://existing.jpg"},
			},
		}

		tit := &model.TagImageType{
			Id:       123,
			Nickname: "banners",
		}

		mockStorage.EXPECT().
			GetFileForWriting(gomock.Any()).
			Return(filepath.Join(tempDir, "storage-file"), nil)

		mockTagWriter.EXPECT().
			CreateOrUpdateTag(tag).
			Return(nil)

		err := autoImageTag(mockStorage, mockTagWriter, tag, tempDir, tit)

		assert.NoError(t, err)
		assert.Len(t, tag.Images, 1) // Should still be only 1 image
	})

	t.Run("error getting storage file", func(t *testing.T) {
		tempDir, cleanup := setupTempDirWithImages(t)
		defer cleanup()

		// Create a specific test file for the tag
		tagImageFile := filepath.Join(tempDir, "test-tag.jpg")
		f, _ := os.Create(tagImageFile)
		f.WriteString("tag image content")
		f.Close()

		tag := &model.Tag{
			Id:     1,
			Title:  "test-tag",
			Images: []*model.TagImage{},
		}

		tit := &model.TagImageType{
			Id:       123,
			Nickname: "banners",
		}

		expectedError := errors.New("storage error")
		mockStorage.EXPECT().
			GetFileForWriting(gomock.Any()).
			Return("", expectedError)

		err := autoImageTag(mockStorage, mockTagWriter, tag, tempDir, tit)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})

	t.Run("error updating tag", func(t *testing.T) {
		tempDir, cleanup := setupTempDirWithImages(t)
		defer cleanup()

		// Create a specific test file for the tag
		tagImageFile := filepath.Join(tempDir, "test-tag.jpg")
		f, _ := os.Create(tagImageFile)
		f.WriteString("tag image content")
		f.Close()

		tag := &model.Tag{
			Id:     1,
			Title:  "test-tag",
			Images: []*model.TagImage{},
		}

		tit := &model.TagImageType{
			Id:       123,
			Nickname: "banners",
		}

		expectedError := errors.New("update error")
		mockStorage.EXPECT().
			GetFileForWriting(gomock.Any()).
			Return(filepath.Join(tempDir, "storage-file"), nil)

		mockStorage.EXPECT().
			GetStorageUrl(gomock.Any()).
			Return("http://storage/tag-image-url")

		mockTagWriter.EXPECT().
			CreateOrUpdateTag(tag).
			Return(expectedError)

		err := autoImageTag(mockStorage, mockTagWriter, tag, tempDir, tit)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})
}

func TestFindExistingImage(t *testing.T) {
	t.Run("find image with exact tag title", func(t *testing.T) {
		tempDir, cleanup := setupTempDirWithImages(t)
		defer cleanup()

		// Create a test image file
		testFile := filepath.Join(tempDir, "My Test Tag.jpg")
		f, _ := os.Create(testFile)
		f.Close()

		result, err := findExistingImage("My Test Tag", tempDir)

		assert.NoError(t, err)
		assert.Equal(t, testFile, result)
	})

	t.Run("find image with directory-formatted tag title", func(t *testing.T) {
		tempDir, cleanup := setupTempDirWithImages(t)
		defer cleanup()

		// Create a test image file with directory format
		testFile := filepath.Join(tempDir, "my-test-tag.png")
		f, _ := os.Create(testFile)
		f.Close()

		result, err := findExistingImage("My Test Tag", tempDir)

		assert.NoError(t, err)
		assert.Equal(t, testFile, result)
	})

	t.Run("find image with different extension", func(t *testing.T) {
		tempDir, cleanup := setupTempDirWithImages(t)
		defer cleanup()

		// Create test image files with different extensions
		testFile := filepath.Join(tempDir, "test-tag.svg")
		f, _ := os.Create(testFile)
		f.Close()

		result, err := findExistingImage("test-tag", tempDir)

		assert.NoError(t, err)
		assert.Equal(t, testFile, result)
	})

	t.Run("no image found", func(t *testing.T) {
		tempDir, cleanup := setupTempDirWithImages(t)
		defer cleanup()

		result, err := findExistingImage("non-existent-tag", tempDir)

		assert.NoError(t, err)
		assert.Equal(t, "", result)
	})

	t.Run("multiple extensions available - returns first found", func(t *testing.T) {
		tempDir, cleanup := setupTempDirWithImages(t)
		defer cleanup()

		// Create multiple files with different extensions
		jpgFile := filepath.Join(tempDir, "test-tag.jpg")
		pngFile := filepath.Join(tempDir, "test-tag.png")

		f1, _ := os.Create(jpgFile)
		f1.Close()
		f2, _ := os.Create(pngFile)
		f2.Close()

		result, err := findExistingImage("test-tag", tempDir)

		assert.NoError(t, err)
		// Should return jpg since it's checked first in the extensions list
		assert.Equal(t, jpgFile, result)
	})

	t.Run("test directories.TagTitleToDirectory functionality", func(t *testing.T) {
		// Test the integration with directories.TagTitleToDirectory
		tempDir, cleanup := setupTempDirWithImages(t)
		defer cleanup()

		// Create a file that matches the directory format
		dirFormat := directories.TagTitleToDirectory("My Test Tag") // Should be "my-test-tag"
		testFile := filepath.Join(tempDir, dirFormat+".jpg")
		f, _ := os.Create(testFile)
		f.Close()

		result, err := findExistingImage("My Test Tag", tempDir)

		assert.NoError(t, err)
		assert.Equal(t, testFile, result)
	})
}
