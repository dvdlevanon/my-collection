package backup

import (
	"encoding/json"
	"io"
	"my-collection/server/pkg/model"

	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("backup")

func Export(ir model.ItemReader, tr model.TagReader, w io.Writer) error {
	items, err := ir.GetAllItems()
	if err != nil {
		return err
	}

	tags, err := tr.GetAllTags()
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
