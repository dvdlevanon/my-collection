package tasks

import (
	"context"
	"my-collection/server/pkg/bl/tasks"
	"my-collection/server/pkg/model"
	"my-collection/server/pkg/server"
	"my-collection/server/pkg/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var logger = utils.MustGetLogger("tasks-handler")

type queueDb interface {
	model.TaskReader
	model.ItemReader
}

type queueProcessor interface {
	model.ProcessorStatus
	Continue()
	Pause()
	ClearFinishedTasks(ctx context.Context) error
}

func NewHandler(db queueDb, processor queueProcessor) *tasksHandler {
	return &tasksHandler{
		db:        db,
		processor: processor,
	}
}

type tasksHandler struct {
	db        queueDb
	processor queueProcessor
}

func (s *tasksHandler) RegisterRoutes(rg *gin.RouterGroup) {

	rg.GET("/queue/metadata", s.getQueueMetadata)
	rg.GET("/queue/tasks", s.getTasks)
	rg.POST("/queue/continue", s.queueContinue)
	rg.POST("/queue/pause", s.queuePause)
	rg.POST("/queue/clear-finished", s.clearFinishedTasks)

}

func (s *tasksHandler) getQueueMetadata(c *gin.Context) {
	ctx := server.ContextWithSubject(c)
	queueMetadata, err := tasks.BuildQueueMetadata(ctx, s.db, s.processor)
	if server.HandleError(ctx, c, err) {
		return
	}

	c.JSON(http.StatusOK, queueMetadata)
}

func (s *tasksHandler) clearFinishedTasks(c *gin.Context) {
	ctx := server.ContextWithSubject(c)
	server.HandleError(ctx, c, s.processor.ClearFinishedTasks(ctx))
}

func (s *tasksHandler) getTasks(c *gin.Context) {
	ctx := server.ContextWithSubject(c)
	page, err := strconv.ParseInt(c.Query("page"), 10, 32)
	if server.HandleError(ctx, c, err) {
		return
	}

	pageSize, err := strconv.ParseInt(c.Query("pageSize"), 10, 32)
	if server.HandleError(ctx, c, err) {
		return
	}

	t, err := s.db.GetTasks(ctx, int((page-1)*pageSize), int(pageSize))
	if server.HandleError(ctx, c, err) {
		return
	}

	c.JSON(http.StatusOK, t)
}

func (s *tasksHandler) queueContinue(c *gin.Context) {
	s.processor.Continue()
	c.Status(http.StatusOK)
}

func (s *tasksHandler) queuePause(c *gin.Context) {
	s.processor.Pause()
	c.Status(http.StatusOK)
}
