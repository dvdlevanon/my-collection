package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"my-collection/server/pkg/model"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-errors/errors"
	"github.com/google/uuid"
	"k8s.io/utils/pointer"
)

func (s *Server) createItem(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if s.handleError(c, err) {
		return
	}

	var item model.Item
	if err = json.Unmarshal(body, &item); err != nil {
		s.handleError(c, err)
		return
	}

	if err = s.gallery.CreateOrUpdateItem(&item); err != nil {
		s.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, model.Item{Id: item.Id})
}

func (s *Server) createTag(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if s.handleError(c, err) {
		return
	}

	var tag model.Tag
	if err = json.Unmarshal(body, &tag); err != nil {
		s.handleError(c, err)
		return
	}

	if err = s.gallery.CreateOrUpdateTag(&tag); err != nil {
		s.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, model.Tag{Id: tag.Id})
}

func (s *Server) updateItem(c *gin.Context) {
	itemId, err := strconv.ParseUint(c.Param("item"), 10, 64)
	if s.handleError(c, err) {
		return
	}

	body, err := io.ReadAll(c.Request.Body)
	if s.handleError(c, err) {
		return
	}

	var item model.Item
	if err = json.Unmarshal(body, &item); err != nil {
		s.handleError(c, err)
		return
	}

	if item.Id != 0 && item.Id != itemId {
		s.handleError(c, errors.Errorf("Mismatch IDs %d != %d", item.Id, itemId))
		return
	}

	item.Id = itemId
	if err = s.gallery.UpdateItem(&item); err != nil {
		s.handleError(c, err)
		return
	}

	c.Status(http.StatusOK)
}

func (s *Server) removeTagFromItem(c *gin.Context) {
	itemId, err := strconv.ParseUint(c.Param("item"), 10, 64)
	if s.handleError(c, err) {
		return
	}

	tagId, err := strconv.ParseUint(c.Param("tag"), 10, 64)
	if s.handleError(c, err) {
		return
	}

	if s.handleError(c, s.gallery.RemoveTagFromItem(itemId, tagId)) {
		return
	}

	c.Status(http.StatusOK)
}

func (s *Server) updateTag(c *gin.Context) {
	tagId, err := strconv.ParseUint(c.Param("tag"), 10, 64)
	if s.handleError(c, err) {
		return
	}

	body, err := io.ReadAll(c.Request.Body)
	if s.handleError(c, err) {
		return
	}

	var tag model.Tag
	if err = json.Unmarshal(body, &tag); err != nil {
		s.handleError(c, err)
		return
	}

	if tag.Id != 0 && tag.Id != tagId {
		s.handleError(c, errors.Errorf("Mismatch IDs %d != %d", tag.Id, tagId))
		return
	}

	tag.Id = tagId
	if err = s.gallery.UpdateTag(&tag); err != nil {
		s.handleError(c, err)
		return
	}

	c.Status(http.StatusOK)
}

func (s *Server) getItem(c *gin.Context) {
	itemId, err := strconv.ParseUint(c.Param("item"), 10, 64)
	if s.handleError(c, err) {
		return
	}

	item, err := s.gallery.GetItem(itemId)
	if s.handleError(c, err) {
		return
	}

	c.JSON(http.StatusOK, item)
}

func (s *Server) getTag(c *gin.Context) {
	tagId, err := strconv.ParseUint(c.Param("tag"), 10, 64)
	if s.handleError(c, err) {
		return
	}

	tag, err := s.gallery.GetTag(tagId)
	if s.handleError(c, err) {
		return
	}

	c.JSON(http.StatusOK, tag)
}

func (s *Server) autoImage(c *gin.Context) {
	tagId, err := strconv.ParseUint(c.Param("tag"), 10, 64)
	if s.handleError(c, err) {
		return
	}

	tag, err := s.gallery.GetTag(tagId)
	if s.handleError(c, err) {
		return
	}

	body, err := io.ReadAll(c.Request.Body)
	if s.handleError(c, err) {
		return
	}

	var fileUrl model.FileUrl
	if err = json.Unmarshal(body, &fileUrl); err != nil {
		s.handleError(c, err)
		return
	}

	if s.handleError(c, s.gallery.AutoImageChildren(tag, fileUrl.Url)) {
		return
	}

	c.JSON(http.StatusOK, tag)
}

func (s *Server) removeTag(c *gin.Context) {
	tagId, err := strconv.ParseUint(c.Param("tag"), 10, 64)
	if s.handleError(c, err) {
		return
	}

	err = s.gallery.RemoveTag(tagId)
	if s.handleError(c, err) {
		return
	}

	c.Status(http.StatusOK)
}

func (s *Server) getTags(c *gin.Context) {
	tags, err := s.gallery.GetAllTags()

	if s.handleError(c, err) {
		return
	}

	logger.Infof("Get tags return %d tags", len(*tags))
	c.JSON(http.StatusOK, tags)
}

func (s *Server) getItems(c *gin.Context) {
	items, err := s.gallery.GetAllItems()

	if s.handleError(c, err) {
		return
	}

	logger.Infof("Get items return %d items", len(*items))
	c.JSON(http.StatusOK, items)
}

func (s *Server) refreshItemsCovers(c *gin.Context) {
	err := s.processor.EnqueueAllItemsCovers()

	if s.handleError(c, err) {
		return
	}

	c.Status(http.StatusOK)
}

