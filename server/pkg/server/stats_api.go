package server

import (
	"my-collection/server/pkg/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) getStats(c *gin.Context) {
	logger.Infof("Getting server stats")

	itemsCount, err := s.db.GetItemsCount()
	if s.handleError(c, err) {
		return
	}

	tagsCount, err := s.db.GetTagsCount()
	if s.handleError(c, err) {
		return
	}

	totalDurationSeconds, err := s.db.GetTotalDurationSeconds()
	if s.handleError(c, err) {
		return
	}

	c.JSON(http.StatusOK, model.Stats{
		ItemsCount:           itemsCount,
		TagsCount:            tagsCount,
		TotalDurationSeconds: totalDurationSeconds,
	})
}
