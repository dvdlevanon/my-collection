package tags

import (
	"my-collection/server/pkg/model"

	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("tags")

func GetItems(ir model.ItemReader, tag *model.Tag) (*[]model.Item, error) {
	itemIds := make([]uint64, 0)
	for _, item := range tag.Items {
		itemIds = append(itemIds, item.Id)
	}

	items, err := ir.GetItems(itemIds)
	if err != nil {
		logger.Errorf("Error getting files of tag %t", err)
		return nil, err
	}

	return items, nil
}
