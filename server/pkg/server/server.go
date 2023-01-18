package server

import (
	"encoding/json"
	"io"
	"my-collection/server/pkg/gallery"
	"my-collection/server/pkg/model"
	"my-collection/server/pkg/storage"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-errors/errors"
	"github.com/op/go-logging"
	"gorm.io/gorm"
)

var logger = logging.MustGetLogger("server")

type Server struct {
	router  *gin.Engine
	gallery *gallery.Gallery
	storage *storage.Storage
}

func New(gallery *gallery.Gallery, storage *storage.Storage) *Server {
	gin.SetMode("release")

	server := &Server{
		router:  gin.New(),
		gallery: gallery,
		storage: storage,
	}

	server.init()
	return server
}

func (s *Server) init() {
	s.router.Use(cors.Default())
	s.router.Use(httpLogger)
	s.router.POST("/items", s.createItem)
	s.router.POST("/tags", s.createTag)
	s.router.POST("/items/:item", s.updateItem)
	s.router.POST("/items/:item/remove-tag/:tag", s.removeTagFromItem)
	s.router.POST("/tags/:tag", s.updateTag)
	s.router.GET("/items/:item", s.getItem)
	s.router.GET("/tags/:tag", s.getTag)
	s.router.GET("/tags", s.getTags)
	s.router.GET("/items", s.getItems)
	s.router.GET("/items/refresh-covers", s.refreshItemsCovers)
	s.router.GET("/stream/*path", s.streamFile)
	s.router.GET("/storage/*path", s.getStorageFile)
	s.router.POST("/upload-file", s.uploadFile)
}

func (s *Server) uploadFile(c *gin.Context) {
	form, err := c.MultipartForm()
	if s.handleError(c, err) {
		return
	}

	path := form.Value["path"][0]
	file := form.File["file"][0]
	relativeFile := filepath.Join(path, file.Filename)
	storageFile, err := s.storage.GetFileForWriting(relativeFile)
	if s.handleError(c, err) {
		return
	}

	if s.handleError(c, c.SaveUploadedFile(file, storageFile)) {
		return
	}

	c.JSON(http.StatusOK, model.FileUrl{Url: relativeFile})
}

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

	c.JSON(http.StatusOK, model.Item{Id: tag.Id})
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
	err := s.gallery.RefreshItemsCovers()

	if s.handleError(c, err) {
		return
	}

	c.Status(http.StatusOK)
}

func (s *Server) getStorageFile(c *gin.Context) {
	path := c.Param("path")
	http.ServeFile(c.Writer, c.Request, s.storage.GetFile(path))
}

func (s *Server) streamFile(c *gin.Context) {
	path := c.Param("path")
	http.ServeFile(c.Writer, c.Request, s.gallery.GetItemAbsolutePath(path))
}

func (s *Server) Run(addr string) {
	logger.Infof("Starting server at address %s", addr)
	s.router.Run(addr)
}

func (s *Server) handleError(c *gin.Context, err error) bool {
	if err == nil {
		return false
	}

	httpError := http.StatusInternalServerError
	if errors.Is(err, gorm.ErrRecordNotFound) {
		httpError = http.StatusNotFound
	}

	c.AbortWithError(httpError, err)
	return true
}
