package app

import (
	"my-collection/server/pkg/db"
	"my-collection/server/pkg/fssync"
	"my-collection/server/pkg/itemsoptimizer"
	"my-collection/server/pkg/mixondemand"
	"my-collection/server/pkg/model"
	"my-collection/server/pkg/opensubtitles"
	"my-collection/server/pkg/processor"
	"my-collection/server/pkg/server/fs"
	"my-collection/server/pkg/server/items"
	"my-collection/server/pkg/server/management"
	storageHandler "my-collection/server/pkg/server/storage"
	"my-collection/server/pkg/server/subtitles"
	"my-collection/server/pkg/server/tags"
	"my-collection/server/pkg/server/tasks"
	"my-collection/server/pkg/spectagger"
	"my-collection/server/pkg/storage"
)

func (mc *MyCollection) registerHandlers(db db.Database, storage *storage.Storage, fsm *fssync.FsManager) {
	mc.server.RegisterHandler(items.NewHandler(db, mc.processor, mc.itemsoptimizer))
	mc.server.RegisterHandler(tags.NewHandler(db, storage, mc.thumbnails))
	mc.server.RegisterHandler(storageHandler.NewHandler(storage))
	mc.server.RegisterHandler(fs.NewHandler(db, fsm))
	mc.server.RegisterHandler(tasks.NewHandler(db, mc.processor))
	mc.server.RegisterHandler(subtitles.NewHandler(db, &struct {
		opensubtitles.OpenSubtitiles
		model.TempFileProvider
	}{*mc.opensubtitles, storage}))
	mc.server.RegisterHandler(mc.push)
	mc.server.RegisterHandler(management.NewHandler(db, &struct {
		processor.Processor
		itemsoptimizer.ItemsOptimizer
		spectagger.Spectagger
		mixondemand.MixOnDemand
	}{*mc.processor, *mc.itemsoptimizer, *mc.spectagger, *mc.mixondemand}))
}
