package server

import "github.com/gin-gonic/gin"

func (s *Server) runSpecTagger(c *gin.Context) {
	logger.Infof("Triggering spec tagger")
	s.spectagger.Trigger()
}
