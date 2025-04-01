package automix

import (
	"fmt"
	"my-collection/server/pkg/bl/items"
	"my-collection/server/pkg/bl/special_tags"
	"my-collection/server/pkg/bl/tag_annotations"
	"my-collection/server/pkg/bl/tags"
	"my-collection/server/pkg/model"
	"my-collection/server/pkg/utils"
	"time"
)

func New(trw model.TagReaderWriter, ir model.ItemReader,
	tarw model.TagAnnotationReaderWriter, dailyMixItemsCount int) (*Automix, error) {
	d, err := trw.GetTag(special_tags.DailymixTag)
	if err != nil {
		if err := trw.CreateOrUpdateTag(special_tags.DailymixTag); err != nil {
			return nil, err
		}
	} else {
		special_tags.DailymixTag = d
	}

	return &Automix{
		trw:                trw,
		ir:                 ir,
		tarw:               tarw,
		dailyMixItemsCount: dailyMixItemsCount,
	}, nil
}

type Automix struct {
	trw                model.TagReaderWriter
	ir                 model.ItemReader
	tarw               model.TagAnnotationReaderWriter
	dailyMixItemsCount int
}

func (d *Automix) Run() {
	for {
		if !isDailymixExists(d.trw, d) {
			if err := d.generateDailymix(d); err != nil {
				utils.LogError(err)
			}
		}

		time.Sleep(1 * time.Minute)
	}
}

func (d *Automix) GetCurrentTime() time.Time {
	return time.Now()
}

func (d *Automix) generateDailymix(ctg model.CurrentTimeGetter) error {
	tag, err := prepareDailymixTag(d.trw, d.tarw, ctg)
	if err != nil {
		return err
	}

	randomItems, err := items.GetRandomItems(d.ir, d.dailyMixItemsCount, func(item *model.Item) bool {
		isShortSubitem := items.IsSubItem(item) && item.DurationSeconds < 60*5
		return !items.IsHighlight(item) && !items.IsSplittedItem(item) && !isShortSubitem
	})

	if err != nil {
		return err
	}

	tag.Items = randomItems
	return d.trw.CreateOrUpdateTag(tag)
}

func isDailymixExists(trw model.TagReaderWriter, ctg model.CurrentTimeGetter) bool {
	_, err := tags.GetChildTag(trw, special_tags.DailymixTag.Id, getCurrentDailymixTitle(ctg))
	return err == nil
}

func getCurrentDailymixTitle(ctg model.CurrentTimeGetter) string {
	return fmt.Sprintf("Mix %s", ctg.GetCurrentTime().Format("02-Jan-2006"))
}

func getCurrentDailymixAnnotation(ctg model.CurrentTimeGetter) string {
	return ctg.GetCurrentTime().Format("Jan-2006")
}

func prepareDailymixTag(trw model.TagReaderWriter, tarw model.TagAnnotationReaderWriter,
	ctg model.CurrentTimeGetter) (*model.Tag, error) {
	tag, err := tags.GetOrCreateChildTag(trw, special_tags.DailymixTag.Id, getCurrentDailymixTitle(ctg))
	if err != nil {
		return nil, err
	}

	_, err = tag_annotations.AddAnnotationToTag(trw, tarw, tag.Id, model.TagAnnotation{
		Title: getCurrentDailymixAnnotation(ctg),
	})
	if err != nil {
		return nil, err
	}

	return tag, err
}

func GetDailymixTagId() uint64 {
	return special_tags.DailymixTag.Id
}
