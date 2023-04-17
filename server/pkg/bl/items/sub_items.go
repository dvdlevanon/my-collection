package items

import (
	"fmt"
	"my-collection/server/pkg/model"

	"github.com/go-errors/errors"
)

func buildSubItem(item *model.Item, startPosition float64, endPosition float64) *model.Item {
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
	}
}

func GetMainItem(ir model.ItemReader, itemId uint64) (*model.Item, error) {
	item, err := ir.GetItem(itemId)
	if err != nil {
		return nil, err
	}

	for IsSubItem(item) {
		item, err = ir.GetItem(item.MainItemId)
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

func splitMain(iw model.ItemWriter, mainItem *model.Item, second float64) ([]*model.Item, error) {
	changedItems := make([]*model.Item, 0)
	sub1 := buildSubItem(mainItem, 0, second)
	sub2 := buildSubItem(mainItem, second, float64(mainItem.DurationSeconds))
	if err := iw.CreateOrUpdateItem(sub1); err != nil {
		return nil, err
	}
	if err := iw.CreateOrUpdateItem(sub2); err != nil {
		return nil, err
	}
	mainItem.SubItems = append(mainItem.SubItems, sub1, sub2)
	changedItems = append(mainItem.SubItems, sub1, sub2)
	return changedItems, nil
}

func shrinkAndSplit(iw model.ItemWriter, mainItem *model.Item, containedItem *model.Item, second float64) ([]*model.Item, error) {
	changedItems := make([]*model.Item, 0)
	sub := buildSubItem(mainItem, second, containedItem.EndPosition)
	if err := iw.CreateOrUpdateItem(sub); err != nil {
		return nil, err
	}
	containedItem.EndPosition = second
	containedItem.DurationSeconds = containedItem.EndPosition - containedItem.StartPosition
	if err := iw.UpdateItem(containedItem); err != nil {
		return nil, err
	}
	mainItem.SubItems = append(mainItem.SubItems, sub)
	changedItems = append(mainItem.SubItems, sub, containedItem)
	return changedItems, nil
}

func Split(irw model.ItemReaderWriter, itemId uint64, second float64) ([]*model.Item, error) {
	mainItem, err := GetMainItem(irw, itemId)
	if err != nil {
		return nil, err
	}

	containedItem, err := GetContainedSubItem(mainItem, second)
	if err != nil {
		return nil, err
	}

	var changedItems []*model.Item
	if IsSubItem(containedItem) {
		changedItems, err = shrinkAndSplit(irw, mainItem, containedItem, second)
	} else {
		changedItems, err = splitMain(irw, mainItem, second)
	}
	if err != nil {
		return nil, err
	}

	return changedItems, irw.UpdateItem(mainItem)
}

func IsSubItem(item *model.Item) bool {
	return item.MainItemId != nil
}
