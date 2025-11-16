package fs

import (
	"encoding/json"
	"io"
	"my-collection/server/pkg/bl/directories"
	"my-collection/server/pkg/bl/fs"
	"my-collection/server/pkg/model"
	"my-collection/server/pkg/server"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("fs-handler")

type fsDb interface {
	model.DirectoryReaderWriter
}

type fsDirectoryChangedListener interface {
	DirectoryChanged()
}

func NewHandler(db fsDb, changeListener fsDirectoryChangedListener) *fsHandler {
	return &fsHandler{
		db:             db,
		changeListener: changeListener,
	}
}

type fsHandler struct {
	db             fsDb
	changeListener fsDirectoryChangedListener
}

func (s *fsHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/directories/scan", s.runDirectoriesScan)
	rg.POST("/directories/tags/*directory", s.SetDirectoryTags)
	rg.GET("/fs", s.getFsDir)
	rg.POST("/fs/include", s.includeDir)
	rg.POST("/fs/exclude", s.excludeDir)
}

func (s *fsHandler) getFsDir(c *gin.Context) {
	ctx := server.ContextWithSubject(c)
	path := c.Query("path")
	depth, err := strconv.ParseInt(c.Query("depth"), 10, 64)
	if server.HandleError(c, err) {
		return
	}

	node, err := fs.GetFsTree(path, int(depth))
	if server.HandleError(c, err) {
		return
	}

	node, err = directories.EnrichFsNode(ctx, s.db, node)
	if server.HandleError(c, err) {
		return
	}

	c.JSON(http.StatusOK, node)
}

func (s *fsHandler) includeDir(c *gin.Context) {
	ctx := server.ContextWithSubject(c)
	path := c.Query("path")
	subdirs, err := strconv.ParseBool(c.Query("subdirs"))
	if server.HandleError(c, err) {
		return
	}
	hierarchy, err := strconv.ParseBool(c.Query("hierarchy"))
	if server.HandleError(c, err) {
		return
	}

	err = fs.IncludeDir(ctx, s.db, path, subdirs, hierarchy)
	if server.HandleError(c, err) {
		return
	}

	s.changeListener.DirectoryChanged()
	c.Status(http.StatusOK)
}

func (s *fsHandler) excludeDir(c *gin.Context) {
	ctx := server.ContextWithSubject(c)
	path := c.Query("path")

	err := fs.ExcludeDir(ctx, s.db, path)
	if server.HandleError(c, err) {
		return
	}

	s.changeListener.DirectoryChanged()
	c.Status(http.StatusOK)
}

func (s *fsHandler) SetDirectoryTags(c *gin.Context) {
	ctx := server.ContextWithSubject(c)
	body, err := io.ReadAll(c.Request.Body)
	if server.HandleError(c, err) {
		return
	}

	var directory model.Directory
	if err = json.Unmarshal(body, &directory); err != nil {
		server.HandleError(c, err)
		return
	}

	if err = directories.UpdateDirectoryTags(ctx, s.db, &directory); err != nil {
		server.HandleError(c, err)
		return
	}

	s.changeListener.DirectoryChanged()
	c.Status(http.StatusOK)
}

func (s *fsHandler) runDirectoriesScan(c *gin.Context) {
	logger.Infof("Triggering directory scan")
	s.changeListener.DirectoryChanged()
}
