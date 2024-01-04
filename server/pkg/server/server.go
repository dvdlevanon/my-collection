package server

import (
	"my-collection/server/pkg/db"
	"my-collection/server/pkg/model"
	processor "my-collection/server/pkg/processor"
	"my-collection/server/pkg/storage"
	"my-collection/server/pkg/utils"
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
	dcc       model.DirectoryChangedCallback
	push      *push
}

func New(db *db.Database, storage *storage.Storage, dcc model.DirectoryChangedCallback, processor processor.Processor) *Server {
	gin.SetMode("release")

	server := &Server{
		router:    gin.New(),
		db:        db,
		storage:   storage,
		dcc:       dcc,
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
	api.DELETE("/items/:item", s.deleteItem)
	api.GET("/items/:item/location", s.getItemLocation)
	api.POST("/items/:item/remove-tag/:tag", s.removeTagFromItem)
	api.POST("/items/:item/main-cover", s.setMainCover)
	api.POST("/items/:item/split", s.splitItem)
	api.POST("/items/:item/make-highlight", s.makeHighlight)
	api.GET("/items/:item/suggestions", s.getSuggestionsForItem)

	api.GET("/tags", s.getTags)
	api.GET("/special-tags", s.getSpecialTags)
	api.GET("/categories", s.getCategories)
	api.POST("/tags", s.createTag)
	api.POST("/tags/:tag", s.updateTag)
	api.GET("/tags/:tag", s.getTag)
	api.DELETE("/tags/:tag", s.removeTag)
	api.POST("/tags/:tag/auto-image", s.autoImage)
	api.GET("/tags/:tag/tag-custom-commands", s.getAllTagCustomCommands)
	api.DELETE("/tags/:tag/tit/:tit", s.removeTagImageFromTag)
	api.POST("/tags/:tag/images/:image", s.updateTagImage)

	api.GET("/directories", s.getDirectories)
	api.POST("/directories", s.createOrUpdateDirectory)
	api.GET("/directories/*directory", s.getDirectory)
	api.POST("/directories/tags/*directory", s.SetDirectoryTags)
	api.DELETE("/directories/*directory", s.excludeDirectory)

	api.POST("/tags/:tag/annotations", s.addAnnotationToTag)
	api.DELETE("/tags/:tag/annotations/:annotation-id", s.removeAnnotationFromTag)
	api.GET("/tags/:tag/available-annotations", s.getTagAvailableAnnotations)

	api.POST("/items/refresh-covers", s.refreshItemsCovers)
	api.POST("/items/refresh-preview", s.refreshItemsPreview)
	api.POST("/items/refresh-video-metadata", s.refreshItemsVideoMetadata)
	api.POST("/items/refresh-file-metadata", s.refreshItemsFileMetadata)
	api.GET("/file/*path", s.getFile)
	api.POST("/upload-file", s.uploadFile)
	api.POST("/upload-file-from-url", s.uploadFileFromUrl)
	api.GET("/export-metadata.json", s.exportMetadata)

	api.GET("/queue/metadata", s.getQueueMetadata)
	api.GET("/queue/tasks", s.getTasks)
	api.POST("/queue/continue", s.queueContinue)
	api.POST("/queue/pause", s.queuePause)
	api.POST("/queue/clear-finished", s.clearFinishedTasks)

	api.GET("/tag-image-types", s.getTagImageTypes)

	api.GET("/ws", s.push.websocket)

	s.router.Static("/ui", "ui/")
	s.router.StaticFile("/", "ui/index.html")
	s.router.GET("/spa/*route", func(c *gin.Context) {
		http.ServeFile(c.Writer, c.Request, "ui/index.html")
	})
}

func (s *Server) Run(addr string) error {
	logger.Infof("Starting server at address %s", addr)
	go s.push.run()
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

	utils.LogError(err)
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
