package server

import (
	"my-collection/server/pkg/directories"
	"my-collection/server/pkg/gallery"
	itemprocessor "my-collection/server/pkg/item-processor"
	"my-collection/server/pkg/storage"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-errors/errors"
	"github.com/op/go-logging"
	"gorm.io/gorm"
)

var logger = logging.MustGetLogger("server")

type Server struct {
	router      *gin.Engine
	gallery     *gallery.Gallery
	storage     *storage.Storage
	processor   itemprocessor.ItemProcessor
	directories directories.Directories
}

func New(gallery *gallery.Gallery, storage *storage.Storage,
	directories directories.Directories, processor itemprocessor.ItemProcessor) *Server {
	gin.SetMode("release")

	server := &Server{
		router:      gin.New(),
		gallery:     gallery,
		storage:     storage,
		directories: directories,
		processor:   processor,
	}

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
