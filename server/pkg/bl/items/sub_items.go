package items

import (
	"context"
	"fmt"
	"my-collection/server/pkg/model"

	"github.com/go-errors/errors"
)

func buildSubItemOrigin(origin string, startPosition float64, endPosition float64) string {
	return fmt.Sprintf("%s-%f-%f", origin, startPosition, endPosition)
}

func buildSubItem(item *model.Item, startPosition float64, endPosition float64) *model.Item {
	return &model.Item{
		Title:           item.Title,
		Origin:          buildSubItemOrigin(item.Origin, startPosition, endPosition),
		Url:             item.Url,
		StartPosition:   startPosition,
		EndPosition:     endPosition,
		Width:           item.Width,
		Height:          item.Height,
		DurationSeconds: endPosition - startPosition,
		VideoCodecName:  item.VideoCodecName,
		AudioCodecName:  item.AudioCodecName,
		LastModified:    item.LastModified,
	}
}

func GetMainItem(ctx context.Context, ir model.ItemReader, itemId uint64) (*model.Item, error) {
	item, err := ir.GetItem(ctx, itemId)
	if err != nil {
		return nil, err
	}

	for IsSubItem(item) {
		item, err = ir.GetItem(ctx, item.MainItemId)
		if err != nil {
			return nil, err
		}
	}

	return item, nil
}

func GetContainedSubItem(mainItem *model.Item, second float64) (*model.Item, error) {
	if len(mainItem.SubItems) == 0 {
		return mainItem, nil
	}

	for _, subItem := range mainItem.SubItems {
		if subItem.StartPosition <= second && subItem.EndPosition >= second {
			return subItem, nil
		}
	}

	return nil, errors.Errorf("sub-item at second %f not found in %v", second, mainItem)
}

func splitMain(ctx context.Context, iw model.ItemWriter, mainItem *model.Item, second float64) ([]*model.Item, error) {
	sub1 := buildSubItem(mainItem, 0, second)
	sub2 := buildSubItem(mainItem, second, float64(mainItem.DurationSeconds))
	if err := iw.CreateOrUpdateItem(ctx, sub1); err != nil {
		return nil, err
	}
	if err := iw.CreateOrUpdateItem(ctx, sub2); err != nil {
		return nil, err
	}
	mainItem.SubItems = append(mainItem.SubItems, sub1, sub2)
	changedItems := []*model.Item{sub1, sub2}
	return changedItems, nil
}

func shrinkAndSplit(ctx context.Context, iw model.ItemWriter, mainItem *model.Item, containedItem *model.Item, second float64) ([]*model.Item, error) {
	sub := buildSubItem(mainItem, second, containedItem.EndPosition)
	if err := iw.CreateOrUpdateItem(ctx, sub); err != nil {
		return nil, err
	}
	containedItem.EndPosition = second
	containedItem.DurationSeconds = containedItem.EndPosition - containedItem.StartPosition
	if err := iw.UpdateItem(ctx, containedItem); err != nil {
		return nil, err
	}
	mainItem.SubItems = append(mainItem.SubItems, sub)
	changedItems := []*model.Item{sub, containedItem}
	return changedItems, nil
}

func Split(ctx context.Context, irw model.ItemReaderWriter, itemId uint64, second float64) ([]*model.Item, error) {
	mainItem, err := GetMainItem(ctx, irw, itemId)
	if err != nil {
		return nil, err
	}

	containedItem, err := GetContainedSubItem(mainItem, second)
	if err != nil {
		return nil, err
	}

	var changedItems []*model.Item
	if IsSubItem(containedItem) {
		changedItems, err = shrinkAndSplit(ctx, irw, mainItem, containedItem, second)
	} else {
		changedItems, err = splitMain(ctx, irw, mainItem, second)
	}
	if err != nil {
		return nil, err
	}

	return changedItems, irw.UpdateItem(ctx, mainItem)
}

func IsSubItem(item *model.Item) bool {
	return item.MainItemId != nil
}

func IsSplittedItem(item *model.Item) bool {
	return len(item.SubItems) > 0
}
