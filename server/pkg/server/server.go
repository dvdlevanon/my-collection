package server

import (
	"my-collection/server/pkg/db"
	"my-collection/server/pkg/fswatch"
	"my-collection/server/pkg/model"
	processor "my-collection/server/pkg/processor"
	"my-collection/server/pkg/storage"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-errors/errors"
	"github.com/op/go-logging"
	"gorm.io/gorm"
	"k8s.io/utils/pointer"
)

var logger = logging.MustGetLogger("server")

type Server struct {
	router    *gin.Engine
	db        *db.Database
	storage   *storage.Storage
	processor processor.Processor
	fswatch   fswatch.Fswatch
	push      *push
}

func New(db *db.Database, storage *storage.Storage, fswatch fswatch.Fswatch, processor processor.Processor) *Server {
	gin.SetMode("release")

	server := &Server{
		router:    gin.New(),
		db:        db,
		storage:   storage,
		fswatch:   fswatch,
		processor: processor,
	}

	server.push = newPush(processor, server)
	server.init()
	return server
}

func (s *Server) init() {
	s.router.Use(cors.Default())
	s.router.Use(httpLogger)

	api := s.router.Group("/api")

	api.GET("/items", s.getItems)
	api.POST("/items", s.createItem)
	api.POST("/items/:item", s.updateItem)
	api.GET("/items/:item", s.getItem)
	api.POST("/items/:item/remove-tag/:tag", s.removeTagFromItem)
	api.POST("/items/:item/main-cover", s.setMainCover)

	api.GET("/tags", s.getTags)
	api.POST("/tags", s.createTag)
	api.POST("/tags/:tag", s.updateTag)
	api.GET("/tags/:tag", s.getTag)
	api.DELETE("/tags/:tag", s.removeTag)
	api.POST("/tags/:tag/auto-image", s.autoImage)

	api.GET("/directories", s.getDirectories)
	api.POST("/directories", s.createOrUpdateDirectory)
	api.GET("/directories/*directory", s.getDirectory)
	api.POST("/directories/tags/*directory", s.SetDirectoryTags)
	api.DELETE("/directories/*directory", s.excludeDirectory)

	api.POST("/tags/:tag/annotations", s.addAnnotationToTag)
	api.DELETE("/tags/:tag/annotations/:annotation-id", s.removeAnnotationFromTag)
	api.GET("/tags/:tag/available-annotations", s.getTagAvailableAnnotations)

	api.GET("/items/refresh-covers", s.refreshItemsCovers)
	api.GET("/items/refresh-preview", s.refreshItemsPreview)
	api.GET("/items/refresh-video-metadata", s.refreshItemsVideoMetadata)
	api.GET("/file/*path", s.getFile)
	api.POST("/upload-file", s.uploadFile)
	api.GET("/export-metadata.json", s.exportMetadata)

	api.GET("/queue/metadata", s.getQueueMetadata)
	api.GET("/queue/tasks", s.getTasks)
	api.POST("/queue/continue", s.queueContinue)
	api.POST("/queue/pause", s.queuePause)
	api.POST("/queue/clear-finished", s.clearFinishedTasks)

	api.GET("/ws", s.push.websocket)

	s.router.Static("/ui", "ui/")
	s.router.StaticFile("/", "ui/index.html")
	s.router.GET("/spa/*route", func(c *gin.Context) {
		http.ServeFile(c.Writer, c.Request, "ui/index.html")
	})
}

func (s *Server) Run(addr string) error {
	logger.Infof("Starting server at address %s", addr)
	return s.router.Run(addr)
}

func (s *Server) handleError(c *gin.Context, err error) bool {
	if err == nil {
		return false
	}

	httpError := http.StatusInternalServerError
	if errors.Is(err, gorm.ErrRecordNotFound) {
		httpError = http.StatusNotFound
	}

	c.AbortWithError(httpError, err)
	return true
}

func (s *Server) buildQueueMetadata() (model.QueueMetadata, error) {
	size, err := s.db.TasksCount("")
	if err != nil {
		logger.Errorf("Unable to get queue size %s", err)
		return model.QueueMetadata{}, nil
	}

	unfinishedTasks, err := s.db.TasksCount("processing_end is null")
	if err != nil {
		logger.Errorf("Unable to get unfinished tasks count %s", err)
		return model.QueueMetadata{}, nil
	}

	queueMetadata := model.QueueMetadata{
		Size:            pointer.Int64(size),
		Paused:          pointer.Bool(s.processor.IsPaused()),
		UnfinishedTasks: pointer.Int64(unfinishedTasks),
	}

	return queueMetadata, err
}
