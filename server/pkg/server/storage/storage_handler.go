package storage

import (
	"fmt"
	"my-collection/server/pkg/model"
	"my-collection/server/pkg/relativasor"
	"my-collection/server/pkg/server"
	"my-collection/server/pkg/utils"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("storage-handler")

type storageHandlerStorage interface {
	model.StorageDownloader
	model.StorageUploader
}

func NewHandler(storage storageHandlerStorage) *storageHandler {
	return &storageHandler{
		storage: storage,
	}
}

type storageHandler struct {
	storage storageHandlerStorage
}

func (s *storageHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/file/*path", s.getFile)
	rg.POST("/upload-file", s.uploadFile)
	rg.POST("/upload-file-from-url", s.uploadFileFromUrl)
}

func (s *storageHandler) uploadFile(c *gin.Context) {
	form, err := c.MultipartForm()
	if server.HandleError(c, err) {
		return
	}

	path := form.Value["path"][0]
	file := form.File["file"][0]
	fileName := fmt.Sprintf("%s-%s", file.Filename, uuid.NewString())
	relativeFile := filepath.Join(path, fileName)
	storageFile, err := s.storage.GetFileForWriting(relativeFile)
	if server.HandleError(c, err) {
		return
	}

	if server.HandleError(c, c.SaveUploadedFile(file, storageFile)) {
		return
	}

	c.JSON(http.StatusOK, model.FileUrl{Url: s.storage.GetStorageUrl(relativeFile)})
}

func (s *storageHandler) uploadFileFromUrl(c *gin.Context) {
	url := c.Query("url")
	path := c.Query("path")
	storageFile, err := s.storage.GetFileForWriting(path)
	if server.HandleError(c, err) {
		return
	}

	if server.HandleError(c, utils.DownloadFile(url, storageFile)) {
		return
	}

	c.JSON(http.StatusOK, model.FileUrl{Url: s.storage.GetStorageUrl(path)})
}

func (s *storageHandler) getFile(c *gin.Context) {
	path := c.Param("path")[1:]
	var file string
	if s.storage.IsStorageUrl(path) {
		file = s.storage.GetFile(path)
	} else {
		file = relativasor.GetAbsoluteFile(path)
	}

	logger.Infof("Getting file %v", file)
	http.ServeFile(c.Writer, c.Request, file)
}
