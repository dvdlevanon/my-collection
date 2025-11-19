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

type subtitleHandlerOp interface {
	subtitles.SubtitlesLister
	subtitles.SubtitlesDownloader
	model.TempFileProvider
}

func NewHandler(db subtitleHandlerDb, op subtitleHandlerOp) *subtitleHandler {
	return &subtitleHandler{
		db: db,
		op: op,
	}
}

type subtitleHandler struct {
	db subtitleHandlerDb
	op subtitleHandlerOp
}

func (s *subtitleHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg = rg.Group("subtitles")
	rg.GET("", s.getSubtitle)
	rg.GET("/:item/available", s.getAvailalbeNames)
	rg.GET("/:item/online", s.getOnlineNames)
	rg.POST("/:item/download", s.downloadSubtitle)
	rg.DELETE("/delete", s.deleteSubtitle)
}

func (s *subtitleHandler) getSubtitle(c *gin.Context) {
	ctx := server.ContextWithSubject(c)
	url := c.Query("url")
	subtitle, err := subtitles.GetSubtitle(ctx, url)
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
	ctx := server.ContextWithSubject(c)
	itemId, err := strconv.ParseUint(c.Param("item"), 10, 64)
	if server.HandleError(c, err) {
		return
	}

	availableNames, err := subtitles.GetAvailableNames(ctx, s.db, itemId)
	if server.HandleError(c, err) {
		return
	}

	c.JSON(http.StatusOK, availableNames)
}

func (s *subtitleHandler) getOnlineNames(c *gin.Context) {
	ctx := server.ContextWithSubject(c)
	itemId, err := strconv.ParseUint(c.Param("item"), 10, 64)
	if server.HandleError(c, err) {
		return
	}

	aiTranslated, err := strconv.ParseBool(c.Query("aiTranslated"))
	if err != nil {
		aiTranslated = false
	}

	lang := c.Query("lang")

	availableNames, err := subtitles.GetOnlineNames(ctx, s.db, s.op, itemId, lang, aiTranslated)
	if server.HandleError(c, err) {
		return
	}

	c.JSON(http.StatusOK, availableNames)
}

func (s *subtitleHandler) downloadSubtitle(c *gin.Context) {
	ctx := server.ContextWithSubject(c)
	itemId, err := strconv.ParseUint(c.Param("item"), 10, 64)
	if server.HandleError(c, err) {
		return
	}

	subtitleId := c.Query("id")
	subtitleTitle := c.Query("title")
	subtitle := model.SubtitleMetadata{
		Id:    subtitleId,
		Title: subtitleTitle,
	}

	url, err := subtitles.Download(ctx, s.db, s.op, s.op, itemId, subtitle)
	if server.HandleError(c, err) {
		return
	}

	c.JSON(http.StatusOK, url)
}

func (s *subtitleHandler) deleteSubtitle(c *gin.Context) {
	ctx := server.ContextWithSubject(c)

	url := c.Query("url")
	err := subtitles.Delete(ctx, url)
	if server.HandleError(c, err) {
		return
	}

	c.Status(http.StatusOK)
}
