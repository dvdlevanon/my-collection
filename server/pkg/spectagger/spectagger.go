package spectagger

import (
	"fmt"
	"my-collection/server/pkg/automix"
	"my-collection/server/pkg/bl/directories"
	"my-collection/server/pkg/bl/items"
	"my-collection/server/pkg/bl/special_tags"
	"my-collection/server/pkg/bl/tag_annotations"
	"my-collection/server/pkg/bl/tags"
	"my-collection/server/pkg/mixondemand"
	"my-collection/server/pkg/model"
	"my-collection/server/pkg/utils"
	"time"

	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("spectagger")

func GetSpecTagId() uint64 {
	return special_tags.SpecTag.Id
}

func New(trw model.TagReaderWriter, irw model.ItemReaderWriter,
	tarw model.TagAnnotationReaderWriter) (*Spectagger, error) {
	s, err := trw.GetTag(special_tags.SpecTag)
	if err != nil {
		if err := trw.CreateOrUpdateTag(special_tags.SpecTag); err != nil {
			return nil, err
		}
	} else {
		special_tags.SpecTag = s
	}

	return &Spectagger{
		trw:            trw,
		irw:            irw,
		tarw:           tarw,
		triggerChannel: make(chan bool),
	}, nil
}

type Spectagger struct {
	trw            model.TagReaderWriter
	irw            model.ItemReaderWriter
	tarw           model.TagAnnotationReaderWriter
	triggerChannel chan bool
}

func (d *Spectagger) Trigger() {
	d.triggerChannel <- true
}

func (d *Spectagger) Run() {
	first := true

	for {
		select {
		case <-d.triggerChannel:
			d.runSpectagger()
		case <-time.After(1 * time.Minute):
			if first {
				d.runSpectagger()
			}

			first = false
		case <-time.After(60 * time.Minute):
			d.runSpectagger()
		}
	}
}

func (d *Spectagger) runSpectagger() {
	logger.Infof("Spectagger started")
	if err := d.autoSpectag(); err != nil {
		utils.LogError(err)
	}
	logger.Infof("Spectagger finished")
}

func (d *Spectagger) autoSpectag() error {
	allItems, err := d.irw.GetAllItems()
	if err != nil {
		return err
	}

	categories, err := d.GetUserCategories()
	if err != nil {
		return err
	}

	tagTitleToId := make(map[string]uint64)

	for _, item := range *allItems {
		resolutionTags, err := getResolutionTags(d.tarw, &item)
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

		typeTag, typeToRemove, err := getTypeTag(d.tarw, &item)
		if err != nil {
			utils.LogError(err)
			continue
		}

		categoryTagsToAdd, categoryTagsToRemove := getCategoryTags(d.tarw, categories, &item)
		if err != nil {
			utils.LogError(err)
			continue
		}

		tagsToAdd := append(categoryTagsToAdd, videoCodecTag, audioCodecTag, durationTag, typeTag)
		tagsToAdd = append(tagsToAdd, resolutionTags...)

		if err := addTagsToItem(&tagTitleToId, d.trw, d.irw, &item, tagsToAdd); err != nil {
			utils.LogError(err)
			continue
		}

		tagsToRemove := categoryTagsToRemove
		if typeToRemove != nil {
			tagsToRemove = append(tagsToRemove, typeToRemove)
		}

		if err := removeTagsFromItem(&tagTitleToId, d.trw, d.irw, &item, tagsToRemove); err != nil {
			utils.LogError(err)
			continue
		}
	}

	return nil
}

func (d *Spectagger) GetUserCategories() (*[]model.Tag, error) {
	categories, err := tags.GetCategories(d.trw)
	if err != nil {
		return nil, err
	}

	userCategories := make([]model.Tag, 0)
	for _, category := range *categories {
		if category.Id == directories.GetDirectoriesTagId() ||
			category.Id == automix.GetDailymixTagId() ||
			category.Id == mixondemand.GetMixOnDemandTagId() ||
			category.Id == GetSpecTagId() ||
			category.Id == items.GetHighlightsTagId() {
			continue
		}

		userCategories = append(userCategories, category)
	}

	return &userCategories, nil
}

func getCategoryTags(tarw model.TagAnnotationReaderWriter, categories *[]model.Tag,
	item *model.Item) ([]*model.Tag, []*model.Tag) {
	categoryTagsToAdd := make([]*model.Tag, 0)
	categoryTagsToRemove := make([]*model.Tag, 0)

	if items.IsHighlight(item) || items.IsSplittedItem(item) {
		return categoryTagsToAdd, categoryTagsToRemove
	}

	categoriesExists := getCategoriesExists(categories, item)
	for i, category := range *categories {
		missing, err := getMissingFromCategoryTag(tarw, &category)
		if err != nil {
			utils.LogError(err)
			continue
		}

		belong, err := getBelongToCategoryTag(tarw, &category)
		if err != nil {
			utils.LogError(err)
			continue
		}

		if categoriesExists[i] {
			categoryTagsToRemove = append(categoryTagsToRemove, missing)
			categoryTagsToAdd = append(categoryTagsToAdd, belong)
		} else {
			categoryTagsToRemove = append(categoryTagsToRemove, belong)
			categoryTagsToAdd = append(categoryTagsToAdd, missing)
		}
	}

	return categoryTagsToAdd, categoryTagsToRemove
}

func getMissingFromCategoryTag(tarw model.TagAnnotationReaderWriter, tag *model.Tag) (*model.Tag, error) {
	ta, err := tag_annotations.GetOrCreateTagAnnoation(tarw, &model.TagAnnotation{Title: "Categories"})
	if err != nil {
		return nil, err
	}

	return &model.Tag{
		ParentID:    &special_tags.SpecTag.Id,
		Title:       fmt.Sprintf("Missing %s", tag.Title),
		Annotations: []*model.TagAnnotation{ta},
	}, nil
}

func getBelongToCategoryTag(tarw model.TagAnnotationReaderWriter, tag *model.Tag) (*model.Tag, error) {
	ta, err := tag_annotations.GetOrCreateTagAnnoation(tarw, &model.TagAnnotation{Title: "Categories"})
	if err != nil {
		return nil, err
	}

	return &model.Tag{
		ParentID:    &special_tags.SpecTag.Id,
		Title:       fmt.Sprintf("Has %s", tag.Title),
		Annotations: []*model.TagAnnotation{ta},
	}, nil
}

func getCategoriesExists(categories *[]model.Tag, item *model.Item) []bool {
	categoriesExists := make([]bool, len(*categories))
	for _, tag := range item.Tags {
		for i, category := range *categories {
			if tags.IsBelongToCategory(tag, &category) {
				categoriesExists[i] = true
				break
			}
		}
	}

	return categoriesExists
}

func getResolutionTags(tarw model.TagAnnotationReaderWriter, item *model.Item) ([]*model.Tag, error) {
	ta, err := tag_annotations.GetOrCreateTagAnnoation(tarw, &model.TagAnnotation{Title: "Resolutions"})
	if err != nil {
		return nil, err
	}

	var tags []*model.Tag

	resolutionTag := &model.Tag{
		ParentID:    &special_tags.SpecTag.Id, // Assuming specTag is defined elsewhere
		Title:       fmt.Sprintf("%d*%d", item.Width, item.Height),
		Annotations: []*model.TagAnnotation{ta},
	}
	tags = append(tags, resolutionTag)

	widthTag := &model.Tag{
		ParentID:    &special_tags.SpecTag.Id,
		Title:       fmt.Sprintf("Width: %d", item.Width),
		Annotations: []*model.TagAnnotation{ta},
	}
	tags = append(tags, widthTag)

	heightTag := &model.Tag{
		ParentID:    &special_tags.SpecTag.Id,
		Title:       fmt.Sprintf("Height: %d", item.Height),
		Annotations: []*model.TagAnnotation{ta},
	}
	tags = append(tags, heightTag)

	return tags, nil
}

func getVideoCodecTag(tarw model.TagAnnotationReaderWriter, item *model.Item) (*model.Tag, error) {
	ta, err := tag_annotations.GetOrCreateTagAnnoation(tarw, &model.TagAnnotation{Title: "Video Codecs"})
	if err != nil {
		return nil, err
	}

	return &model.Tag{
		ParentID:    &special_tags.SpecTag.Id,
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
		ParentID:    &special_tags.SpecTag.Id,
		Title:       item.AudioCodecName,
		Annotations: []*model.TagAnnotation{ta},
	}, nil
}

func getTypeTag(tarw model.TagAnnotationReaderWriter, item *model.Item) (*model.Tag, *model.Tag, error) {
	ta, err := tag_annotations.GetOrCreateTagAnnoation(tarw, &model.TagAnnotation{Title: "Type"})
	if err != nil {
		return nil, nil, err
	}

	regularTag := &model.Tag{
		ParentID:    &special_tags.SpecTag.Id,
		Title:       "Regular",
		Annotations: []*model.TagAnnotation{ta},
	}
	var tagToRemove *model.Tag

	title := ""
	if items.IsSplittedItem(item) {
		title = "Splitted"
		tagToRemove = regularTag
	} else if items.IsSubItem(item) {
		title = "Sub Item"
		tagToRemove = regularTag
	} else if items.IsHighlight(item) {
		title = "Hightlight"
		tagToRemove = regularTag
	} else {
		title = "Regular"
	}

	return &model.Tag{
		ParentID:    &special_tags.SpecTag.Id,
		Title:       title,
		Annotations: []*model.TagAnnotation{ta},
	}, tagToRemove, nil
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
	} else if item.DurationSeconds < (60 * 120) {
		title = "90-120 mintues"
	} else {
		title = "> 120 minutes"
	}

	return &model.Tag{
		ParentID:    &special_tags.SpecTag.Id,
		Title:       title,
		Annotations: []*model.TagAnnotation{ta},
	}, nil
}

