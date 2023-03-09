package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"my-collection/server/pkg/backup"
	"my-collection/server/pkg/bl/directories"
	"my-collection/server/pkg/bl/tag_annotations"
	"my-collection/server/pkg/bl/tags"
	"my-collection/server/pkg/bl/tasks"
	"my-collection/server/pkg/model"
	"my-collection/server/pkg/relativasor"
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
	if s.handleError(c, json.Unmarshal(body, &item)) {
		return
	}

	item.Url = relativasor.GetRelativePath(item.Url)
	item.Origin = relativasor.GetRelativePath(item.Origin)
	if s.handleError(c, s.db.CreateOrUpdateItem(&item)) {
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
	if s.handleError(c, json.Unmarshal(body, &tag)) {
		return
	}

	if s.handleError(c, s.db.CreateOrUpdateTag(&tag)) {
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
	if s.handleError(c, json.Unmarshal(body, &item)) {
		return
	}

	if item.Id != 0 && item.Id != itemId {
		s.handleError(c, errors.Errorf("Mismatch IDs %d != %d", item.Id, itemId))
		return
	}

	item.Id = itemId
	if s.handleError(c, s.db.UpdateItem(&item)) {
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

	if s.handleError(c, s.db.RemoveTagFromItem(itemId, tagId)) {
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
	if s.handleError(c, json.Unmarshal(body, &tag)) {
		return
	}

	if tag.Id != 0 && tag.Id != tagId {
		s.handleError(c, errors.Errorf("Mismatch IDs %d != %d", tag.Id, tagId))
		return
	}

	tag.Id = tagId
	if s.handleError(c, s.db.UpdateTag(&tag)) {
		return
	}

	c.Status(http.StatusOK)
}

func (s *Server) getItem(c *gin.Context) {
	itemId, err := strconv.ParseUint(c.Param("item"), 10, 64)
	if s.handleError(c, err) {
		return
	}

	item, err := s.db.GetItem(itemId)
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

	tag, err := s.db.GetTag(tagId)
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

	tag, err := s.db.GetTag(tagId)
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

	if s.handleError(c, tags.AutoImageChildren(s.storage, s.db, tag, fileUrl.Url)) {
		return
	}

	c.JSON(http.StatusOK, tag)
}

func (s *Server) removeTag(c *gin.Context) {
	tagId, err := strconv.ParseUint(c.Param("tag"), 10, 64)
	if s.handleError(c, err) {
		return
	}

	if s.handleError(c, s.db.RemoveTag(tagId)) {
		return
	}

	c.Status(http.StatusOK)
}

func (s *Server) getTags(c *gin.Context) {
	tags, err := s.db.GetAllTags()
	if s.handleError(c, err) {
		return
	}

	logger.Infof("Get tags return %d tags", len(*tags))
	c.JSON(http.StatusOK, tags)
}

func (s *Server) getItems(c *gin.Context) {
	items, err := s.db.GetAllItems()
	if s.handleError(c, err) {
		return
	}

	logger.Infof("Get items return %d items", len(*items))
	c.JSON(http.StatusOK, items)
}

func (s *Server) refreshItemsCovers(c *gin.Context) {
	force, err := strconv.ParseBool(c.Query("force"))
	if err != nil {
		force = false
	}

	if s.handleError(c, s.processor.EnqueueAllItemsCovers(force)) {
		return
	}

	c.Status(http.StatusOK)
}

func (s *Server) refreshItemsPreview(c *gin.Context) {
	force, err := strconv.ParseBool(c.Query("force"))
	if err != nil {
		force = false
	}

	if s.handleError(c, s.processor.EnqueueAllItemsPreview(force)) {
		return
	}

	c.Status(http.StatusOK)
}

func (s *Server) refreshItemsVideoMetadata(c *gin.Context) {
	forceParam := c.Query("force")
	force, err := strconv.ParseBool(forceParam)
	if err != nil {
		force = false
	}

	if s.handleError(c, s.processor.EnqueueAllItemsVideoMetadata(force)) {
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
		file = relativasor.GetAbsoluteFile(path)
	}

	logger.Infof("Getting file %v", file)
	http.ServeFile(c.Writer, c.Request, file)
}

func (s *Server) exportMetadata(c *gin.Context) {
	jsonBytes := bytes.Buffer{}
	if s.handleError(c, backup.Export(s.db, s.db, &jsonBytes)) {
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

	annotationId, err := tag_annotations.AddAnnotationToTag(s.db, s.db, tagId, annotation)
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

	if s.handleError(c, s.db.RemoveTagAnnotationFromTag(tagId, annotationId)) {
		return
	}

	c.Status(http.StatusOK)
}

func (s *Server) getTagAvailableAnnotations(c *gin.Context) {
	tagId, err := strconv.ParseUint(c.Param("tag"), 10, 64)
	if s.handleError(c, err) {
		return
	}

	availableAnnotations, err := tag_annotations.GetTagAvailableAnnotations(s.db, s.db, tagId)
	if s.handleError(c, err) {
		return
	}

	c.JSON(http.StatusOK, availableAnnotations)
}

func (s *Server) getDirectories(c *gin.Context) {
	directories, err := s.db.GetAllDirectories()

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
	directory.Path = relativasor.GetRelativePath(directory.Path)
	if err = s.db.CreateOrUpdateDirectory(&directory); err != nil {
		s.handleError(c, err)
		return
	}

	s.fswatch.DirectoryChanged(&directory)
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

	if err = directories.SetDirectoryTags(s.db, &directory); err != nil {
		s.handleError(c, err)
		return
	}

	s.fswatch.DirectoryChanged(&directory)
	c.Status(http.StatusOK)
}

func (s *Server) getDirectory(c *gin.Context) {
	directoryPath := c.Param("directory")[1:]
	directory, err := s.db.GetDirectory("path = ?", directoryPath)
	if s.handleError(c, err) {
		return
	}

	c.JSON(http.StatusOK, directory)
}

func (s *Server) excludeDirectory(c *gin.Context) {
	directoryPath := c.Param("directory")[1:]

	err := directories.ExcludeDirectory(s.db, directoryPath)
	if s.handleError(c, err) {
		return
	}

	s.fswatch.DirectoryExcluded(directoryPath)
	c.Status(http.StatusOK)
}

func (s *Server) setMainCover(c *gin.Context) {
	itemId, err := strconv.ParseUint(c.Param("item"), 10, 64)
	if s.handleError(c, err) {
		return
	}

	second, err := strconv.ParseFloat(c.Query("second"), 64)
	if s.handleError(c, err) {
		return
	}

	logger.Infof("Setting main cover for item %d at %d", itemId, second)
	s.processor.EnqueueMainCover(itemId, second)
}

func (s *Server) getQueueMetadata(c *gin.Context) {
	queueMetadata, err := s.buildQueueMetadata()
	if s.handleError(c, err) {
		return
	}

	c.JSON(http.StatusOK, queueMetadata)
}

func (s *Server) clearFinishedTasks(c *gin.Context) {
	s.handleError(c, s.processor.ClearFinishedTasks())
}

func (s *Server) getTasks(c *gin.Context) {
	page, err := strconv.ParseInt(c.Query("page"), 10, 32)
	if s.handleError(c, err) {
		return
	}

	pageSize, err := strconv.ParseInt(c.Query("pageSize"), 10, 32)
	if s.handleError(c, err) {
		return
	}

	t, err := s.db.GetTasks(int((page-1)*pageSize), int(pageSize))
	if s.handleError(c, err) {
		return
	}

	tasks.AddDescriptionToTasks(s.db, t)
	c.JSON(http.StatusOK, t)
}

func (s *Server) queueContinue(c *gin.Context) {
	s.processor.Continue()
	c.Status(http.StatusOK)
}

func (s *Server) queuePause(c *gin.Context) {
	s.processor.Pause()
	c.Status(http.StatusOK)
}
