package subtitles

import (
	"my-collection/server/pkg/bl/subtitles"
	"my-collection/server/pkg/model"
	"my-collection/server/pkg/server"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("subtitles-handler")

type subtitleHandlerDb interface {
	model.ItemReader
}

func NewHandler(db subtitleHandlerDb) *subtitleHandler {
	return &subtitleHandler{
		db: db,
	}
}

type subtitleHandler struct {
	db subtitleHandlerDb
}

func (s *subtitleHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg = rg.Group("subtitles")
	rg.GET("/:item", s.getSubtitle)
	rg.GET("/:item/available", s.getAvailalbeNames)
}

func (s *subtitleHandler) getSubtitle(c *gin.Context) {
	itemId, err := strconv.ParseUint(c.Param("item"), 10, 64)
	if server.HandleError(c, err) {
		return
	}

	subtitleName := c.Query("name")
	subtitle, err := subtitles.GetSubtitle(s.db, itemId, subtitleName)
	if err == subtitles.ErrSubtitileNotFound {
		c.Status(http.StatusNoContent)
		return
	}
	if server.HandleError(c, err) {
		return
	}

	c.JSON(http.StatusOK, subtitle)
}

func (s *subtitleHandler) getAvailalbeNames(c *gin.Context) {
	itemId, err := strconv.ParseUint(c.Param("item"), 10, 64)
	if server.HandleError(c, err) {
		return
	}

	availableNames, err := subtitles.GetAvailableNames(s.db, itemId)
	if server.HandleError(c, err) {
		return
	}

	c.JSON(http.StatusOK, availableNames)
}
