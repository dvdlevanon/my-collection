package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"my-collection/server/pkg/db"
	"my-collection/server/pkg/gallery"
	itemprocessor "my-collection/server/pkg/item-processor"
	"my-collection/server/pkg/model"
	"my-collection/server/pkg/storage"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/utils/pointer"
)

type DirectoriesMock struct{}
type ItemProcessorMock struct{}

func (d *DirectoriesMock) Init() error                                                     { return nil }
func (d *DirectoriesMock) DirectoryChanged(directory *model.Directory)                     {}
func (d *DirectoriesMock) DirectoryExcluded(path string)                                   {}
func (d *ItemProcessorMock) Run()                                                          {}
func (d *ItemProcessorMock) EnqueueAllItemsCovers(force bool) error                        { return nil }
func (d *ItemProcessorMock) EnqueueAllItemsPreview(force bool) error                       { return nil }
func (d *ItemProcessorMock) EnqueueAllItemsVideoMetadata(force bool) error                 { return nil }
func (d *ItemProcessorMock) EnqueueItemVideoMetadata(id uint64)                            {}
func (d *ItemProcessorMock) EnqueueItemPreview(id uint64)                                  {}
func (d *ItemProcessorMock) EnqueueItemCovers(id uint64)                                   {}
func (d *ItemProcessorMock) EnqueueMainCover(id uint64, second float64)                    {}
func (d *ItemProcessorMock) IsPaused() bool                                                { return false }
func (d *ItemProcessorMock) Pause()                                                        {}
func (d *ItemProcessorMock) Continue()                                                     {}
func (d *ItemProcessorMock) SetProcessorNotifier(notifier itemprocessor.ProcessorNotifier) {}

func setupNewServer(t *testing.T, filename string) *Server {
	assert.NoError(t, os.MkdirAll(".tests", 0750))
	dbpath := fmt.Sprintf(".tests/%s", filename)
	_, err := os.Create(dbpath)
	assert.NoError(t, err)
	assert.NoError(t, os.Remove(dbpath))
	db, err := db.New("", dbpath)
	assert.NoError(t, err)
	storage, err := storage.New("/tmp/root-directory/.storage")
	assert.NoError(t, err)
	gallery := gallery.New(db, storage, "")
	return New(gallery, storage, &DirectoriesMock{}, &ItemProcessorMock{})
}

func TestCreateAndGetItem(t *testing.T) {
	server := setupNewServer(t, "create-item-test.sqlite")
	item := model.Item{Title: "title1", Url: "url1", Origin: "origin"}
	payload, err := json.Marshal(item)
	assert.NoError(t, err)
	req := httptest.NewRequest("POST", "/api/items", bytes.NewReader(payload))
	resp := httptest.NewRecorder()
	server.router.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, http.StatusOK)
	returnedId := model.Item{}
	err = json.Unmarshal(resp.Body.Bytes(), &returnedId)
	assert.NoError(t, err)
	assert.Equal(t, returnedId.Id, uint64(1))
	req = httptest.NewRequest("GET", fmt.Sprintf("/api/items/%d", returnedId.Id), nil)
	resp = httptest.NewRecorder()
	server.router.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, http.StatusOK)
	returnedItem := model.Item{}
	err = json.Unmarshal(resp.Body.Bytes(), &returnedItem)
	assert.NoError(t, err)
	assert.Equal(t, returnedItem.Id, returnedId.Id)
	assert.Equal(t, returnedItem.Title, item.Title)
	assert.Equal(t, returnedItem.Url, item.Url)
}

func TestCreateAndGetTag(t *testing.T) {
	server := setupNewServer(t, "create-tag-test.sqlite")
	tag := model.Tag{Title: "title1"}
	payload, err := json.Marshal(tag)
	assert.NoError(t, err)
	req := httptest.NewRequest("POST", "/api/tags", bytes.NewReader(payload))
	resp := httptest.NewRecorder()
	server.router.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, http.StatusOK)
	returnedId := model.Tag{}
	err = json.Unmarshal(resp.Body.Bytes(), &returnedId)
	assert.NoError(t, err)
	assert.Equal(t, returnedId.Id, uint64(1))
	req = httptest.NewRequest("GET", fmt.Sprintf("/api/tags/%d", returnedId.Id), nil)
	resp = httptest.NewRecorder()
	server.router.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, http.StatusOK)
	returnedTag := model.Tag{}
	err = json.Unmarshal(resp.Body.Bytes(), &returnedTag)
	assert.NoError(t, err)
	assert.Equal(t, returnedTag.Id, returnedId.Id)
	assert.Equal(t, returnedTag.Title, tag.Title)
}

