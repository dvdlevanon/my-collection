package general_tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"my-collection/server/pkg/model"
	"my-collection/server/pkg/relativasor"
	"os"

	"github.com/go-errors/errors"
	"github.com/op/go-logging"
)

var fmLogger = logging.MustGetLogger("file-metadata")

type fileMetadataParams struct {
	ItemId uint64 `json:"id,omitempty"`
}

func MetadataDesc(id uint64, title string) string {
	return fmt.Sprintf("Metadata %s", title)
}

func MarshalFileMetadataParams(id uint64) (string, error) {
	p := fileMetadataParams{ItemId: id}
	res, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	return string(res), nil
}

func unmarshalFileMetadataParams(params string) (fileMetadataParams, error) {
	var p fileMetadataParams
	if err := json.Unmarshal([]byte(params), &p); err != nil {
		return p, err
	}
	return p, nil
}

func UpdateFileMetadata(ctx context.Context, irw model.ItemReaderWriter, params string) error {
	p, err := unmarshalFileMetadataParams(params)
	if err != nil {
		return err
	}

	return updateFileMetadata(ctx, irw, p)
}

func updateFileMetadata(ctx context.Context, irw model.ItemReaderWriter, p fileMetadataParams) error {
	item, err := irw.GetItem(ctx, p.ItemId)
	if err != nil {
		return err
	}

	if err := updateMetadata(item); err != nil {
		return err
	}

	return irw.UpdateItem(ctx, item)
}

func updateMetadata(item *model.Item) error {
	path := relativasor.GetAbsoluteFile(item.Url)
	fmLogger.Infof("Refreshing file metadata for item %d  [file: %s]", item.Id, path)

	file, err := os.Stat(path)
	if err != nil {
		return errors.Wrap(err, 1)
	}

	item.LastModified = file.ModTime().UnixMilli()
	item.FileSize = file.Size()
	return nil
}
