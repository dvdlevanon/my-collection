package items

import (
	"fmt"
	"my-collection/server/pkg/model"
)

func buildHighlight(item *model.Item, startPosition float64, endPosition float64) *model.Item {
	return &model.Item{
		Title:           item.Title,
		Origin:          fmt.Sprintf("%s-%f-%f", item.Origin, startPosition, endPosition),
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
	}
}

func MakeHighlight(irw model.ItemReaderWriter, itemId uint64, startPosition float64, endPosition float64) (*model.Item, error) {
	item, err := irw.GetItem(itemId)
	if err != nil {
		return nil, err
	}

	highlight := buildHighlight(item, startPosition, endPosition)
	if err := irw.CreateOrUpdateItem(highlight); err != nil {
		return nil, err
	}

	item.Highlights = append(item.Highlights, highlight)
	if err := irw.UpdateItem(item); err != nil {
		return nil, err
	}

	return highlight, nil
}

func IsHighlight(item *model.Item) bool {
	return item.HighlightParentItemId != nil
}
