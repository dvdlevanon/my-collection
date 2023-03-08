package db

import (
	"my-collection/server/pkg/model"
)

func (d *Database) CreateTagAnnotation(tagAnnotation *model.TagAnnotation) error {
	return d.create(tagAnnotation)
}

func (d *Database) RemoveTag(tagId uint64) error {
	return d.delete(model.Tag{Id: tagId})
}

func (d *Database) RemoveTagAnnotationFromTag(tagId uint64, annotationId uint64) error {
	return d.deleteAssociation(model.Tag{Id: tagId}, model.TagAnnotation{Id: annotationId}, "Annotations")
}

func (d *Database) GetTagAnnotation(conds ...interface{}) (*model.TagAnnotation, error) {
	tagAnnotation := &model.TagAnnotation{}
	err := d.handleError(d.db.Model(tagAnnotation).First(tagAnnotation, conds...).Error)
	return tagAnnotation, err
}

func (d *Database) GetTagAnnotations(tagId uint64) ([]model.TagAnnotation, error) {
	var annotations []model.TagAnnotation
	err := d.handleError(d.db.Model(&model.Tag{Id: tagId}).Association("Annotations").Find(&annotations))
	return annotations, err
}
