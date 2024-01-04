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