func getTagsWithId(tagTitleToId *map[string]uint64, trw model.TagReaderWriter,
	unsavedTags []*model.Tag) ([]*model.Tag, error) {
	savedTags := make([]*model.Tag, 0)
	for _, unsavedTag := range unsavedTags {
		id, ok := (*tagTitleToId)[unsavedTag.Title]
		if !ok {
			tag, err := tags.GetOrCreateTag(trw, unsavedTag)
			if err != nil {
				return nil, err
			}
			id = tag.Id
			(*tagTitleToId)[unsavedTag.Title] = id
		}

		savedTags = append(savedTags, &model.Tag{Id: id})
	}

	return savedTags, nil
}

func addTagsToItem(tagTitleToId *map[string]uint64, trw model.TagReaderWriter, irw model.ItemReaderWriter,
	item *model.Item, unsavedTags []*model.Tag) error {
	if len(unsavedTags) == 0 {
		return nil
	}

	savedTags, err := getTagsWithId(tagTitleToId, trw, unsavedTags)
	if err != nil {
		return err
	}

	return items.EnsureItemHaveTags(irw, item, savedTags)
}

func removeTagsFromItem(tagTitleToId *map[string]uint64, trw model.TagReaderWriter, irw model.ItemReaderWriter,
	item *model.Item, unsavedTags []*model.Tag) error {
	if len(unsavedTags) == 0 {
		return nil
	}

	savedTags, err := getTagsWithId(tagTitleToId, trw, unsavedTags)
	if err != nil {
		return err
	}

	return items.EnsureItemMissingTags(irw, item, savedTags)
}
