package server

import (
	"encoding/json"
	"io"
	"my-collection/server/pkg/automix"
	"my-collection/server/pkg/bl/directories"
	"my-collection/server/pkg/bl/items"
	"my-collection/server/pkg/bl/tags"
	"my-collection/server/pkg/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-errors/errors"
)

func (s *Server) createTag(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if s.handleError(c, err) {
		return
	}

	var tag model.Tag
	if s.handleError(c, json.Unmarshal(body, &tag)) {
		return
	}

	if s.handleError(c, s.db.CreateOrUpdateTag(&tag)) {
		return
	}

	c.JSON(http.StatusOK, model.Tag{Id: tag.Id})
}

func (s *Server) updateTag(c *gin.Context) {
	tagId, err := strconv.ParseUint(c.Param("tag"), 10, 64)
	if s.handleError(c, err) {
		return
	}

	body, err := io.ReadAll(c.Request.Body)
	if s.handleError(c, err) {
		return
	}

	var tag model.Tag
	if s.handleError(c, json.Unmarshal(body, &tag)) {
		return
	}

	if tag.Id != 0 && tag.Id != tagId {
		s.handleError(c, errors.Errorf("Mismatch IDs %d != %d", tag.Id, tagId))
		return
	}

	tag.Id = tagId
	if s.handleError(c, s.db.UpdateTag(&tag)) {
		return
	}

	c.Status(http.StatusOK)
}

func (s *Server) getTag(c *gin.Context) {
	tagId, err := strconv.ParseUint(c.Param("tag"), 10, 64)
	if s.handleError(c, err) {
		return
	}

	tag, err := s.db.GetTag(tagId)
	if s.handleError(c, err) {
		return
	}

	c.JSON(http.StatusOK, tag)
}

func (s *Server) removeTag(c *gin.Context) {
	tagId, err := strconv.ParseUint(c.Param("tag"), 10, 64)
	if s.handleError(c, err) {
		return
	}

	if s.handleError(c, s.db.RemoveTag(tagId)) {
		return
	}

	c.Status(http.StatusOK)
}

func (s *Server) getTags(c *gin.Context) {
	tags, err := s.db.GetAllTags()
	if s.handleError(c, err) {
		return
	}

	logger.Infof("Get tags return %d tags", len(*tags))
	c.JSON(http.StatusOK, tags)
}

func (s *Server) getSpecialTags(c *gin.Context) {
	tags, err := s.db.GetTagsWithoutChildren(
		directories.GetDirectoriesTagId(),
		automix.GetDailymixTagId(),
		items.GetHighlightsTagId())

	if s.handleError(c, err) {
		return
	}

	logger.Infof("Get tags return %d tags", len(*tags))
	c.JSON(http.StatusOK, tags)
}

func (s *Server) autoImage(c *gin.Context) {
	tagId, err := strconv.ParseUint(c.Param("tag"), 10, 64)
	if s.handleError(c, err) {
		return
	}

	tag, err := s.db.GetTag(tagId)
	if s.handleError(c, err) {
		return
	}

	body, err := io.ReadAll(c.Request.Body)
	if s.handleError(c, err) {
		return
	}

	var fileUrl model.FileUrl
	if err = json.Unmarshal(body, &fileUrl); err != nil {
		s.handleError(c, err)
		return
	}

	if s.handleError(c, tags.AutoImageChildren(s.storage, s.db, s.db, tag, fileUrl.Url)) {
		return
	}

	c.JSON(http.StatusOK, tag)
}
