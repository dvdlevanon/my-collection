package server

import (
	"my-collection/server/pkg/bl/tasks"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

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
