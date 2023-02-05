package gallery

import (
	"encoding/json"
	"io"
	"my-collection/server/pkg/model"
)

func (g *Gallery) Export(w io.Writer) error {
	items, err := g.GetAllItems()
	if err != nil {
		return err
	}

	tags, err := g.GetAllTags()
	if err != nil {
		return err
	}

	logger.Infof("Exporting %d items and %d tags", len(*items), len(*tags))

	itemsAndTags := model.ItemsAndTags{
		Items: *items,
		Tags:  *tags,
	}

	jsonBytes, err := json.Marshal(itemsAndTags)
	if err != nil {
		return err
	}

	_, err = w.Write(jsonBytes)
	return err
}
