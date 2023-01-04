package server

import (
	"encoding/json"
	"io"
	"my-collection/server/pkg/db"
	"my-collection/server/pkg/model"
	"net/http"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-errors/errors"
)

type Server struct {
	router *gin.Engine
	db     *db.Database
}

func New(db *db.Database) *Server {
	server := &Server{
		router: gin.Default(),
		db:     db,
	}

	server.init()
	return server
}

func (s *Server) init() {
	s.router.Use(cors.Default())
	s.router.POST("/items", s.createItem)
	s.router.POST("/tags", s.createTag)
	s.router.POST("/items/:item", s.updateItem)
	s.router.POST("/tags/:tag", s.updateTag)
	s.router.GET("/items/:item", s.getItem)
	s.router.GET("/tags/:tag", s.getTag)
	s.router.GET("/tags", s.getTags)
	s.router.GET("/items", s.getItems)
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

	if err = s.db.CreateOrUpdateItem(&item); err != nil {
		c.AbortWithError(http.StatusInternalServerError, errors.Wrap(err, 0))
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

	if err = s.db.CreateOrUpdateTag(&tag); err != nil {
		c.AbortWithError(http.StatusInternalServerError, errors.Wrap(err, 0))
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

	if err = s.db.UpdateItem(&item); err != nil {
		c.AbortWithError(http.StatusInternalServerError, errors.Wrap(err, 0))
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

	if err = s.db.UpdateTag(&tag); err != nil {
		c.AbortWithError(http.StatusInternalServerError, errors.Wrap(err, 0))
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

	item, err := s.db.GetItem(itemId)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, errors.Wrap(err, 0))
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

	tag, err := s.db.GetTag(tagId)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, errors.Wrap(err, 0))
		return
	}

	c.JSON(http.StatusOK, tag)
}

func (s *Server) getTags(c *gin.Context) {
	tags, err := s.db.GetAllTags()

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, errors.Wrap(err, 0))
		return
	}

	c.JSON(http.StatusOK, tags)
}

func (s *Server) getItems(c *gin.Context) {
	items, err := s.db.GetAllItems()

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, errors.Wrap(err, 0))
		return
	}

	c.JSON(http.StatusOK, items)
}

func (s *Server) Run() {
	s.router.Run()
}
