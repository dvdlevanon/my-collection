package server

import (
	"bytes"
	"my-collection/server/pkg/backup"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) exportMetadata(c *gin.Context) {
	jsonBytes := bytes.Buffer{}
	if s.handleError(c, backup.Export(s.db, s.db, &jsonBytes)) {
		return
	}

	c.Header("Content-Type", "application/json")
	c.Header("Content-Disposition", "gallery-metadata.json")
	c.String(http.StatusOK, jsonBytes.String())
}
