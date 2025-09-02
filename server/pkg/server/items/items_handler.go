package items

import (
	"encoding/json"
	"io"
	"my-collection/server/pkg/bl/items"
	"my-collection/server/pkg/model"
	"my-collection/server/pkg/relativasor"
	"my-collection/server/pkg/server"
	"my-collection/server/pkg/suggestions"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-errors/errors"
	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("items-handler")

type itemsHandlerDb interface {
	model.ItemReaderWriter
	model.TagReader
}

type itemsHandlerProcessor interface {
	EnqueueItemVideoMetadata(id uint64)
	EnqueueItemCovers(id uint64)
	EnqueueCropFrame(id uint64, second float64, rect model.RectFloat)
	EnqueueItemPreview(id uint64)
	EnqueueItemFileMetadata(id uint64)
	EnqueueMainCover(id uint64, second float64)
}

type itemsHandlerOptimizer interface {
	HandleItem(item *model.Item)
}

func NewHandler(db itemsHandlerDb, processor itemsHandlerProcessor, optimizer itemsHandlerOptimizer) *itemsHandler {
	return &itemsHandler{
		db:        db,
		processor: processor,
		optimizer: optimizer,
	}
}

type itemsHandler struct {
	db        itemsHandlerDb
	processor itemsHandlerProcessor
	optimizer itemsHandlerOptimizer
}

func (s *itemsHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg = rg.Group("items")
	rg.GET("", s.getItems)
	rg.POST("", s.createItem)
	rg.POST("/:item", s.updateItem)
	rg.GET("/:item", s.getItem)
	rg.DELETE("/:item", s.deleteItem)
	rg.GET("/:item/location", s.getItemLocation)
	rg.POST("/:item/remove-tag/:tag", s.removeTagFromItem)
	rg.POST("/:item/main-cover", s.setMainCover)
	rg.POST("/:item/split", s.splitItem)
	rg.POST("/:item/make-highlight", s.makeHighlight)
	rg.POST("/:item/crop-frame", s.cropFrame)
	rg.GET("/:item/suggestions", s.getSuggestionsForItem)
	rg.POST("/:item/process", s.refreshItem)
	rg.POST("/:item/optimize", s.optimizeItem)
}

func (s *itemsHandler) createItem(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if server.HandleError(c, err) {
		return
	}

	var item model.Item
	if server.HandleError(c, json.Unmarshal(body, &item)) {
		return
	}

	item.Url = relativasor.GetRelativePath(item.Url)
	item.Origin = relativasor.GetRelativePath(item.Origin)
	if server.HandleError(c, s.db.CreateOrUpdateItem(&item)) {
		return
	}

	c.JSON(http.StatusOK, model.Item{Id: item.Id})
}

func (s *itemsHandler) updateItem(c *gin.Context) {
	itemId, err := strconv.ParseUint(c.Param("item"), 10, 64)
	if server.HandleError(c, err) {
		return
	}

	body, err := io.ReadAll(c.Request.Body)
	if server.HandleError(c, err) {
		return
	}

	var item model.Item
	if server.HandleError(c, json.Unmarshal(body, &item)) {
		return
	}

	if item.Id != 0 && item.Id != itemId {
		server.HandleError(c, errors.Errorf("Mismatch IDs %d != %d", item.Id, itemId))
		return
	}

	item.Id = itemId
	if server.HandleError(c, s.db.UpdateItem(&item)) {
		return
	}

	c.Status(http.StatusOK)
}

func (s *itemsHandler) getItem(c *gin.Context) {
	itemId, err := strconv.ParseUint(c.Param("item"), 10, 64)
	if server.HandleError(c, err) {
		return
	}

	item, err := s.db.GetItem(itemId)
	if server.HandleError(c, err) {
		return
	}

	c.JSON(http.StatusOK, item)
}

func (s *itemsHandler) deleteItem(c *gin.Context) {
	itemId, err := strconv.ParseUint(c.Param("item"), 10, 64)
	if server.HandleError(c, err) {
		return
	}

	deleteRealItem, err := strconv.ParseBool(c.Query("deleteRealFile"))
	if server.HandleError(c, err) {
		return
	}

	if deleteRealItem {
		if server.HandleError(c, items.DeleteRealFile(s.db, itemId)) {
			return
		}
	}

	errs := items.RemoveItemAndItsAssociations(s.db, itemId)
	if len(errs) > 0 {
		if server.HandleError(c, errs[0]) {
			return
		}
	}

	c.Status(http.StatusOK)
}

func (s *itemsHandler) getItemLocation(c *gin.Context) {
	itemId, err := strconv.ParseUint(c.Param("item"), 10, 64)
	if server.HandleError(c, err) {
		return
	}

	item, err := s.db.GetItem(itemId)
	if server.HandleError(c, err) {
		return
	}

	c.JSON(http.StatusOK, model.FileUrl{
		Url: relativasor.GetAbsoluteFile(item.Url),
	})
}

func (s *itemsHandler) getItems(c *gin.Context) {
	items, err := s.db.GetAllItems()
	if server.HandleError(c, err) {
		return
	}

	logger.Infof("Get items return %d items", len(*items))
	c.JSON(http.StatusOK, items)
}

func (s *itemsHandler) removeTagFromItem(c *gin.Context) {
	itemId, err := strconv.ParseUint(c.Param("item"), 10, 64)
	if server.HandleError(c, err) {
		return
	}

	tagId, err := strconv.ParseUint(c.Param("tag"), 10, 64)
	if server.HandleError(c, err) {
		return
	}

	if server.HandleError(c, s.db.RemoveTagFromItem(itemId, tagId)) {
		return
	}

	c.Status(http.StatusOK)
}

func (s *itemsHandler) setMainCover(c *gin.Context) {
	itemId, err := strconv.ParseUint(c.Param("item"), 10, 64)
	if server.HandleError(c, err) {
		return
	}

	second, err := strconv.ParseFloat(c.Query("second"), 64)
	if server.HandleError(c, err) {
		return
	}

	logger.Infof("Setting main cover for item %d at %d", itemId, second)
	s.processor.EnqueueMainCover(itemId, second)
	c.Status(http.StatusOK)
}

func (s *itemsHandler) splitItem(c *gin.Context) {
	itemId, err := strconv.ParseUint(c.Param("item"), 10, 64)
	if server.HandleError(c, err) {
		return
	}

	second, err := strconv.ParseFloat(c.Query("second"), 64)
	if server.HandleError(c, err) {
		return
	}

	logger.Infof("Splitting item %d at %f", itemId, second)
	changedItems, err := items.Split(s.db, itemId, second)
	if server.HandleError(c, err) {
		return
	}

	for _, item := range changedItems {
		s.processor.EnqueueItemVideoMetadata(item.Id)
		s.processor.EnqueueItemCovers(item.Id)
		s.processor.EnqueueItemFileMetadata(item.Id)
		// s.processor.EnqueueItemPreview(item.Id)
	}

	c.Status(http.StatusOK)
}

func (s *itemsHandler) makeHighlight(c *gin.Context) {
	itemId, err := strconv.ParseUint(c.Param("item"), 10, 64)
	if server.HandleError(c, err) {
		return
	}

	startSecond, err := strconv.ParseFloat(c.Query("start"), 64)
	if server.HandleError(c, err) {
		return
	}

	endSecond, err := strconv.ParseFloat(c.Query("end"), 64)
	if server.HandleError(c, err) {
		return
	}

	highlightId, err := strconv.ParseUint(c.Query("highlight-id"), 10, 64)
	if server.HandleError(c, err) {
		return
	}

	logger.Infof("Making highlight for item %d from %f to %f", itemId, startSecond, endSecond)
	highlightItem, err := items.MakeHighlight(s.db, itemId, startSecond, endSecond, highlightId)
	if server.HandleError(c, err) {
		return
	}

	s.processor.EnqueueItemVideoMetadata(highlightItem.Id)
	s.processor.EnqueueItemCovers(highlightItem.Id)
	s.processor.EnqueueItemPreview(highlightItem.Id)
	s.processor.EnqueueItemFileMetadata(highlightItem.Id)
	c.Status(http.StatusOK)
}

func (s *itemsHandler) cropFrame(c *gin.Context) {
	itemId, err := strconv.ParseUint(c.Param("item"), 10, 64)
	if server.HandleError(c, err) {
		return
	}

	second, err := strconv.ParseFloat(c.Query("second"), 64)
	if server.HandleError(c, err) {
		return
	}

	cropX, err := strconv.ParseFloat(c.Query("crop-x"), 64)
	if server.HandleError(c, err) {
		return
	}

	cropY, err := strconv.ParseFloat(c.Query("crop-y"), 64)
	if server.HandleError(c, err) {
		return
	}

	cropWidth, err := strconv.ParseFloat(c.Query("crop-width"), 64)
	if server.HandleError(c, err) {
		return
	}

	cropHeight, err := strconv.ParseFloat(c.Query("crop-height"), 64)
	if server.HandleError(c, err) {
		return
	}

	rect := model.RectFloat{
		X: cropX,
		Y: cropY,
		W: cropWidth,
		H: cropHeight,
	}

	logger.Infof("Cropping frame for item %d at %f %s", itemId, second, rect)
	s.processor.EnqueueCropFrame(itemId, second, rect)
}

func (s *itemsHandler) getSuggestionsForItem(c *gin.Context) {
	itemId, err := strconv.ParseUint(c.Param("item"), 10, 64)
	if server.HandleError(c, err) {
		return
	}

	result, err := suggestions.GetSuggestionsForItem(s.db, s.db, itemId, 8)
	if server.HandleError(c, err) {
		return
	}

	c.JSON(http.StatusOK, result)
}

func (s *itemsHandler) refreshItem(c *gin.Context) {
	itemId, err := strconv.ParseUint(c.Param("item"), 10, 64)
	if server.HandleError(c, err) {
		return
	}

	s.processor.EnqueueItemVideoMetadata(itemId)
	s.processor.EnqueueItemCovers(itemId)
	s.processor.EnqueueItemPreview(itemId)
	s.processor.EnqueueItemFileMetadata(itemId)

	c.Status(http.StatusOK)
}

func (s *itemsHandler) optimizeItem(c *gin.Context) {
	itemId, err := strconv.ParseUint(c.Param("item"), 10, 64)
	if server.HandleError(c, err) {
		return
	}

	item, err := s.db.GetItem(itemId)
	if server.HandleError(c, err) {
		return
	}

	s.optimizer.HandleItem(item)
	c.Status(http.StatusOK)
}
