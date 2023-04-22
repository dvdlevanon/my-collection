package tag_annotations

import (
	"errors"
	"my-collection/server/pkg/model"

	"gorm.io/gorm"
)

func GetOrCreateTagAnnoation(tarw model.TagAnnotationReaderWriter, a *model.TagAnnotation) (*model.TagAnnotation, error) {
	annotation, err := tarw.GetTagAnnotation(a)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = tarw.CreateTagAnnotation(a)
		if err != nil {
			return a, err
		}

		annotation = a
	}

	return annotation, nil
}

func AddAnnotationToTag(trw model.TagReaderWriter, tarw model.TagAnnotationReaderWriter,
	tagId uint64, a model.TagAnnotation) (uint64, error) {
	annotation, err := GetOrCreateTagAnnoation(tarw, &a)
	if err != nil {
		return 0, err
	}

	tag, err := trw.GetTag(tagId)
	if err != nil {
		return 0, err
	}

	tag.Annotations = append(tag.Annotations, annotation)
	return annotation.Id, trw.CreateOrUpdateTag(tag)
}

func GetTagAvailableAnnotations(tr model.TagReader, tar model.TagAnnotationReader, tagId uint64) ([]model.TagAnnotation, error) {
	tag, err := tr.GetTag(tagId)
	if err != nil {
		return nil, err
	}

	availableAnnotations := make(map[uint64]model.TagAnnotation)
	allRelevantTags := make([]*model.Tag, 0)
	allRelevantTags = append(allRelevantTags, tag)
	allRelevantTags = append(allRelevantTags, tag.Children...)
	for _, child := range allRelevantTags {
		annotations, err := tar.GetTagAnnotations(child.Id)
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