func (s *Server) refreshItemsPreview(c *gin.Context) {
	err := s.processor.EnqueueAllItemsPreview()

	if s.handleError(c, err) {
		return
	}

	c.Status(http.StatusOK)
}

func (s *Server) refreshItemsVideoMetadata(c *gin.Context) {
	err := s.processor.EnqueueAllItemsVideoMetadata()

	if s.handleError(c, err) {
		return
	}

	c.Status(http.StatusOK)
}

func (s *Server) uploadFile(c *gin.Context) {
	form, err := c.MultipartForm()
	if s.handleError(c, err) {
		return
	}

	path := form.Value["path"][0]
	file := form.File["file"][0]
	fileName := fmt.Sprintf("%s-%s", file.Filename, uuid.NewString())
	relativeFile := filepath.Join(path, fileName)
	storageFile, err := s.storage.GetFileForWriting(relativeFile)
	if s.handleError(c, err) {
		return
	}

	if s.handleError(c, c.SaveUploadedFile(file, storageFile)) {
		return
	}

	c.JSON(http.StatusOK, model.FileUrl{Url: s.storage.GetStorageUrl(relativeFile)})
}

func (s *Server) getFile(c *gin.Context) {
	path := c.Param("path")[1:]
	var file string
	if s.storage.IsStorageUrl(path) {
		file = s.storage.GetFile(path)
	} else {
		file = s.gallery.GetFile(path)
	}

	logger.Infof("Getting file %v", file)
	http.ServeFile(c.Writer, c.Request, file)
}

func (s *Server) exportMetadata(c *gin.Context) {
	jsonBytes := bytes.Buffer{}
	if s.handleError(c, s.gallery.Export(&jsonBytes)) {
		return
	}

	c.Header("Content-Type", "application/json")
	c.Header("Content-Disposition", "gallery-metadata.json")
	c.String(http.StatusOK, jsonBytes.String())
}

func (s *Server) addAnnotationToTag(c *gin.Context) {
	tagId, err := strconv.ParseUint(c.Param("tag"), 10, 64)
	if s.handleError(c, err) {
		return
	}

	body, err := io.ReadAll(c.Request.Body)
	if s.handleError(c, err) {
		return
	}

	var annotation model.TagAnnotation
	if err = json.Unmarshal(body, &annotation); err != nil {
		s.handleError(c, err)
		return
	}

	annotationId, err := s.gallery.AddAnnotationToTag(tagId, annotation)
	if s.handleError(c, err) {
		return
	}

	c.JSON(http.StatusOK, model.TagAnnotation{Id: annotationId})
}

func (s *Server) removeAnnotationFromTag(c *gin.Context) {
	tagId, err := strconv.ParseUint(c.Param("tag"), 10, 64)
	if s.handleError(c, err) {
		return
	}

	annotationId, err := strconv.ParseUint(c.Param("annotation-id"), 10, 64)
	if s.handleError(c, err) {
		return
	}

	if s.handleError(c, s.gallery.RemoveTagAnnotationFromTag(tagId, annotationId)) {
		return
	}

	c.Status(http.StatusOK)
}

func (s *Server) getTagAvailableAnnotations(c *gin.Context) {
	tagId, err := strconv.ParseUint(c.Param("tag"), 10, 64)
	if s.handleError(c, err) {
		return
	}

	availableAnnotations, err := s.gallery.GetTagAvailableAnnotations(tagId)
	if s.handleError(c, err) {
		return
	}

	c.JSON(http.StatusOK, availableAnnotations)
}

func (s *Server) getDirectories(c *gin.Context) {
	directories, err := s.gallery.GetAllDirectories()

	if s.handleError(c, err) {
		return
	}

	logger.Infof("Get directories return %d directories", len(*directories))
	c.JSON(http.StatusOK, directories)
}

func (s *Server) createOrUpdateDirectory(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if s.handleError(c, err) {
		return
	}

	var directory model.Directory
	if err = json.Unmarshal(body, &directory); err != nil {
		s.handleError(c, err)
		return
	}

	directory.ProcessingStart = pointer.Int64(time.Now().UnixMilli())
	if err = s.gallery.CreateOrUpdateDirectory(&directory); err != nil {
		s.handleError(c, err)
		return
	}

	s.directories.DirectoryChanged(&directory)
	c.Status(http.StatusOK)
}

func (s *Server) SetDirectoryTags(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if s.handleError(c, err) {
		return
	}

	var directory model.Directory
	if err = json.Unmarshal(body, &directory); err != nil {
		s.handleError(c, err)
		return
	}

	if err = s.gallery.SetDirectoryTags(&directory); err != nil {
		s.handleError(c, err)
		return
	}

	s.directories.DirectoryChanged(&directory)
	c.Status(http.StatusOK)
}

func (s *Server) getDirectory(c *gin.Context) {
	directoryPath := c.Param("directory")[1:]
	directory, err := s.gallery.GetDirectory("path = ?", directoryPath)
	if s.handleError(c, err) {
		return
	}

	c.JSON(http.StatusOK, directory)
}

func (s *Server) excludeDirectory(c *gin.Context) {
	directoryPath := c.Param("directory")[1:]

	err := s.gallery.ExcludeDirectory(directoryPath)
	if s.handleError(c, err) {
		return
	}

	s.directories.DirectoryExcluded(directoryPath)
	c.Status(http.StatusOK)
}
