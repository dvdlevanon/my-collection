package gallery

import (
	"my-collection/server/pkg/db"
	"my-collection/server/pkg/relativasor"
	"my-collection/server/pkg/storage"

	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("gallery")

type Gallery struct {
	*db.Database
	relativasor          *relativasor.PathRelativasor
	storage              *storage.Storage
	CoversCount          int
	PreviewSceneCount    int
	PreviewSceneDuration int
	AutomaticProcessing  bool
	TrustFileExtenssion  bool
}

func New(db *db.Database, storage *storage.Storage, relativasor *relativasor.PathRelativasor) *Gallery {
	return &Gallery{
		Database:             db,
		storage:              storage,
		relativasor:          relativasor,
		CoversCount:          3,
		PreviewSceneCount:    4,
		PreviewSceneDuration: 3,
		AutomaticProcessing:  false,
		TrustFileExtenssion:  true,
	}
}
