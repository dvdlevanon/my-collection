package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"my-collection/server/pkg/db"
	"my-collection/server/pkg/gallery"
	"my-collection/server/pkg/model"
	"my-collection/server/pkg/storage"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupNewServer(t *testing.T, filename string) *Server {
	assert.NoError(t, os.MkdirAll(".tests", 0750))
	dbpath := fmt.Sprintf(".tests/%s", filename)
	_, err := os.Create(dbpath)
	assert.NoError(t, err)
	assert.NoError(t, os.Remove(dbpath))
	db, err := db.New(dbpath)
	assert.NoError(t, err)
	storage, err := storage.New("/tmp/root-directory/.storage")
	assert.NoError(t, err)
	gallery := gallery.New(db, storage, "")
	return New(gallery, storage)
}

func TestCreateAndGetItem(t *testing.T) {
	server := setupNewServer(t, "create-item-test.sqlite")
	item := model.Item{Title: "title1", Url: "url1"}
	payload, err := json.Marshal(item)
	assert.NoError(t, err)
	req := httptest.NewRequest("POST", "/items", bytes.NewReader(payload))
	resp := httptest.NewRecorder()
	server.router.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, http.StatusOK)
	returnedId := model.Item{}
	err = json.Unmarshal(resp.Body.Bytes(), &returnedId)
	assert.NoError(t, err)
	assert.Equal(t, returnedId.Id, uint64(1))
	req = httptest.NewRequest("GET", fmt.Sprintf("/items/%d", returnedId.Id), nil)
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
	req := httptest.NewRequest("POST", "/tags", bytes.NewReader(payload))
	resp := httptest.NewRecorder()
	server.router.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, http.StatusOK)
	returnedId := model.Item{}
	err = json.Unmarshal(resp.Body.Bytes(), &returnedId)
	assert.NoError(t, err)
	assert.Equal(t, returnedId.Id, uint64(1))
	req = httptest.NewRequest("GET", fmt.Sprintf("/tags/%d", returnedId.Id), nil)
	resp = httptest.NewRecorder()
	server.router.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, http.StatusOK)
	returnedTag := model.Item{}
	err = json.Unmarshal(resp.Body.Bytes(), &returnedTag)
	assert.NoError(t, err)
	assert.Equal(t, returnedTag.Id, returnedId.Id)
	assert.Equal(t, returnedTag.Title, tag.Title)
}

func TestUpdateItem(t *testing.T) {
	server := setupNewServer(t, "update-item-test.sqlite")
	item := model.Item{Title: "title1"}
	payload, err := json.Marshal(item)
	assert.NoError(t, err)
	req := httptest.NewRequest("POST", "/items", bytes.NewReader(payload))
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
	req = httptest.NewRequest("POST", fmt.Sprintf("/items/%d", returnedId.Id), bytes.NewReader(payload))
	resp = httptest.NewRecorder()
	server.router.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, http.StatusOK)
	req = httptest.NewRequest("GET", fmt.Sprintf("/items/%d", returnedId.Id), nil)
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
	req := httptest.NewRequest("POST", "/tags", bytes.NewReader(payload))
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
	req = httptest.NewRequest("POST", fmt.Sprintf("/tags/%d", returnedId.Id), bytes.NewReader(payload))
	resp = httptest.NewRecorder()
	server.router.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, http.StatusOK)
	req = httptest.NewRequest("GET", fmt.Sprintf("/tags/%d", returnedId.Id), nil)
	resp = httptest.NewRecorder()
	server.router.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, http.StatusOK)
	returnedTag := model.Tag{}
	err = json.Unmarshal(resp.Body.Bytes(), &returnedTag)
	assert.NoError(t, err)
	assert.Equal(t, returnedTag.Id, returnedId.Id)
	assert.Equal(t, returnedTag.Title, "updated-title")
}
