package server

import (
	"encoding/json"
	"io"
	"my-collection/server/pkg/bl/items"
	"my-collection/server/pkg/model"
	"my-collection/server/pkg/relativasor"
	"my-collection/server/pkg/suggestions"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-errors/errors"
)

func (s *Server) createItem(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if s.handleError(c, err) {
		return
	}

	var item model.Item
	if s.handleError(c, json.Unmarshal(body, &item)) {
		return
	}

	item.Url = relativasor.GetRelativePath(item.Url)
	item.Origin = relativasor.GetRelativePath(item.Origin)
	if s.handleError(c, s.db.CreateOrUpdateItem(&item)) {
		return
	}

	c.JSON(http.StatusOK, model.Item{Id: item.Id})
}

func (s *Server) updateItem(c *gin.Context) {
	itemId, err := strconv.ParseUint(c.Param("item"), 10, 64)
	if s.handleError(c, err) {
		return
	}

	body, err := io.ReadAll(c.Request.Body)
	if s.handleError(c, err) {
		return
	}

	var item model.Item
	if s.handleError(c, json.Unmarshal(body, &item)) {
		return
	}

	if item.Id != 0 && item.Id != itemId {
		s.handleError(c, errors.Errorf("Mismatch IDs %d != %d", item.Id, itemId))
		return
	}

	item.Id = itemId
	if s.handleError(c, s.db.UpdateItem(&item)) {
		return
	}

	c.Status(http.StatusOK)
}

func (s *Server) getItem(c *gin.Context) {
	itemId, err := strconv.ParseUint(c.Param("item"), 10, 64)
	if s.handleError(c, err) {
		return
	}

	item, err := s.db.GetItem(itemId)
	if s.handleError(c, err) {
		return
	}

	c.JSON(http.StatusOK, item)
}

func (s *Server) getItemLocation(c *gin.Context) {
	itemId, err := strconv.ParseUint(c.Param("item"), 10, 64)
	if s.handleError(c, err) {
		return
	}

	item, err := s.db.GetItem(itemId)
	if s.handleError(c, err) {
		return
	}

	c.JSON(http.StatusOK, model.FileUrl{
		Url: relativasor.GetAbsoluteFile(item.Url),
	})
}

func (s *Server) getItems(c *gin.Context) {
	items, err := s.db.GetAllItems()
	if s.handleError(c, err) {
		return
	}

	logger.Infof("Get items return %d items", len(*items))
	c.JSON(http.StatusOK, items)
}

func (s *Server) removeTagFromItem(c *gin.Context) {
	itemId, err := strconv.ParseUint(c.Param("item"), 10, 64)
	if s.handleError(c, err) {
		return
	}

	tagId, err := strconv.ParseUint(c.Param("tag"), 10, 64)
	if s.handleError(c, err) {
		return
	}

	if s.handleError(c, s.db.RemoveTagFromItem(itemId, tagId)) {
		return
	}

	c.Status(http.StatusOK)
}

func (s *Server) setMainCover(c *gin.Context) {
	itemId, err := strconv.ParseUint(c.Param("item"), 10, 64)
	if s.handleError(c, err) {
		return
	}

	second, err := strconv.ParseFloat(c.Query("second"), 64)
	if s.handleError(c, err) {
		return
	}

	logger.Infof("Setting main cover for item %d at %d", itemId, second)
	s.processor.EnqueueMainCover(itemId, second)
	c.Status(http.StatusOK)
}

func (s *Server) splitItem(c *gin.Context) {
	itemId, err := strconv.ParseUint(c.Param("item"), 10, 64)
	if s.handleError(c, err) {
		return
	}

	second, err := strconv.ParseFloat(c.Query("second"), 64)
	if s.handleError(c, err) {
		return
	}

	logger.Infof("Splitting item %d at %f", itemId, second)
	changedItems, err := items.Split(s.db, itemId, second)
	if s.handleError(c, err) {
		return
	}

	for _, item := range changedItems {
		s.processor.EnqueueItemVideoMetadata(item.Id)
		s.processor.EnqueueItemCovers(item.Id)
		s.processor.EnqueueItemPreview(item.Id)
	}

	c.Status(http.StatusOK)
}

func (s *Server) makeHighlight(c *gin.Context) {
	itemId, err := strconv.ParseUint(c.Param("item"), 10, 64)
	if s.handleError(c, err) {
		return
	}

	startSecond, err := strconv.ParseFloat(c.Query("start"), 64)
	if s.handleError(c, err) {
		return
	}

	endSecond, err := strconv.ParseFloat(c.Query("end"), 64)
	if s.handleError(c, err) {
		return
	}

	highlightId, err := strconv.ParseUint(c.Query("highlight-id"), 10, 64)
	if s.handleError(c, err) {
		return
	}

	logger.Infof("Making highlight for item %d from %f to %f", itemId, startSecond, endSecond)
	highlightItem, err := items.MakeHighlight(s.db, itemId, startSecond, endSecond, highlightId)
	if s.handleError(c, err) {
		return
	}

	s.processor.EnqueueItemVideoMetadata(highlightItem.Id)
	s.processor.EnqueueItemCovers(highlightItem.Id)
	s.processor.EnqueueItemPreview(highlightItem.Id)
	c.Status(http.StatusOK)
}

func (s *Server) getSuggestionsForItem(c *gin.Context) {
	itemId, err := strconv.ParseUint(c.Param("item"), 10, 64)
	if s.handleError(c, err) {
		return
	}

	result, err := suggestions.GetSuggestionsForItem(s.db, s.db, itemId, 8)
	if s.handleError(c, err) {
		return
	}

	c.JSON(http.StatusOK, result)
}
