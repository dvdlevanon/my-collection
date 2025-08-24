package server

import (
	"encoding/json"
	"io"
	"my-collection/server/pkg/bl/directories"
	"my-collection/server/pkg/model"
	"net/http"

	"github.com/gin-gonic/gin"
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

	if err = directories.CreateOrUpdateDirectory(s.db, &directory); err != nil {
		s.handleError(c, err)
		return
	}

	s.dcc.DirectoryChanged()
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

	if err = directories.UpdateDirectoryTags(s.db, &directory); err != nil {
		s.handleError(c, err)
		return
	}

	s.dcc.DirectoryChanged()
	c.Status(http.StatusOK)
}

func (s *Server) excludeDirectory(c *gin.Context) {
	directoryPath := c.Param("directory")[1:]

	err := directories.ExcludeDirectory(s.db, directoryPath)
	if s.handleError(c, err) {
		return
	}

	s.dcc.DirectoryChanged()
	c.Status(http.StatusOK)
}

func (s *Server) runDirectoriesScan(c *gin.Context) {
	logger.Infof("Triggering directory scan")
	s.dcc.DirectoryChanged()
}
