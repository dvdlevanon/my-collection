package server

import (
	"encoding/json"
	"io"
	"my-collection/server/pkg/bl/directories"
	"my-collection/server/pkg/model"
	"my-collection/server/pkg/relativasor"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"k8s.io/utils/pointer"
)

func (s *Server) getDirectories(c *gin.Context) {
	directories, err := s.db.GetAllDirectories()

	if s.handleError(c, err) {
		return
	}

	logger.Infof("Get directories return %d directories", len(*directories))
	c.JSON(http.StatusOK, directories)
}

func (s *Server) createOrUpdateDirectory(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if s.handleError(c, err) {
		return
	}

	var directory model.Directory
	if err = json.Unmarshal(body, &directory); err != nil {
		s.handleError(c, err)
		return
	}

	directory.ProcessingStart = pointer.Int64(time.Now().UnixMilli())
	directory.Path = relativasor.GetRelativePath(directory.Path)
	if err = s.db.CreateOrUpdateDirectory(&directory); err != nil {
		s.handleError(c, err)
		return
	}

	s.fswatch.DirectoryChanged(&directory)
	c.Status(http.StatusOK)
}

func (s *Server) SetDirectoryTags(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if s.handleError(c, err) {
		return
	}

	var directory model.Directory
	if err = json.Unmarshal(body, &directory); err != nil {
		s.handleError(c, err)
		return
	}

	if err = directories.SetDirectoryTags(s.db, &directory); err != nil {
		s.handleError(c, err)
		return
	}

	s.fswatch.DirectoryChanged(&directory)
	c.Status(http.StatusOK)
}

func (s *Server) getDirectory(c *gin.Context) {
	directoryPath := c.Param("directory")[1:]
	directory, err := s.db.GetDirectory("path = ?", directoryPath)
	if s.handleError(c, err) {
		return
	}

	c.JSON(http.StatusOK, directory)
}

func (s *Server) excludeDirectory(c *gin.Context) {
	directoryPath := c.Param("directory")[1:]

	err := directories.ExcludeDirectory(s.db, directoryPath)
	if s.handleError(c, err) {
		return
	}

	s.fswatch.DirectoryExcluded(directoryPath)
	c.Status(http.StatusOK)
}
