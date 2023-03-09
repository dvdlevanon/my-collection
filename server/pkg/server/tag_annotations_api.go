package server

import (
	"encoding/json"
	"io"
	"my-collection/server/pkg/bl/tag_annotations"
	"my-collection/server/pkg/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (s *Server) addAnnotationToTag(c *gin.Context) {
	tagId, err := strconv.ParseUint(c.Param("tag"), 10, 64)
	if s.handleError(c, err) {
		return
	}

	body, err := io.ReadAll(c.Request.Body)
	if s.handleError(c, err) {
		return
	}

	var annotation model.TagAnnotation
	if err = json.Unmarshal(body, &annotation); err != nil {
		s.handleError(c, err)
		return
	}

	annotationId, err := tag_annotations.AddAnnotationToTag(s.db, s.db, tagId, annotation)
	if s.handleError(c, err) {
		return
	}

	c.JSON(http.StatusOK, model.TagAnnotation{Id: annotationId})
}

func (s *Server) removeAnnotationFromTag(c *gin.Context) {
	tagId, err := strconv.ParseUint(c.Param("tag"), 10, 64)
	if s.handleError(c, err) {
		return
	}

	annotationId, err := strconv.ParseUint(c.Param("annotation-id"), 10, 64)
	if s.handleError(c, err) {
		return
	}

	if s.handleError(c, s.db.RemoveTagAnnotationFromTag(tagId, annotationId)) {
		return
	}

	c.Status(http.StatusOK)
}

func (s *Server) getTagAvailableAnnotations(c *gin.Context) {
	tagId, err := strconv.ParseUint(c.Param("tag"), 10, 64)
	if s.handleError(c, err) {
		return
	}

	availableAnnotations, err := tag_annotations.GetTagAvailableAnnotations(s.db, s.db, tagId)
	if s.handleError(c, err) {
		return
	}

	c.JSON(http.StatusOK, availableAnnotations)
}
