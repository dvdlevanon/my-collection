package app

import (
	"context"
	"fmt"
	"my-collection/server/pkg/automix"
	"my-collection/server/pkg/bl/directories"
	"my-collection/server/pkg/bl/items"
	"my-collection/server/pkg/db"
	"my-collection/server/pkg/fssync"
	"my-collection/server/pkg/itemsoptimizer"
	"my-collection/server/pkg/mixondemand"
	"my-collection/server/pkg/opensubtitles"
	"my-collection/server/pkg/processor"
	"my-collection/server/pkg/relativasor"
	"my-collection/server/pkg/server"
	"my-collection/server/pkg/server/push"
	"my-collection/server/pkg/spectagger"
	"my-collection/server/pkg/storage"
	"my-collection/server/pkg/thumbnails"
	"my-collection/server/pkg/utils"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"syscall"
	"time"

	"github.com/op/go-logging"
	"golang.org/x/sync/errgroup"
)

var logger = logging.MustGetLogger("mycollection")

func New(config MyCollectionConfig) (*MyCollection, error) {
	mc := &MyCollection{}
	if err := mc.initialize(config); err != nil {
		return nil, err
	}
	return mc, nil
}

type MyCollection struct {
	processor      *processor.Processor
	fsManager      *fssync.FsManager
	automix        *automix.Automix
	mixondemand    *mixondemand.MixOnDemand
	spectagger     *spectagger.Spectagger
	itemsoptimizer *itemsoptimizer.ItemsOptimizer
	thumbnails     *thumbnails.Thumbnails
	server         *server.Server
	push           push.PushHandler
	opensubtitles  *opensubtitles.OpenSubtitiles
}

func (mc *MyCollection) initialize(config MyCollectionConfig) error {
	ctx := utils.ContextWithSubject(context.TODO(), "init")

	if err := relativasor.Init(config.RootDir); err != nil {
		return err
	}

	dataDir := path.Join(relativasor.GetRootDirectory(), ".mycollection")
	if err := os.MkdirAll(dataDir, 0750); err != nil {
		return err
	}

	logger.Infof("Root directory is: %s", relativasor.GetRootDirectory())

	db, err := db.New(filepath.Join(dataDir, "db.sqlite"), false)
	if err != nil {
		return err
	}

	storage, err := storage.New(filepath.Join(dataDir, "storage"))
	if err != nil {
		return err
	}

	mc.processor, err = processor.New(db, storage, config.ProcessorPaused, config.CoversCount, config.PreviewSceneCount, config.PreviewSceneDuration)
	if err != nil {
		return err
	}

	if err := items.InitHighlights(ctx, db); err != nil {
		return err
	}
	if err := directories.Init(ctx, db); err != nil {
		return err
	}
	mc.fsManager, err = fssync.NewFsManager(ctx, db, config.FilesFilter, 60*time.Second)
	if err != nil {
		return err
	}

	mc.automix, err = automix.New(ctx, db, config.AutoMixItemsCount)
	if err != nil {
		return err
	}

	mc.mixondemand, err = mixondemand.New(ctx, db, config.MixOnDemandItemsCount)
	if err != nil {
		return err
	}

	mc.spectagger, err = spectagger.New(ctx, db)
	if err != nil {
		return err
	}

	mc.opensubtitles = opensubtitles.NewOpenSubtitles(config.OpenSubtitleApiKeys)

	mc.itemsoptimizer = itemsoptimizer.New(db, mc.processor, config.ItemsOptimizerMaxResolution)
	mc.thumbnails = thumbnails.New(db, db, storage, 100, 100)
	mc.server = server.New(config.ListenAddress)
	mc.push = push.NewPush()

	mc.fsManager.AddPushListener(mc.push)
	mc.processor.AddPushListener(mc.push)

	mc.registerHandlers(db, storage, mc.fsManager)

	return nil
}

func (mc *MyCollection) Run(ctx context.Context) error {
	eg, ctx := errgroup.WithContext(ctx)
	if err := mc.runWithErrgroup(ctx, eg); err != nil {
		return err
	}

	return mc.waitErrgroup(eg)
}

func (mc *MyCollection) runWithErrgroup(ctx context.Context, eg *errgroup.Group) error {
	ctx, stopSignal := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stopSignal()

	eg.Go(func() error {
		return mc.processor.Run(ctx)
	})

	eg.Go(func() error {
		return mc.fsManager.Watch(ctx)
	})

	eg.Go(func() error {
		return mc.automix.Run(ctx)
	})

	eg.Go(func() error {
		return mc.spectagger.Run(ctx)
	})

	eg.Go(func() error {
		return mc.itemsoptimizer.Run(ctx)
	})

	eg.Go(func() error {
		return mc.thumbnails.Run(ctx)
	})

	eg.Go(func() error {
		return mc.push.Run(ctx)
	})

	eg.Go(func() error {
		return mc.server.Run(ctx)
	})

	<-ctx.Done()
	return nil
}

func (mc *MyCollection) waitErrgroup(eg *errgroup.Group) error {
	done := make(chan error, 1)
	go func() {
		done <- eg.Wait()
	}()

	select {
	case err := <-done:
		if err == context.Canceled {
			return nil
		}
		return err
	case <-time.After(5 * time.Second):
		return fmt.Errorf("operation timed out")
	}
}
