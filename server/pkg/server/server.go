package server

import (
	"my-collection/server/pkg/gallery"
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
	router  *gin.Engine
	gallery *gallery.Gallery
	storage *storage.Storage
}

func New(gallery *gallery.Gallery, storage *storage.Storage) *Server {
	gin.SetMode("release")

	server := &Server{
		router:  gin.New(),
		gallery: gallery,
		storage: storage,
	}

	server.init()
	return server
}

func (s *Server) init() {
	s.router.Use(cors.Default())
	s.router.Use(httpLogger)

	api := s.router.Group("/api")
	api.POST("/items", s.createItem)
	api.POST("/tags", s.createTag)
	api.POST("/items/:item", s.updateItem)
	api.POST("/items/:item/remove-tag/:tag", s.removeTagFromItem)
	api.POST("/tags/:tag", s.updateTag)
	api.GET("/items/:item", s.getItem)
	api.GET("/tags/:tag", s.getTag)
	api.GET("/tags/:tag/available-annotations", s.getTagAvailableAnnotations)
	api.GET("/tags", s.getTags)
	api.GET("/items", s.getItems)
	api.GET("/items/refresh-covers", s.refreshItemsCovers)
	api.GET("/items/refresh-preview", s.refreshItemsPreview)
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
