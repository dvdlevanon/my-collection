package server

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

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

func (s *Server) refreshItemsFileMetadata(c *gin.Context) {
	if s.handleError(c, s.processor.EnqueueAllItemsFileMetadata()) {
		return
	}

	c.Status(http.StatusOK)
}

func (s *Server) refreshItem(c *gin.Context) {
	itemId, err := strconv.ParseUint(c.Param("item"), 10, 64)
	if s.handleError(c, err) {
		return
	}

	s.processor.EnqueueItemVideoMetadata(itemId)
	s.processor.EnqueueItemCovers(itemId)
	s.processor.EnqueueItemPreview(itemId)
	s.processor.EnqueueItemFileMetadata(itemId)

	c.Status(http.StatusOK)
}

func (s *Server) optimizeItem(c *gin.Context) {
	itemId, err := strconv.ParseUint(c.Param("item"), 10, 64)
	if s.handleError(c, err) {
		return
	}

	item, err := s.db.GetItem(itemId)
	if s.handleError(c, err) {
		return
	}

	s.itemsOptimizer.HandleItem(item)
	c.Status(http.StatusOK)
}
