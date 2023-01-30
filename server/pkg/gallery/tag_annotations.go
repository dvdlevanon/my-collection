package gallery

import (
	"errors"
	"my-collection/server/pkg/model"

	"gorm.io/gorm"
)

func (g *Gallery) AddAnnotationToTag(tagId uint64, a model.TagAnnotation) (uint64, error) {
	annotation, err := g.GetTagAnnotation(a)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = g.CreateTagAnnotation(&a)
		if err != nil {
			return 0, err
		}

		annotation = &a
	}

	tag, err := g.GetTag(tagId)
	if err != nil {
		return 0, err
	}

	tag.Annotations = append(tag.Annotations, annotation)
	return annotation.Id, g.CreateOrUpdateTag(tag)
}

func (g *Gallery) GetTagAvailableAnnotations(tagId uint64) ([]model.TagAnnotation, error) {
	tag, err := g.GetTag(tagId)
	if err != nil {
		return nil, err
	}

	availableAnnotations := make(map[uint64]model.TagAnnotation)
	allRelevantTags := make([]*model.Tag, 0)
	allRelevantTags = append(allRelevantTags, tag)
	allRelevantTags = append(allRelevantTags, tag.Children...)
	for _, child := range allRelevantTags {
		annotations, err := g.GetTagAnnotations(child.Id)
		if err != nil {
			return nil, err
		}

		for _, annotation := range annotations {
			availableAnnotations[annotation.Id] = annotation
		}
	}

	result := make([]model.TagAnnotation, 0, len(availableAnnotations))
	for _, v := range availableAnnotations {
		result = append(result, v)
	}

	return result, nil
}
