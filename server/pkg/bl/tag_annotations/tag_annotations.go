package tag_annotations

import (
	"context"
	"errors"
	"my-collection/server/pkg/model"

	"gorm.io/gorm"
)

func GetOrCreateTagAnnoation(ctx context.Context, tarw model.TagAnnotationReaderWriter, a *model.TagAnnotation) (*model.TagAnnotation, error) {
	annotation, err := tarw.GetTagAnnotation(ctx, a)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = tarw.CreateTagAnnotation(ctx, a)
		if err != nil {
			return a, err
		}

		annotation = a
	}

	return annotation, nil
}

func AddAnnotationToTag(ctx context.Context, trw model.TagReaderWriter, tarw model.TagAnnotationReaderWriter,
	tagId uint64, a model.TagAnnotation) (uint64, error) {
	annotation, err := GetOrCreateTagAnnoation(ctx, tarw, &a)
	if err != nil {
		return 0, err
	}

	tag, err := trw.GetTag(ctx, tagId)
	if err != nil {
		return 0, err
	}

	tag.Annotations = append(tag.Annotations, annotation)
	return annotation.Id, trw.CreateOrUpdateTag(ctx, tag)
}

func GetTagAvailableAnnotations(ctx context.Context, tr model.TagReader, tar model.TagAnnotationReader, tagId uint64) ([]model.TagAnnotation, error) {
	tag, err := tr.GetTag(ctx, tagId)
	if err != nil {
		return nil, err
	}

	availableAnnotations := make(map[uint64]model.TagAnnotation)
	allRelevantTags := make([]*model.Tag, 0)
	allRelevantTags = append(allRelevantTags, tag)
	allRelevantTags = append(allRelevantTags, tag.Children...)
	for _, child := range allRelevantTags {
		annotations, err := tar.GetTagAnnotations(ctx, child.Id)
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
