package mixondemand

import (
	"context"
	"fmt"

	"my-collection/server/pkg/bl/special_tags"
	"my-collection/server/pkg/bl/tag_annotations"
	"my-collection/server/pkg/bl/tags"
	"my-collection/server/pkg/model"
	"my-collection/server/pkg/suggestions"
	"time"
)

type mixOnDemandDb interface {
	model.TagReaderWriter
	model.ItemReader
	model.TagAnnotationReaderWriter
}

func New(ctx context.Context, db mixOnDemandDb, mixOnDemandItemsCount int) (*MixOnDemand, error) {
	d, err := db.GetTag(ctx, special_tags.MixOnDemandTag)
	if err != nil {
		if err := db.CreateOrUpdateTag(ctx, special_tags.MixOnDemandTag); err != nil {
			return nil, err
		}
	} else {
		special_tags.MixOnDemandTag = d
	}

	return &MixOnDemand{
		db:                    db,
		mixOnDemandItemsCount: mixOnDemandItemsCount,
	}, nil
}

type MixOnDemand struct {
	db                    mixOnDemandDb
	mixOnDemandItemsCount int
}

func (d *MixOnDemand) GetCurrentTime() time.Time {
	return time.Now()
}

func (d *MixOnDemand) GenerateMixOnDemand(ctx context.Context, ctg model.CurrentTimeGetter, desc string, tags []model.Tag) (*model.Tag, error) {
	tag, err := prepareMixOnDemandTag(ctx, d.db, ctg, desc)
	if err != nil {
		return nil, err
	}

	result, err := suggestions.GetSuggestionsForTags(ctx, d.db, d.db, &tags, d.mixOnDemandItemsCount)
	if err != nil {
		return nil, err
	}

	tag.Items = result
	return tag, d.db.CreateOrUpdateTag(ctx, tag)
}

func getCurrentMixOnDemandTitle(desc string, ctg model.CurrentTimeGetter) string {
	return fmt.Sprintf("%s - %s", desc, ctg.GetCurrentTime().Format("02-Jan-2006 15:04:05"))
}

func getCurrentMixOnDemandAnnotation(ctg model.CurrentTimeGetter) string {
	return ctg.GetCurrentTime().Format("Jan-2006")
}

func prepareMixOnDemandTag(ctx context.Context, db mixOnDemandDb, ctg model.CurrentTimeGetter, desc string) (*model.Tag, error) {
	tag, err := tags.GetOrCreateChildTag(ctx, db, special_tags.MixOnDemandTag.Id, getCurrentMixOnDemandTitle(desc, ctg))
	if err != nil {
		return nil, err
	}

	_, err = tag_annotations.AddAnnotationToTag(ctx, db, db, tag.Id, model.TagAnnotation{
		Title: getCurrentMixOnDemandAnnotation(ctg),
	})
	if err != nil {
		return nil, err
	}

	return tag, err
}

func GetMixOnDemandTagId() uint64 {
	return special_tags.MixOnDemandTag.Id
}
