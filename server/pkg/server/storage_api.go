package server

import (
	"fmt"
	"my-collection/server/pkg/model"
	"my-collection/server/pkg/relativasor"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *Server) uploadFile(c *gin.Context) {
	form, err := c.MultipartForm()
	if s.handleError(c, err) {
		return
	}

	path := form.Value["path"][0]
	file := form.File["file"][0]
	fileName := fmt.Sprintf("%s-%s", file.Filename, uuid.NewString())
	relativeFile := filepath.Join(path, fileName)
	storageFile, err := s.storage.GetFileForWriting(relativeFile)
	if s.handleError(c, err) {
		return
	}

	if s.handleError(c, c.SaveUploadedFile(file, storageFile)) {
		return
	}

	c.JSON(http.StatusOK, model.FileUrl{Url: s.storage.GetStorageUrl(relativeFile)})
}

func (s *Server) getFile(c *gin.Context) {
	path := c.Param("path")[1:]
	var file string
	if s.storage.IsStorageUrl(path) {
		file = s.storage.GetFile(path)
	} else {
		file = relativasor.GetAbsoluteFile(path)
	}

	logger.Infof("Getting file %v", file)
	http.ServeFile(c.Writer, c.Request, file)
}