func TestUpdateItem(t *testing.T) {
	server := setupNewServer(t, "update-item-test.sqlite")
	item := model.Item{Title: "title1", Origin: "origin"}
	payload, err := json.Marshal(item)
	assert.NoError(t, err)
	req := httptest.NewRequest("POST", "/api/items", bytes.NewReader(payload))
	resp := httptest.NewRecorder()
	server.router.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, http.StatusOK)
	returnedId := model.Item{}
	err = json.Unmarshal(resp.Body.Bytes(), &returnedId)
	assert.NoError(t, err)
	assert.Equal(t, returnedId.Id, uint64(1))
	item = model.Item{Id: returnedId.Id, Url: "update-url"}
	payload, err = json.Marshal(item)
	assert.NoError(t, err)
	req = httptest.NewRequest("POST", fmt.Sprintf("/api/items/%d", returnedId.Id), bytes.NewReader(payload))
	resp = httptest.NewRecorder()
	server.router.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, http.StatusOK)
	req = httptest.NewRequest("GET", fmt.Sprintf("/api/items/%d", returnedId.Id), nil)
	resp = httptest.NewRecorder()
	server.router.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, http.StatusOK)
	returnedItem := model.Item{}
	err = json.Unmarshal(resp.Body.Bytes(), &returnedItem)
	assert.NoError(t, err)
	assert.Equal(t, returnedItem.Id, returnedId.Id)
	assert.Equal(t, returnedItem.Title, "title1")
	assert.Equal(t, returnedItem.Url, "update-url")
}

func TestUpdateTag(t *testing.T) {
	server := setupNewServer(t, "update-tag-test.sqlite")
	tag := model.Tag{Title: "title1"}
	payload, err := json.Marshal(tag)
	assert.NoError(t, err)
	req := httptest.NewRequest("POST", "/api/tags", bytes.NewReader(payload))
	resp := httptest.NewRecorder()
	server.router.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, http.StatusOK)
	returnedId := model.Tag{}
	err = json.Unmarshal(resp.Body.Bytes(), &returnedId)
	assert.NoError(t, err)
	assert.Equal(t, returnedId.Id, uint64(1))
	tag = model.Tag{Id: returnedId.Id, Title: "updated-title"}
	payload, err = json.Marshal(tag)
	assert.NoError(t, err)
	req = httptest.NewRequest("POST", fmt.Sprintf("/api/tags/%d", returnedId.Id), bytes.NewReader(payload))
	resp = httptest.NewRecorder()
	server.router.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, http.StatusOK)
	req = httptest.NewRequest("GET", fmt.Sprintf("/api/tags/%d", returnedId.Id), nil)
	resp = httptest.NewRecorder()
	server.router.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, http.StatusOK)
	returnedTag := model.Tag{}
	err = json.Unmarshal(resp.Body.Bytes(), &returnedTag)
	assert.NoError(t, err)
	assert.Equal(t, returnedTag.Id, returnedId.Id)
	assert.Equal(t, returnedTag.Title, "updated-title")
}

func TestItemNotFound(t *testing.T) {
	server := setupNewServer(t, "item-not-found-test.sqlite")
	req := httptest.NewRequest("GET", fmt.Sprintf("/api/items/%d", 666), nil)
	resp := httptest.NewRecorder()
	server.router.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, http.StatusNotFound)
}

func TestTagNotFound(t *testing.T) {
	server := setupNewServer(t, "tag-not-found-test.sqlite")
	req := httptest.NewRequest("GET", fmt.Sprintf("/api/tags/%d", 666), nil)
	resp := httptest.NewRecorder()
	server.router.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, http.StatusNotFound)
}

