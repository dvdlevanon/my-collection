package settings

import (
	"my-collection/server/pkg/model"
)

var Settings model.Settings

func Init(settings model.Settings) {
	Settings = settings
}
