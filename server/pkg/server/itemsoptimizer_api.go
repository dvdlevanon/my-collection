package server

import "github.com/gin-gonic/gin"

func (s *Server) runItemsOptimizer(c *gin.Context) {
	logger.Infof("Triggering items optimizer")
	s.itemsOptimizer.Trigger()
}