func TestTagAnnotations(t *testing.T) {
	server := setupNewServer(t, "tag-annotations-test.sqlite")

	tag := model.Tag{Title: "title1"}
	payload, err := json.Marshal(tag)
	assert.NoError(t, err)
	req := httptest.NewRequest("POST", "/api/tags", bytes.NewReader(payload))
	resp := httptest.NewRecorder()
	server.router.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, http.StatusOK)
	returnedId := model.Tag{}
	err = json.Unmarshal(resp.Body.Bytes(), &returnedId)
	assert.NoError(t, err)

	annotation := model.TagAnnotation{Title: "annotation1"}
	payload, err = json.Marshal(annotation)
	assert.NoError(t, err)
	req = httptest.NewRequest("POST", fmt.Sprintf("/api/tags/%d/annotations", returnedId.Id), bytes.NewReader(payload))
	resp = httptest.NewRecorder()
	server.router.ServeHTTP(resp, req)
	returnedTagAnnotation := model.TagAnnotation{}
	err = json.Unmarshal(resp.Body.Bytes(), &returnedTagAnnotation)
	assert.NoError(t, err)
	assert.Equal(t, resp.Code, http.StatusOK)
	assert.Equal(t, uint64(1), returnedTagAnnotation.Id)

	req = httptest.NewRequest("GET", fmt.Sprintf("/api/tags/%d", returnedId.Id), nil)
	resp = httptest.NewRecorder()
	server.router.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, http.StatusOK)
	returnedTag := model.Tag{}
	err = json.Unmarshal(resp.Body.Bytes(), &returnedTag)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(returnedTag.Annotations))

	req = httptest.NewRequest("GET", fmt.Sprintf("/api/tags/%d/available-annotations", returnedId.Id), nil)
	resp = httptest.NewRecorder()
	server.router.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, http.StatusOK)
	returnedTagAnnotations := []model.Tag{}
	err = json.Unmarshal(resp.Body.Bytes(), &returnedTagAnnotations)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(returnedTagAnnotations))
	assert.Equal(t, "annotation1", returnedTagAnnotations[0].Title)

	req = httptest.NewRequest("DELETE", fmt.Sprintf("/api/tags/%d/annotations/%d", returnedId.Id, returnedTagAnnotations[0].Id), nil)
	resp = httptest.NewRecorder()
	server.router.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, http.StatusOK)

	req = httptest.NewRequest("GET", fmt.Sprintf("/api/tags/%d/available-annotations", returnedId.Id), nil)
	resp = httptest.NewRecorder()
	server.router.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, http.StatusOK)
	returnedTagAnnotations = []model.Tag{}
	err = json.Unmarshal(resp.Body.Bytes(), &returnedTagAnnotations)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(returnedTagAnnotations))
}

func TestDirectories(t *testing.T) {
	server := setupNewServer(t, "directories-test.sqlite")

	directory := model.Directory{
		Path:       "path/to/file",
		Excluded:   pointer.Bool(false),
		FilesCount: pointer.Int(10),
		Tags: []*model.Tag{
			{
				Title: "tag1",
			},
		},
	}

	payload, err := json.Marshal(directory)
	assert.NoError(t, err)
	req := httptest.NewRequest("POST", "/api/directories", bytes.NewReader(payload))
	resp := httptest.NewRecorder()
	server.router.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, http.StatusOK)

	req = httptest.NewRequest("GET", fmt.Sprintf("/api/directories/%s", directory.Path), nil)
	resp = httptest.NewRecorder()
	server.router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
	returnedDirectory := model.Directory{}
	err = json.Unmarshal(resp.Body.Bytes(), &returnedDirectory)
	assert.NoError(t, err)

	assert.Equal(t, directory.Path, returnedDirectory.Path)
	assert.Equal(t, directory.FilesCount, returnedDirectory.FilesCount)
	assert.Equal(t, *directory.Excluded, *returnedDirectory.Excluded)
	assert.Len(t, returnedDirectory.Tags, 1)
	assert.Empty(t, returnedDirectory.Tags[0].Title)
	assert.Equal(t, uint64(1), returnedDirectory.Tags[0].Id)
}
