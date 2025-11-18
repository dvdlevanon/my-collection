package spectagger

import (
	"context"
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

type specTaggerDb interface {
	model.TagReaderWriter
	model.ItemReaderWriter
	model.TagAnnotationReaderWriter
}

func New(ctx context.Context, db specTaggerDb) (*Spectagger, error) {
	s, err := db.GetTag(ctx, special_tags.SpecTag)
	if err != nil {
		if err := db.CreateOrUpdateTag(ctx, special_tags.SpecTag); err != nil {
			return nil, err
		}
	} else {
		special_tags.SpecTag = s
	}

	return &Spectagger{
		db:             db,
		triggerChannel: make(chan bool),
	}, nil
}

type Spectagger struct {
	db             specTaggerDb
	triggerChannel chan bool
}

func (d *Spectagger) EnqueueSpecTagger() {
	d.triggerChannel <- true
}

func (d *Spectagger) Run(ctx context.Context) error {
	first := true
	ctx = utils.ContextWithSubject(ctx, "spectagger")

	for {
		select {
		case <-d.triggerChannel:
			d.runSpectagger(ctx)
		case <-time.After(1 * time.Minute):
			if first {
				d.runSpectagger(ctx)
			}

			first = false
		case <-time.After(60 * time.Minute):
			d.runSpectagger(ctx)
		case <-ctx.Done():
			return nil
		}
	}
}

func (d *Spectagger) runSpectagger(ctx context.Context) {
	logger.Infof("Spectagger started")
	if err := d.autoSpectag(ctx); err != nil {
		utils.LogError("Error in autoSpectag", err)
	}
	logger.Infof("Spectagger finished")
}

func (d *Spectagger) autoSpectag(ctx context.Context) error {
	allItems, err := d.db.GetAllItems(ctx)
	if err != nil {
		return err
	}

	categories, err := d.GetUserCategories(ctx)
	if err != nil {
		return err
	}

	tagTitleToId := make(map[string]uint64)

	for _, item := range *allItems {
		resolutionTags, err := getResolutionTags(ctx, d.db, &item)
		if err != nil {
			utils.LogError("Error getting resolution tags", err)
			continue
		}

		videoCodecTag, err := getVideoCodecTag(ctx, d.db, &item)
		if err != nil {
			utils.LogError("Error getting video codec tag", err)
			continue
		}

		audioCodecTag, err := getAudioCodecTag(ctx, d.db, &item)
		if err != nil {
			utils.LogError("Error getting audio codec tag", err)
			continue
		}

		durationTag, err := getDurationTag(ctx, d.db, &item)
		if err != nil {
			utils.LogError("Error getting duration tag", err)
			continue
		}

		typeTag, typeToRemove, err := getTypeTag(ctx, d.db, &item)
		if err != nil {
			utils.LogError("Error getting type tag", err)
			continue
		}

		categoryTagsToAdd, categoryTagsToRemove := getCategoryTags(ctx, d.db, categories, &item)
		tagsToAdd := append(categoryTagsToAdd, videoCodecTag, audioCodecTag, durationTag, typeTag)
		tagsToAdd = append(tagsToAdd, resolutionTags...)
		tagsToAdd = removeNils(tagsToAdd)

		if err := addTagsToItem(ctx, &tagTitleToId, d.db, &item, tagsToAdd); err != nil {
			utils.LogError("Error adding tags to item", err)
			continue
		}

		tagsToRemove := categoryTagsToRemove
		if typeToRemove != nil {
			tagsToRemove = append(tagsToRemove, typeToRemove)
		}

		if err := removeTagsFromItem(ctx, &tagTitleToId, d.db, &item, tagsToRemove); err != nil {
			utils.LogError("Error removing tags from item", err)
			continue
		}
	}

	return nil
}

func removeNils(tags []*model.Tag) []*model.Tag {
	result := make([]*model.Tag, 0)
	for _, item := range tags {
		if item != nil {
			result = append(result, item)
		}
	}

	return result
}

func (d *Spectagger) GetUserCategories(ctx context.Context) (*[]model.Tag, error) {
	categories, err := tags.GetCategories(ctx, d.db)
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

func getCategoryTags(ctx context.Context, tarw model.TagAnnotationReaderWriter, categories *[]model.Tag,
	item *model.Item) ([]*model.Tag, []*model.Tag) {
	categoryTagsToAdd := make([]*model.Tag, 0)
	categoryTagsToRemove := make([]*model.Tag, 0)

	if items.IsHighlight(item) || items.IsSplittedItem(item) {
		return categoryTagsToAdd, categoryTagsToRemove
	}

	categoriesExists := getCategoriesExists(categories, item)
	for i, category := range *categories {
		missing, err := getMissingFromCategoryTag(ctx, tarw, &category)
		if err != nil {
			utils.LogError("Error getting missing from category tag", err)
			continue
		}

		belong, err := getBelongToCategoryTag(ctx, tarw, &category)
		if err != nil {
			utils.LogError("Error getting belong to category tag", err)
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

func getMissingFromCategoryTag(ctx context.Context, tarw model.TagAnnotationReaderWriter, tag *model.Tag) (*model.Tag, error) {
	ta, err := tag_annotations.GetOrCreateTagAnnoation(ctx, tarw, &model.TagAnnotation{Title: "Categories"})
	if err != nil {
		return nil, err
	}

	return &model.Tag{
		ParentID:    &special_tags.SpecTag.Id,
		Title:       fmt.Sprintf("Missing %s", tag.Title),
		Annotations: []*model.TagAnnotation{ta},
	}, nil
}

func getBelongToCategoryTag(ctx context.Context, tarw model.TagAnnotationReaderWriter, tag *model.Tag) (*model.Tag, error) {
	ta, err := tag_annotations.GetOrCreateTagAnnoation(ctx, tarw, &model.TagAnnotation{Title: "Categories"})
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

func getResolutionTags(ctx context.Context, tarw model.TagAnnotationReaderWriter, item *model.Item) ([]*model.Tag, error) {
	if item.Width == 0 || item.Height == 0 {
		return nil, nil
	}

	ta, err := tag_annotations.GetOrCreateTagAnnoation(ctx, tarw, &model.TagAnnotation{Title: "Resolutions"})
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

func getVideoCodecTag(ctx context.Context, tarw model.TagAnnotationReaderWriter, item *model.Item) (*model.Tag, error) {
	if item.VideoCodecName == "" {
		return nil, nil
	}

	ta, err := tag_annotations.GetOrCreateTagAnnoation(ctx, tarw, &model.TagAnnotation{Title: "Video Codecs"})
	if err != nil {
		return nil, err
	}

	return &model.Tag{
		ParentID:    &special_tags.SpecTag.Id,
		Title:       item.VideoCodecName,
		Annotations: []*model.TagAnnotation{ta},
	}, nil
}

func getAudioCodecTag(ctx context.Context, tarw model.TagAnnotationReaderWriter, item *model.Item) (*model.Tag, error) {
	if item.AudioCodecName == "" {
		return nil, nil
	}

	ta, err := tag_annotations.GetOrCreateTagAnnoation(ctx, tarw, &model.TagAnnotation{Title: "Audio Codecs"})
	if err != nil {
		return nil, err
	}

	return &model.Tag{
		ParentID:    &special_tags.SpecTag.Id,
		Title:       item.AudioCodecName,
		Annotations: []*model.TagAnnotation{ta},
	}, nil
}

func getTypeTag(ctx context.Context, tarw model.TagAnnotationReaderWriter, item *model.Item) (*model.Tag, *model.Tag, error) {
	ta, err := tag_annotations.GetOrCreateTagAnnoation(ctx, tarw, &model.TagAnnotation{Title: "Type"})
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

func getDurationTag(ctx context.Context, tarw model.TagAnnotationReaderWriter, item *model.Item) (*model.Tag, error) {
	if item.DurationSeconds == 0 {
		return nil, nil
	}

	ta, err := tag_annotations.GetOrCreateTagAnnoation(ctx, tarw, &model.TagAnnotation{Title: "Duration"})
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

func getTagsWithId(ctx context.Context, tagTitleToId *map[string]uint64, trw model.TagReaderWriter,
	unsavedTags []*model.Tag) ([]*model.Tag, error) {
	savedTags := make([]*model.Tag, 0)
	for _, unsavedTag := range unsavedTags {
		id, ok := (*tagTitleToId)[unsavedTag.Title]
		if !ok {
			tag, err := tags.GetOrCreateTag(ctx, trw, unsavedTag)
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

func addTagsToItem(ctx context.Context, tagTitleToId *map[string]uint64, db specTaggerDb, item *model.Item, unsavedTags []*model.Tag) error {
	if len(unsavedTags) == 0 {
		return nil
	}

	savedTags, err := getTagsWithId(ctx, tagTitleToId, db, unsavedTags)
	if err != nil {
		return err
	}

	_, err = items.EnsureItemHaveTags(ctx, db, item, savedTags)
	return err
}

func removeTagsFromItem(ctx context.Context, tagTitleToId *map[string]uint64, db specTaggerDb,
	item *model.Item, unsavedTags []*model.Tag) error {
	if len(unsavedTags) == 0 {
		return nil
	}

	savedTags, err := getTagsWithId(ctx, tagTitleToId, db, unsavedTags)
	if err != nil {
		return err
	}

	return items.EnsureItemMissingTags(ctx, db, item, savedTags)
}
