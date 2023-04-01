package automix

import (
	"fmt"
	"my-collection/server/pkg/bl/items"
	"my-collection/server/pkg/bl/tag_annotations"
	"my-collection/server/pkg/bl/tags"
	"my-collection/server/pkg/model"
	"my-collection/server/pkg/utils"
	"time"
)

const DAILYMIX_TAG_ID = uint64(343) // tags-util.js

var DailymixTag = model.Tag{
	Id:    DAILYMIX_TAG_ID,
	Title: "DailyMix",
}

func New(trw model.TagReaderWriter, ir model.ItemReader,
	tarw model.TagAnnotationReaderWriter, dailyMixItemsCount int) (*Automix, error) {
	_, err := trw.GetTag(DailymixTag)
	if err != nil {
		if err := trw.CreateOrUpdateTag(&DailymixTag); err != nil {
			return nil, err
		}
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

	randomItems, err := items.GetRandomItems(d.ir, d.dailyMixItemsCount)
	if err != nil {
		return err
	}

	tag.Items = randomItems
	return d.trw.CreateOrUpdateTag(tag)
}

func isDailymixExists(trw model.TagReaderWriter, ctg model.CurrentTimeGetter) bool {
	_, err := tags.GetChildTag(trw, DAILYMIX_TAG_ID, getCurrentDailymixTitle(ctg))
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
	tag, err := tags.GetOrCreateChildTag(trw, DAILYMIX_TAG_ID, getCurrentDailymixTitle(ctg))
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
