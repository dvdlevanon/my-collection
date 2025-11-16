package management

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"my-collection/server/pkg/backup"
	"my-collection/server/pkg/model"
	"my-collection/server/pkg/server"
	"my-collection/server/pkg/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("management-handler")

type managementDb interface {
	model.ItemReader
	model.TagReader
	model.DbMetadataReader
}

type managementProcessor interface {
	EnqueueAllItemsCovers(ctx context.Context, force bool) error
	EnqueueAllItemsFileMetadata(ctx context.Context) error
	EnqueueAllItemsPreview(ctx context.Context, force bool) error
	EnqueueAllItemsVideoMetadata(ctx context.Context, force bool) error
	GenerateMixOnDemand(ctx context.Context, ctg model.CurrentTimeGetter, desc string, tags []model.Tag) (*model.Tag, error)
	EnqueueItemOptimizer()
	EnqueueSpecTagger()
}

func NewHandler(db managementDb, processor managementProcessor) *managementHandler {
	return &managementHandler{
		db:        db,
		processor: processor,
	}
}

type managementHandler struct {
	db        managementDb
	processor managementProcessor
}

func (s *managementHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/items/refresh-covers", s.refreshItemsCovers)
	rg.POST("/items/refresh-preview", s.refreshItemsPreview)
	rg.POST("/items/refresh-video-metadata", s.refreshItemsVideoMetadata)
	rg.POST("/items/refresh-file-metadata", s.refreshItemsFileMetadata)
	rg.POST("/spectagger/run", s.runSpecTagger)
	rg.POST("/itemsoptimizer/run", s.runItemsOptimizer)
	rg.POST("/mix-on-demand", s.generateMixOnDemand)
	rg.GET("/export-metadata.json", s.exportMetadata)
	rg.GET("/stats", s.getStats)
}

func (s *managementHandler) runSpecTagger(c *gin.Context) {
	logger.Infof("Triggering spec tagger")
	s.processor.EnqueueSpecTagger()
}

func (s *managementHandler) getStats(c *gin.Context) {
	ctx := server.ContextWithSubject(c)
	logger.Infof("Getting server stats")

	itemsCount, err := s.db.GetItemsCount(ctx)
	if server.HandleError(c, err) {
		return
	}

	tagsCount, err := s.db.GetTagsCount(ctx)
	if server.HandleError(c, err) {
		return
	}

	totalDurationSeconds, err := s.db.GetTotalDurationSeconds(ctx)
	if server.HandleError(c, err) {
		return
	}

	c.JSON(http.StatusOK, model.Stats{
		ItemsCount:           itemsCount,
		TagsCount:            tagsCount,
		TotalDurationSeconds: totalDurationSeconds,
	})
}

func (s *managementHandler) refreshItemsCovers(c *gin.Context) {
	ctx := server.ContextWithSubject(c)
	force, err := strconv.ParseBool(c.Query("force"))
	if err != nil {
		force = false
	}

	if server.HandleError(c, s.processor.EnqueueAllItemsCovers(ctx, force)) {
		return
	}

	c.Status(http.StatusOK)
}

func (s *managementHandler) refreshItemsPreview(c *gin.Context) {
	ctx := server.ContextWithSubject(c)
	force, err := strconv.ParseBool(c.Query("force"))
	if err != nil {
		force = false
	}

	if server.HandleError(c, s.processor.EnqueueAllItemsPreview(ctx, force)) {
		return
	}

	c.Status(http.StatusOK)
}

func (s *managementHandler) refreshItemsVideoMetadata(c *gin.Context) {
	ctx := server.ContextWithSubject(c)
	forceParam := c.Query("force")
	force, err := strconv.ParseBool(forceParam)
	if err != nil {
		force = false
	}

	if server.HandleError(c, s.processor.EnqueueAllItemsVideoMetadata(ctx, force)) {
		return
	}

	c.Status(http.StatusOK)
}

func (s *managementHandler) refreshItemsFileMetadata(c *gin.Context) {
	ctx := server.ContextWithSubject(c)
	if server.HandleError(c, s.processor.EnqueueAllItemsFileMetadata(ctx)) {
		return
	}

	c.Status(http.StatusOK)
}

func (s *managementHandler) runItemsOptimizer(c *gin.Context) {
	logger.Infof("Triggering items optimizer")
	s.processor.EnqueueItemOptimizer()
}

func (s *managementHandler) exportMetadata(c *gin.Context) {
	ctx := server.ContextWithSubject(c)
	jsonBytes := bytes.Buffer{}
	if server.HandleError(c, backup.Export(ctx, s.db, s.db, &jsonBytes)) {
		return
	}

	c.Header("Content-Type", "application/json")
	c.Header("Content-Disposition", "gallery-metadata.json")
	c.String(http.StatusOK, jsonBytes.String())
}

func (s *managementHandler) generateMixOnDemand(c *gin.Context) {
	ctx := server.ContextWithSubject(c)
	body, err := io.ReadAll(c.Request.Body)
	if server.HandleError(c, err) {
		return
	}

	desc := c.Query("desc")
	if desc == "" {
		server.HandleError(c, fmt.Errorf("empty desc"))
		return
	}

	var tags []model.Tag
	if server.HandleError(c, json.Unmarshal(body, &tags)) {
		return
	}

	result, err := s.processor.GenerateMixOnDemand(ctx, utils.NowTimeGetter{}, desc, tags)
	if server.HandleError(c, err) {
		return
	}

	c.JSON(http.StatusOK, result)
}
