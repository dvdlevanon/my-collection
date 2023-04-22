package spectagger

import (
	"fmt"
	"my-collection/server/pkg/bl/items"
	"my-collection/server/pkg/bl/tag_annotations"
	"my-collection/server/pkg/bl/tags"
	"my-collection/server/pkg/model"
	"my-collection/server/pkg/utils"
	"time"

	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("spectagger")

var specTag = &model.Tag{
	Title:          "Spec", // tags-utils.js
	ParentID:       nil,
	DisplayStyle:   "chip",
	DefaultSorting: "items-count",
}

func GetSpecTagId() uint64 {
	return specTag.Id
}

func New(trw model.TagReaderWriter, irw model.ItemReaderWriter,
	tarw model.TagAnnotationReaderWriter) (*Spectagger, error) {
	s, err := trw.GetTag(specTag)
	if err != nil {
		if err := trw.CreateOrUpdateTag(specTag); err != nil {
			return nil, err
		}
	} else {
		specTag = s
	}

	return &Spectagger{
		trw:  trw,
		irw:  irw,
		tarw: tarw,
	}, nil
}

type Spectagger struct {
	trw  model.TagReaderWriter
	irw  model.ItemReaderWriter
	tarw model.TagAnnotationReaderWriter
}

func (d *Spectagger) Run() {
	for {
		time.Sleep(1 * time.Minute)
		logger.Infof("Spectagger started")
		if err := d.autoSpectag(); err != nil {
			utils.LogError(err)
		}
		logger.Infof("Spectagger finished")
		time.Sleep(1 * time.Hour)
	}
}

func (d *Spectagger) autoSpectag() error {
	allItems, err := d.irw.GetAllItems()
	if err != nil {
		return err
	}

	tagTitleToId := make(map[string]uint64)

	for _, item := range *allItems {
		resolutionTag, err := getResolutionTag(d.tarw, &item)
		if err != nil {
			utils.LogError(err)
			continue
		}

		videoCodecTag, err := getVideoCodecTag(d.tarw, &item)
		if err != nil {
			utils.LogError(err)
			continue
		}

		audioCodecTag, err := getAudioCodecTag(d.tarw, &item)
		if err != nil {
			utils.LogError(err)
			continue
		}

		durationTag, err := getDurationTag(d.tarw, &item)
		if err != nil {
			utils.LogError(err)
			continue
		}

		if err := addTagsToItem(&tagTitleToId, d.trw, d.irw, &item, []*model.Tag{
			resolutionTag, videoCodecTag, audioCodecTag, durationTag}); err != nil {
			utils.LogError(err)
			continue
		}
	}

	return nil
}

func getResolutionTag(tarw model.TagAnnotationReaderWriter, item *model.Item) (*model.Tag, error) {
	ta, err := tag_annotations.GetOrCreateTagAnnoation(tarw, &model.TagAnnotation{Title: "Resolutions"})

	if err != nil {
		return nil, err
	}

	return &model.Tag{
		ParentID:    &specTag.Id,
		Title:       fmt.Sprintf("%d*%d", item.Width, item.Height),
		Annotations: []*model.TagAnnotation{ta},
	}, nil
}

func getVideoCodecTag(tarw model.TagAnnotationReaderWriter, item *model.Item) (*model.Tag, error) {
	ta, err := tag_annotations.GetOrCreateTagAnnoation(tarw, &model.TagAnnotation{Title: "Video Codecs"})

	if err != nil {
		return nil, err
	}

	return &model.Tag{
		ParentID:    &specTag.Id,
		Title:       item.VideoCodecName,
		Annotations: []*model.TagAnnotation{ta},
	}, nil
}

func getAudioCodecTag(tarw model.TagAnnotationReaderWriter, item *model.Item) (*model.Tag, error) {
	ta, err := tag_annotations.GetOrCreateTagAnnoation(tarw, &model.TagAnnotation{Title: "Audio Codecs"})

	if err != nil {
		return nil, err
	}

	return &model.Tag{
		ParentID:    &specTag.Id,
		Title:       item.AudioCodecName,
		Annotations: []*model.TagAnnotation{ta},
	}, nil
}

func getDurationTag(tarw model.TagAnnotationReaderWriter, item *model.Item) (*model.Tag, error) {
	ta, err := tag_annotations.GetOrCreateTagAnnoation(tarw, &model.TagAnnotation{Title: "Duration"})

	if err != nil {
		return nil, err
	}

	title := ""

	if item.DurationSeconds < (60 * 15) {
		title = "< 15 mintues"
	} else if item.DurationSeconds < (60 * 30) {
		title = "15-30 mintues"
	} else if item.DurationSeconds < (60 * 45) {
		title = "30-45 mintues"
	} else if item.DurationSeconds < (60 * 60) {
		title = "45-60 mintues"
	} else if item.DurationSeconds < (60 * 90) {
		title = "60-90 mintues"
	} else {
		title = "> 90 minutes"
	}

	return &model.Tag{
		ParentID:    &specTag.Id,
		Title:       title,
		Annotations: []*model.TagAnnotation{ta},
	}, nil
}

func addTagsToItem(tagTitleToId *map[string]uint64, trw model.TagReaderWriter, irw model.ItemReaderWriter,
	item *model.Item, unsavedTags []*model.Tag) error {
	savedTags := make([]*model.Tag, 0)
	for _, unsavedTag := range unsavedTags {
		id, ok := (*tagTitleToId)[unsavedTag.Title]
		if !ok {
			tag, err := tags.GetOrCreateTag(trw, unsavedTag)
			if err != nil {
				return err
			}
			id = tag.Id
			(*tagTitleToId)[unsavedTag.Title] = id
		}

		savedTags = append(savedTags, &model.Tag{Id: id})
	}

	return items.EnsureItemHaveTags(irw, item, savedTags)
}
