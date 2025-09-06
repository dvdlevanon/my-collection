package tags

import (
	"encoding/json"
	"io"
	"my-collection/server/pkg/automix"
	"my-collection/server/pkg/bl/directories"
	"my-collection/server/pkg/bl/items"
	"my-collection/server/pkg/bl/tag_annotations"
	"my-collection/server/pkg/bl/tags"
	"my-collection/server/pkg/mixondemand"
	"my-collection/server/pkg/model"
	"my-collection/server/pkg/server"
	"my-collection/server/pkg/spectagger"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-errors/errors"
	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("tags-handler")

type tagsHandlerDb interface {
	model.TagReaderWriter
	model.TagCustomCommandsReader
	model.TagImageWriter
	model.TagImageTypeReaderWriter
	model.TagAnnotationReaderWriter
}

type tagsHandlerProcessor interface {
	ProcessThumbnail(image *model.TagImage) error
}

func NewHandler(db tagsHandlerDb, storage model.StorageUploader, processor tagsHandlerProcessor) *tagsHandler {
	return &tagsHandler{
		db:        db,
		storage:   storage,
		processor: processor,
	}
}

type tagsHandler struct {
	db        tagsHandlerDb
	storage   model.StorageUploader
	processor tagsHandlerProcessor
}

func (s *tagsHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/special-tags", s.getSpecialTags)
	rg.GET("/categories", s.getCategories)
	rg = rg.Group("tags")
	rg.GET("", s.getTags)
	rg.POST("", s.createTag)
	rg.POST("/:tag", s.updateTag)
	rg.GET("/:tag", s.getTag)
	rg.DELETE("/:tag", s.removeTag)
	rg.POST("/:tag/auto-image", s.autoImage)
	rg.GET("/:tag/tag-custom-commands", s.getAllTagCustomCommands)
	rg.DELETE("/:tag/tit/:tit", s.removeTagImageFromTag)
	rg.POST("/:tag/images/:image", s.updateTagImage)
	rg.POST("/:tag/random-mix/include", s.randomMixInclude)
	rg.POST("/:tag/random-mix/exclude", s.randomMixExclude)
	rg.POST("/:tag/annotations", s.addAnnotationToTag)
	rg.DELETE("/:tag/annotations/:annotation-id", s.removeAnnotationFromTag)
	rg.GET("/:tag/available-annotations", s.getTagAvailableAnnotations)
	rg.GET("/tag-image-types", s.getTagImageTypes)
}

func (s *tagsHandler) createTag(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if server.HandleError(c, err) {
		return
	}

	var tag model.Tag
	if server.HandleError(c, json.Unmarshal(body, &tag)) {
		return
	}

	if server.HandleError(c, s.db.CreateOrUpdateTag(&tag)) {
		return
	}

	c.JSON(http.StatusOK, model.Tag{Id: tag.Id})
}

func (s *tagsHandler) updateTag(c *gin.Context) {
	tagId, err := strconv.ParseUint(c.Param("tag"), 10, 64)
	if server.HandleError(c, err) {
		return
	}

	body, err := io.ReadAll(c.Request.Body)
	if server.HandleError(c, err) {
		return
	}

	var tag model.Tag
	if server.HandleError(c, json.Unmarshal(body, &tag)) {
		return
	}

	if tag.Id != 0 && tag.Id != tagId {
		server.HandleError(c, errors.Errorf("Mismatch IDs %d != %d", tag.Id, tagId))
		return
	}

	tag.Id = tagId
	if server.HandleError(c, s.db.UpdateTag(&tag)) {
		return
	}

	c.Status(http.StatusOK)
}

func (s *tagsHandler) getTag(c *gin.Context) {
	tagId, err := strconv.ParseUint(c.Param("tag"), 10, 64)
	if server.HandleError(c, err) {
		return
	}

	tag, err := s.db.GetTag(tagId)
	if server.HandleError(c, err) {
		return
	}

	c.JSON(http.StatusOK, tag)
}

func (s *tagsHandler) removeTag(c *gin.Context) {
	tagId, err := strconv.ParseUint(c.Param("tag"), 10, 64)
	if server.HandleError(c, err) {
		return
	}

	if server.HandleError(c, s.db.RemoveTag(tagId)) {
		return
	}

	c.Status(http.StatusOK)
}

func (s *tagsHandler) getTags(c *gin.Context) {
	tags, err := s.db.GetAllTags()
	if server.HandleError(c, err) {
		return
	}

	logger.Infof("Get tags return %d tags", len(*tags))
	c.JSON(http.StatusOK, tags)
}

func (s *tagsHandler) getCategories(c *gin.Context) {
	categories, err := tags.GetCategories(s.db)
	if server.HandleError(c, err) {
		return
	}

	logger.Infof("Get categories return %d tags", len(*categories))
	c.JSON(http.StatusOK, categories)
}

func (s *tagsHandler) getSpecialTags(c *gin.Context) {
	tags, err := s.db.GetTagsWithoutChildren(
		directories.GetDirectoriesTagId(),
		automix.GetDailymixTagId(),
		mixondemand.GetMixOnDemandTagId(),
		spectagger.GetSpecTagId(),
		items.GetHighlightsTagId())

	if server.HandleError(c, err) {
		return
	}

	logger.Infof("Get special tags return %d tags", len(*tags))
	c.JSON(http.StatusOK, tags)
}

func (s *tagsHandler) autoImage(c *gin.Context) {
	tagId, err := strconv.ParseUint(c.Param("tag"), 10, 64)
	if server.HandleError(c, err) {
		return
	}

	tag, err := s.db.GetTag(tagId)
	if server.HandleError(c, err) {
		return
	}

	body, err := io.ReadAll(c.Request.Body)
	if server.HandleError(c, err) {
		return
	}

	var fileUrl model.FileUrl
	if err = json.Unmarshal(body, &fileUrl); err != nil {
		server.HandleError(c, err)
		return
	}

	if server.HandleError(c, tags.AutoImageChildren(s.storage, s.db, s.db, tag, fileUrl.Url)) {
		return
	}

	c.JSON(http.StatusOK, tag)
}

func (s *tagsHandler) getAllTagCustomCommands(c *gin.Context) {
	commands, err := s.db.GetAllTagCustomCommands()
	if server.HandleError(c, err) {
		return
	}

	logger.Infof("Get tags custom commands return %d commands", len(*commands))
	c.JSON(http.StatusOK, commands)
}

func (s *tagsHandler) removeTagImageFromTag(c *gin.Context) {
	tagId, err := strconv.ParseUint(c.Param("tag"), 10, 64)
	if server.HandleError(c, err) {
		return
	}

	titId, err := strconv.ParseUint(c.Param("tit"), 10, 64)
	if server.HandleError(c, err) {
		return
	}

	if server.HandleError(c, tags.RemoveTagImages(s.db, tagId, titId)) {
		return
	}

	c.Status(http.StatusOK)
}

func (s *tagsHandler) updateTagImage(c *gin.Context) {
	tagId, err := strconv.ParseUint(c.Param("tag"), 10, 64)
	if server.HandleError(c, err) {
		return
	}

	imageId, err := strconv.ParseUint(c.Param("image"), 10, 64)
	if server.HandleError(c, err) {
		return
	}

	body, err := io.ReadAll(c.Request.Body)
	if server.HandleError(c, err) {
		return
	}

	var image model.TagImage
	if server.HandleError(c, json.Unmarshal(body, &image)) {
		return
	}

	if image.Id != imageId {
		server.HandleError(c, errors.Errorf("Mismatch IDs %d != %d", image.Id, imageId))
		return
	}

	if image.TagId != tagId {
		server.HandleError(c, errors.Errorf("Mismatch IDs %d != %d", image.TagId, tagId))
		return
	}

	if server.HandleError(c, s.db.UpdateTagImage(&image)) {
		return
	}

	go s.processor.ProcessThumbnail(&image)
	c.Status(http.StatusOK)
}

func (s *tagsHandler) randomMixExclude(c *gin.Context) {
	tagId, err := strconv.ParseUint(c.Param("tag"), 10, 64)
	if server.HandleError(c, err) {
		return
	}

	tag, err := s.db.GetTag(tagId)
	if server.HandleError(c, err) {
		return
	}

	noRandom := true
	tag.NoRandom = &noRandom
	if server.HandleError(c, s.db.UpdateTag(tag)) {
		return
	}

	c.Status(http.StatusOK)
}

func (s *tagsHandler) randomMixInclude(c *gin.Context) {
	tagId, err := strconv.ParseUint(c.Param("tag"), 10, 64)
	if server.HandleError(c, err) {
		return
	}

	tag, err := s.db.GetTag(tagId)
	if server.HandleError(c, err) {
		return
	}

	noRandom := false
	tag.NoRandom = &noRandom
	if server.HandleError(c, s.db.UpdateTag(tag)) {
		return
	}

	c.Status(http.StatusOK)
}

func (s *tagsHandler) addAnnotationToTag(c *gin.Context) {
	tagId, err := strconv.ParseUint(c.Param("tag"), 10, 64)
	if server.HandleError(c, err) {
		return
	}

	body, err := io.ReadAll(c.Request.Body)
	if server.HandleError(c, err) {
		return
	}

	var annotation model.TagAnnotation
	if err = json.Unmarshal(body, &annotation); err != nil {
		server.HandleError(c, err)
		return
	}

	annotationId, err := tag_annotations.AddAnnotationToTag(s.db, s.db, tagId, annotation)
	if server.HandleError(c, err) {
		return
	}

	c.JSON(http.StatusOK, model.TagAnnotation{Id: annotationId})
}

func (s *tagsHandler) removeAnnotationFromTag(c *gin.Context) {
	tagId, err := strconv.ParseUint(c.Param("tag"), 10, 64)
	if server.HandleError(c, err) {
		return
	}

	annotationId, err := strconv.ParseUint(c.Param("annotation-id"), 10, 64)
	if server.HandleError(c, err) {
		return
	}

	if server.HandleError(c, s.db.RemoveTagAnnotationFromTag(tagId, annotationId)) {
		return
	}

	c.Status(http.StatusOK)
}

func (s *tagsHandler) getTagAvailableAnnotations(c *gin.Context) {
	tagId, err := strconv.ParseUint(c.Param("tag"), 10, 64)
	if server.HandleError(c, err) {
		return
	}

	availableAnnotations, err := tag_annotations.GetTagAvailableAnnotations(s.db, s.db, tagId)
	if server.HandleError(c, err) {
		return
	}

	c.JSON(http.StatusOK, availableAnnotations)
}

func (s *tagsHandler) getTagImageTypes(c *gin.Context) {
	tits, err := s.db.GetAllTagImageTypes()
	if server.HandleError(c, err) {
		return
	}

	logger.Infof("Get tag image types return %d tag image types", len(*tits))
	c.JSON(http.StatusOK, tits)
}
