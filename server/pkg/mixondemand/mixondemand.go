package mixondemand

import (
	"fmt"

	"my-collection/server/pkg/bl/special_tags"
	"my-collection/server/pkg/bl/tag_annotations"
	"my-collection/server/pkg/bl/tags"
	"my-collection/server/pkg/model"
	"my-collection/server/pkg/suggestions"
	"time"
)

func New(trw model.TagReaderWriter, ir model.ItemReader,
	tarw model.TagAnnotationReaderWriter, mixOnDemandItemsCount int) (*MixOnDemand, error) {
	d, err := trw.GetTag(special_tags.MixOnDemandTag)
	if err != nil {
		if err := trw.CreateOrUpdateTag(special_tags.MixOnDemandTag); err != nil {
			return nil, err
		}
	} else {
		special_tags.MixOnDemandTag = d
	}

	return &MixOnDemand{
		trw:                   trw,
		ir:                    ir,
		tarw:                  tarw,
		mixOnDemandItemsCount: mixOnDemandItemsCount,
	}, nil
}

type MixOnDemand struct {
	trw                   model.TagReaderWriter
	ir                    model.ItemReader
	tarw                  model.TagAnnotationReaderWriter
	mixOnDemandItemsCount int
}

func (d *MixOnDemand) GetCurrentTime() time.Time {
	return time.Now()
}

func (d *MixOnDemand) GenerateMixOnDemand(ctg model.CurrentTimeGetter, desc string, tags []model.Tag) (*model.Tag, error) {
	tag, err := prepareMixOnDemandTag(d.trw, d.tarw, ctg, desc)
	if err != nil {
		return nil, err
	}

	result, err := suggestions.GetSuggestionsForTags(d.ir, d.trw, &tags, d.mixOnDemandItemsCount)
	if err != nil {
		return nil, err
	}

	tag.Items = result
	return tag, d.trw.CreateOrUpdateTag(tag)
}

func getCurrentMixOnDemandTitle(desc string, ctg model.CurrentTimeGetter) string {
	return fmt.Sprintf("%s - %s", desc, ctg.GetCurrentTime().Format("02-Jan-2006 15:04:05"))
}

func getCurrentMixOnDemandAnnotation(ctg model.CurrentTimeGetter) string {
	return ctg.GetCurrentTime().Format("Jan-2006")
}

func prepareMixOnDemandTag(trw model.TagReaderWriter, tarw model.TagAnnotationReaderWriter,
	ctg model.CurrentTimeGetter, desc string) (*model.Tag, error) {
	tag, err := tags.GetOrCreateChildTag(trw, special_tags.MixOnDemandTag.Id, getCurrentMixOnDemandTitle(desc, ctg))
	if err != nil {
		return nil, err
	}

	_, err = tag_annotations.AddAnnotationToTag(trw, tarw, tag.Id, model.TagAnnotation{
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
