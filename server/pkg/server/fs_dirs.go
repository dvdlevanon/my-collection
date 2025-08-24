package server

import (
	"my-collection/server/pkg/bl/directories"
	"my-collection/server/pkg/bl/fs"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (s *Server) getFsDir(c *gin.Context) {
	path := c.Query("path")
	depth, err := strconv.ParseInt(c.Query("depth"), 10, 64)
	if s.handleError(c, err) {
		return
	}

	node, err := fs.GetFsTree(path, int(depth))
	if s.handleError(c, err) {
		return
	}

	node, err = directories.EnrichFsNode(s.db, node)
	if s.handleError(c, err) {
		return
	}

	c.JSON(http.StatusOK, node)
}

func (s *Server) includeDir(c *gin.Context) {
	path := c.Query("path")
	subdirs, err := strconv.ParseBool(c.Query("subdirs"))
	if s.handleError(c, err) {
		return
	}
	hierarchy, err := strconv.ParseBool(c.Query("hierarchy"))
	if s.handleError(c, err) {
		return
	}

	err = fs.IncludeDir(s.db, path, subdirs, hierarchy)
	if s.handleError(c, err) {
		return
	}

	s.dcc.DirectoryChanged()
	c.Status(http.StatusOK)
}

func (s *Server) excludeDir(c *gin.Context) {
	path := c.Query("path")

	err := fs.ExcludeDir(s.db, path)
	if s.handleError(c, err) {
		return
	}

	s.dcc.DirectoryChanged()
	c.Status(http.StatusOK)
}
