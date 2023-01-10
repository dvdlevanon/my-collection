package server

import (
	"encoding/json"
	"io"
	"my-collection/server/pkg/gallery"
	"my-collection/server/pkg/model"
	"my-collection/server/pkg/storage"
	"net/http"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-errors/errors"
	"github.com/op/go-logging"
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
	s.router.POST("/tags/:tag", s.updateTag)
	s.router.GET("/items/:item", s.getItem)
	s.router.GET("/tags/:tag", s.getTag)
	s.router.GET("/tags", s.getTags)
	s.router.GET("/items", s.getItems)
	s.router.GET("/items/refresh-preview", s.refreshItemsPreview)
	s.router.GET("/stream/*path", s.streamFile)
	s.router.GET("/storage/*path", s.getStorageFile)
}

func (s *Server) createItem(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, errors.Wrap(err, 0))
		return
	}

	var item model.Item
	if err = json.Unmarshal(body, &item); err != nil {
		c.AbortWithError(http.StatusInternalServerError, errors.Wrap(err, 0))
		return
	}

	if err = s.gallery.CreateOrUpdateItem(&item); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, model.Item{Id: item.Id})
}

func (s *Server) createTag(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, errors.Wrap(err, 0))
		return
	}

	var tag model.Tag
	if err = json.Unmarshal(body, &tag); err != nil {
		c.AbortWithError(http.StatusInternalServerError, errors.Wrap(err, 0))
		return
	}

	if err = s.gallery.CreateOrUpdateTag(&tag); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, model.Item{Id: tag.Id})
}

func (s *Server) updateItem(c *gin.Context) {
	itemId, err := strconv.ParseUint(c.Param("item"), 10, 64)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, errors.Wrap(err, 0))
		return
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, errors.Wrap(err, 0))
		return
	}

	var item model.Item
	if err = json.Unmarshal(body, &item); err != nil {
		c.AbortWithError(http.StatusInternalServerError, errors.Wrap(err, 0))
		return
	}

	if item.Id != 0 && item.Id != itemId {
		c.AbortWithError(http.StatusInternalServerError, errors.Errorf("Mismatch IDs %d != %d", item.Id, itemId))
		return
	}

	item.Id = itemId

	if err = s.gallery.UpdateItem(&item); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}

func (s *Server) updateTag(c *gin.Context) {
	tagId, err := strconv.ParseUint(c.Param("tag"), 10, 64)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, errors.Wrap(err, 0))
		return
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, errors.Wrap(err, 0))
		return
	}

	var tag model.Tag
	if err = json.Unmarshal(body, &tag); err != nil {
		c.AbortWithError(http.StatusInternalServerError, errors.Wrap(err, 0))
		return
	}

	if tag.Id != 0 && tag.Id != tagId {
		c.AbortWithError(http.StatusInternalServerError, errors.Errorf("Mismatch IDs %d != %d", tag.Id, tagId))
		return
	}

	tag.Id = tagId

	if err = s.gallery.UpdateTag(&tag); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}

func (s *Server) getItem(c *gin.Context) {
	itemId, err := strconv.ParseUint(c.Param("item"), 10, 64)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, errors.Wrap(err, 0))
		return
	}

	item, err := s.gallery.GetItem(itemId)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, item)
}

func (s *Server) getTag(c *gin.Context) {
	tagId, err := strconv.ParseUint(c.Param("tag"), 10, 64)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, errors.Wrap(err, 0))
		return
	}

	tag, err := s.gallery.GetTag(tagId)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, tag)
}

func (s *Server) getTags(c *gin.Context) {
	tags, err := s.gallery.GetAllTags()

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	logger.Infof("Get tags return %d tags", len(*tags))
	c.JSON(http.StatusOK, tags)
}

func (s *Server) getItems(c *gin.Context) {
	items, err := s.gallery.GetAllItems()

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	logger.Infof("Get items return %d items", len(*items))
	c.JSON(http.StatusOK, items)
}

func (s *Server) refreshItemsPreview(c *gin.Context) {
	err := s.gallery.RefreshItemsPreview()

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
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
