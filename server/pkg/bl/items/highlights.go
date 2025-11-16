package items

import (
	"context"
	"fmt"
	"my-collection/server/pkg/model"
)

var highlightsTag = &model.Tag{
	Title:    "Highlights", // tags-utils.js
	ParentID: nil,
}

func InitHighlights(ctx context.Context, trw model.TagReaderWriter) error {
	h, err := trw.GetTag(ctx, highlightsTag)
	if err != nil {
		if err := trw.CreateOrUpdateTag(ctx, highlightsTag); err != nil {
			return err
		}
	} else {
		highlightsTag = h
	}
	return nil
}

func buildHighlightUrl(origin string, startPosition float64, endPosition float64) string {
	return fmt.Sprintf("%s-%f-%f", origin, startPosition, endPosition)
}

func buildHighlight(item *model.Item, startPosition float64, endPosition float64, highlightId uint64) *model.Item {
	return &model.Item{
		Title:           item.Title,
		Origin:          buildHighlightUrl(item.Origin, startPosition, endPosition),
		Url:             item.Url,
		StartPosition:   startPosition,
		EndPosition:     endPosition,
		Width:           item.Width,
		Height:          item.Height,
		DurationSeconds: endPosition - startPosition,
		VideoCodecName:  item.VideoCodecName,
		AudioCodecName:  item.AudioCodecName,
		LastModified:    item.LastModified,
		PreviewMode:     PREVIEW_FROM_START_POSITION,
		Tags:            []*model.Tag{{Id: highlightId}},
	}
}

func MakeHighlight(ctx context.Context, irw model.ItemReaderWriter, itemId uint64, startPosition float64, endPosition float64, highlightId uint64) (*model.Item, error) {
	item, err := irw.GetItem(ctx, itemId)
	if err != nil {
		return nil, err
	}

	highlight := buildHighlight(item, startPosition, endPosition, highlightId)
	if err := irw.CreateOrUpdateItem(ctx, highlight); err != nil {
		return nil, err
	}

	item.Highlights = append(item.Highlights, highlight)
	if err := irw.UpdateItem(ctx, item); err != nil {
		return nil, err
	}

	return highlight, nil
}

func IsHighlight(item *model.Item) bool {
	return item.HighlightParentItemId != nil
}

func GetHighlightsTagId() uint64 {
	return highlightsTag.Id
}
