package db

import (
	"context"
	"my-collection/server/pkg/model"
)

func (d *databaseImpl) CreateTagAnnotation(ctx context.Context, tagAnnotation *model.TagAnnotation) error {
	return d.create(ctx, tagAnnotation)
}

func (d *databaseImpl) RemoveTag(ctx context.Context, tagId uint64) error {
	return d.delete(ctx, model.Tag{Id: tagId})
}

func (d *databaseImpl) RemoveTagAnnotationFromTag(ctx context.Context, tagId uint64, annotationId uint64) error {
	return d.deleteAssociation(ctx, model.Tag{Id: tagId}, model.TagAnnotation{Id: annotationId}, "Annotations")
}

func (d *databaseImpl) GetTagAnnotation(ctx context.Context, conds ...interface{}) (*model.TagAnnotation, error) {
	tagAnnotation := &model.TagAnnotation{}
	err := d.handleError(d.db.WithContext(ctx).Model(tagAnnotation).First(tagAnnotation, conds...).Error)
	return tagAnnotation, err
}

func (d *databaseImpl) GetTagAnnotations(ctx context.Context, tagId uint64) ([]model.TagAnnotation, error) {
	var annotations []model.TagAnnotation
	err := d.handleError(d.db.WithContext(ctx).Model(&model.Tag{Id: tagId}).Association("Annotations").Find(&annotations))
	return annotations, err
}
