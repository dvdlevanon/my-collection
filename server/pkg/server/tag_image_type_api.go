package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) getTagImageTypes(c *gin.Context) {
	tits, err := s.db.GetAllTagImageTypes()
	if s.handleError(c, err) {
		return
	}

	logger.Infof("Get tag image types return %d tag image types", len(*tits))
	c.JSON(http.StatusOK, tits)
}
