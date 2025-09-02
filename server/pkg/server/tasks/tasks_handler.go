package tasks

import (
	"my-collection/server/pkg/bl/tasks"
	"my-collection/server/pkg/model"
	"my-collection/server/pkg/server"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("tasks-handler")

type queueDb interface {
	model.TaskReader
	model.ItemReader
}

type queueProcessor interface {
	model.ProcessorStatus
	Continue()
	Pause()
	ClearFinishedTasks() error
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
	queueMetadata, err := tasks.BuildQueueMetadata(s.db, s.processor)
	if server.HandleError(c, err) {
		return
	}

	c.JSON(http.StatusOK, queueMetadata)
}

func (s *tasksHandler) clearFinishedTasks(c *gin.Context) {
	server.HandleError(c, s.processor.ClearFinishedTasks())
}

func (s *tasksHandler) getTasks(c *gin.Context) {
	page, err := strconv.ParseInt(c.Query("page"), 10, 32)
	if server.HandleError(c, err) {
		return
	}

	pageSize, err := strconv.ParseInt(c.Query("pageSize"), 10, 32)
	if server.HandleError(c, err) {
		return
	}

	t, err := s.db.GetTasks(int((page-1)*pageSize), int(pageSize))
	if server.HandleError(c, err) {
		return
	}

	tasks.AddDescriptionToTasks(s.db, t)
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
