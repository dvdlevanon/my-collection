package spectagger

import (
	"context"
	"my-collection/server/pkg/model"
)

func newCachedTarw(db model.TagAnnotationReaderWriter) *cachedTarw {
	return &cachedTarw{
		db:    db,
		cache: make(map[string]*model.TagAnnotation),
	}
}

type cachedTarw struct {
	db    model.TagAnnotationReaderWriter
	cache map[string]*model.TagAnnotation
}

func (c *cachedTarw) GetTagAnnotation(ctx context.Context, conds ...interface{}) (*model.TagAnnotation, error) {
	if len(conds) == 1 {
		if tagAnnotation, ok := conds[0].(*model.TagAnnotation); ok && tagAnnotation != nil && tagAnnotation.Title != "" {
			cached, ok := c.cache[tagAnnotation.Title]
			if ok {
				return cached, nil
			}

			result, err := c.db.GetTagAnnotation(ctx, conds...)
			if err != nil {
				return nil, err
			}

			if result != nil {
				c.cache[tagAnnotation.Title] = result
			}
			return result, nil
		}
	}

	return c.db.GetTagAnnotation(ctx, conds...)
}

func (c *cachedTarw) GetTagAnnotations(ctx context.Context, tagId uint64) ([]model.TagAnnotation, error) {
	return c.db.GetTagAnnotations(ctx, tagId)
}

func (c *cachedTarw) CreateTagAnnotation(ctx context.Context, tagAnnotation *model.TagAnnotation) error {
	return c.db.CreateTagAnnotation(ctx, tagAnnotation)
}

func (c *cachedTarw) RemoveTag(ctx context.Context, tagId uint64) error {
	return c.db.RemoveTag(ctx, tagId)
}

func (c *cachedTarw) RemoveTagAnnotationFromTag(ctx context.Context, tagId uint64, annotationId uint64) error {
	return c.db.RemoveTagAnnotationFromTag(ctx, tagId, annotationId)
}
