package automix

import (
	"context"
	"fmt"
	"my-collection/server/pkg/bl/items"
	"my-collection/server/pkg/bl/special_tags"
	"my-collection/server/pkg/bl/tag_annotations"
	"my-collection/server/pkg/bl/tags"
	"my-collection/server/pkg/model"
	"my-collection/server/pkg/utils"
	"time"
)

type autoMixDb interface {
	model.TagReaderWriter
	model.ItemReader
	model.TagAnnotationReaderWriter
}

func New(db autoMixDb, dailyMixItemsCount int) (*Automix, error) {
	d, err := db.GetTag(special_tags.DailymixTag)
	if err != nil {
		if err := db.CreateOrUpdateTag(special_tags.DailymixTag); err != nil {
			return nil, err
		}
	} else {
		special_tags.DailymixTag = d
	}

	return &Automix{
		db:                 db,
		dailyMixItemsCount: dailyMixItemsCount,
	}, nil
}

type Automix struct {
	db                 autoMixDb
	dailyMixItemsCount int
}

func (d *Automix) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(1 * time.Minute):
			if !isDailymixExists(d.db, d) {
				if err := d.generateDailymix(d); err != nil {
					utils.LogError("Error in generateDailymix", err)
				}
			}
		}
	}
}

func (d *Automix) GetCurrentTime() time.Time {
	return time.Now()
}

func (d *Automix) generateDailymix(ctg model.CurrentTimeGetter) error {
	tag, err := prepareDailymixTag(d.db, ctg)
	if err != nil {
		return err
	}

	randomItems, err := items.GetRandomItems(d.db, d.dailyMixItemsCount, func(item *model.Item) bool {
		isShortSubitem := items.IsSubItem(item) && item.DurationSeconds < 60*5
		return !items.IsHighlight(item) && !items.IsSplittedItem(item) && !isShortSubitem
	})

	if err != nil {
		return err
	}

	tag.Items = randomItems
	return d.db.CreateOrUpdateTag(tag)
}

func isDailymixExists(db autoMixDb, ctg model.CurrentTimeGetter) bool {
	_, err := tags.GetChildTag(db, special_tags.DailymixTag.Id, getCurrentDailymixTitle(ctg))
	return err == nil
}

func getCurrentDailymixTitle(ctg model.CurrentTimeGetter) string {
	return fmt.Sprintf("Mix %s", ctg.GetCurrentTime().Format("02-Jan-2006"))
}

func getCurrentDailymixAnnotation(ctg model.CurrentTimeGetter) string {
	return ctg.GetCurrentTime().Format("Jan-2006")
}

func prepareDailymixTag(db autoMixDb, ctg model.CurrentTimeGetter) (*model.Tag, error) {
	tag, err := tags.GetOrCreateChildTag(db, special_tags.DailymixTag.Id, getCurrentDailymixTitle(ctg))
	if err != nil {
		return nil, err
	}

	_, err = tag_annotations.AddAnnotationToTag(db, db, tag.Id, model.TagAnnotation{
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
