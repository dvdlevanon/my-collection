package server

import (
	"encoding/json"
	"io"
	"my-collection/server/pkg/model"
	"my-collection/server/pkg/relativasor"
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
}
