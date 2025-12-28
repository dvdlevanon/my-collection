package backup

import (
	"context"
	"encoding/json"
	"io"
	"my-collection/server/pkg/model"
	"my-collection/server/pkg/utils"
)

var logger = utils.MustGetLogger("backup")

func Export(ctx context.Context, ir model.ItemReader, tr model.TagReader, w io.Writer) error {
	items, err := ir.GetAllItems(ctx)
	if err != nil {
		return err
	}

	tags, err := tr.GetAllTags(ctx)
	if err != nil {
		return err
	}

	logger.Infof(ctx, "Exporting %d items and %d tags", len(*items), len(*tags))

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
