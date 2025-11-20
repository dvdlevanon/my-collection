package general_tasks

import (
	"context"
	"encoding/json"
	"my-collection/server/pkg/model"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestMetadataDesc(t *testing.T) {
	result := MetadataDesc(123, "test.mp4")
	assert.Equal(t, "Metadata test.mp4", result)
}

func TestMarshalFileMetadataParams(t *testing.T) {
	params, err := MarshalFileMetadataParams(123)
	assert.NoError(t, err)

	var p fileMetadataParams
	err = json.Unmarshal([]byte(params), &p)
	assert.NoError(t, err)
	assert.Equal(t, uint64(123), p.ItemId)
}

func TestUnmarshalFileMetadataParams(t *testing.T) {
	p := fileMetadataParams{ItemId: 456}
	data, _ := json.Marshal(p)

	result, err := unmarshalFileMetadataParams(string(data))
	assert.NoError(t, err)
	assert.Equal(t, uint64(456), result.ItemId)
}

func TestUnmarshalFileMetadataParams_InvalidJSON(t *testing.T) {
	_, err := unmarshalFileMetadataParams("invalid json")
	assert.Error(t, err)
}

func TestUpdateFileMetadata(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIRW := model.NewMockItemReaderWriter(ctrl)
	ctx := context.Background()

	// Create a temporary file for testing
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.mp4")
	file, err := os.Create(tmpFile)
	assert.NoError(t, err)
	file.Close()

	item := &model.Item{
		Id:    123,
		Title: "test.mp4",
		Url:   tmpFile,
	}

	mockIRW.EXPECT().GetItem(ctx, uint64(123)).Return(item, nil)
	mockIRW.EXPECT().UpdateItem(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, item *model.Item) error {
		assert.Greater(t, item.LastModified, int64(0))
		assert.GreaterOrEqual(t, item.FileSize, int64(0))
		return nil
	})

	params, _ := MarshalFileMetadataParams(123)
	err = UpdateFileMetadata(ctx, mockIRW, params)
	assert.NoError(t, err)
}

func TestUpdateFileMetadata_GetItemError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIRW := model.NewMockItemReaderWriter(ctrl)
	ctx := context.Background()

	mockIRW.EXPECT().GetItem(ctx, uint64(123)).Return(nil, assert.AnError)

	params, _ := MarshalFileMetadataParams(123)
	err := UpdateFileMetadata(ctx, mockIRW, params)
	assert.Error(t, err)
}

func TestUpdateFileMetadata_FileNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIRW := model.NewMockItemReaderWriter(ctrl)
	ctx := context.Background()

	item := &model.Item{
		Id:    123,
		Title: "nonexistent.mp4",
		Url:   "/nonexistent/path/file.mp4",
	}

	mockIRW.EXPECT().GetItem(ctx, uint64(123)).Return(item, nil)

	params, _ := MarshalFileMetadataParams(123)
	err := UpdateFileMetadata(ctx, mockIRW, params)
	assert.Error(t, err)
}

func TestUpdateFileMetadata_UpdateItemError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIRW := model.NewMockItemReaderWriter(ctrl)
	ctx := context.Background()

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.mp4")
	file, err := os.Create(tmpFile)
	assert.NoError(t, err)
	file.Close()

	item := &model.Item{
		Id:    123,
		Title: "test.mp4",
		Url:   tmpFile,
	}

	mockIRW.EXPECT().GetItem(ctx, uint64(123)).Return(item, nil)
	mockIRW.EXPECT().UpdateItem(ctx, gomock.Any()).Return(assert.AnError)

	params, _ := MarshalFileMetadataParams(123)
	err = UpdateFileMetadata(ctx, mockIRW, params)
	assert.Error(t, err)
}

func TestUpdateFileMetadata_InvalidParams(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIRW := model.NewMockItemReaderWriter(ctrl)
	ctx := context.Background()

	err := UpdateFileMetadata(ctx, mockIRW, "invalid json")
	assert.Error(t, err)
}

func TestUpdateMetadata(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.mp4")
	file, err := os.Create(tmpFile)
	assert.NoError(t, err)
	file.WriteString("test content")
	file.Close()

	// Get file info before
	info, _ := os.Stat(tmpFile)
	expectedModTime := info.ModTime().UnixMilli()
	expectedSize := info.Size()

	// Wait a bit to ensure mod time might change
	time.Sleep(10 * time.Millisecond)

	item := &model.Item{
		Id:    123,
		Title: "test.mp4",
		Url:   tmpFile,
	}

	err = updateMetadata(item)
	assert.NoError(t, err)
	assert.Equal(t, expectedSize, item.FileSize)
	// ModTime should be set (might be slightly different due to timing)
	assert.GreaterOrEqual(t, item.LastModified, expectedModTime)
}

func TestUpdateMetadata_FileNotFound(t *testing.T) {
	item := &model.Item{
		Id:    123,
		Title: "nonexistent.mp4",
		Url:   "/nonexistent/path/file.mp4",
	}

	err := updateMetadata(item)
	assert.Error(t, err)
}
