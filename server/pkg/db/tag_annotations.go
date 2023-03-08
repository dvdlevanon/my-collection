package db

import (
	"my-collection/server/pkg/model"

	"github.com/go-errors/errors"
)

func (d *Database) CreateTagAnnotation(tagAnnotation *model.TagAnnotation) error {
	return d.create(tagAnnotation)
}

func (d *Database) GetTagAnnotation(conds ...interface{}) (*model.TagAnnotation, error) {
	tagAnnotation := &model.TagAnnotation{}
	err := d.db.Model(tagAnnotation).First(tagAnnotation, conds...).Error
	return tagAnnotation, err
}

func (d *Database) RemoveTagAnnotationFromTag(tagId uint64, annotationId uint64) error {
	return d.db.Model(&model.Tag{Id: tagId}).Association("Annotations").Delete(model.TagAnnotation{Id: annotationId})
}

func (d *Database) GetTagAnnotations(tagId uint64) ([]model.TagAnnotation, error) {
	var annotations []model.TagAnnotation
	err := d.db.Model(&model.Tag{Id: tagId}).Association("Annotations").Find(&annotations)

	if err != nil {
		err = errors.Wrap(err, 0)
	}

	return annotations, err
}

func (d *Database) RemoveTag(tagId uint64) error {
	return d.db.Delete(model.Tag{Id: tagId}).Error
}
