package main

import (
	"context"
	"flag"
	"my-collection/server/pkg/automix"
	"my-collection/server/pkg/bl/directories"
	"my-collection/server/pkg/bl/items"
	"my-collection/server/pkg/db"
	"my-collection/server/pkg/fssync"
	"my-collection/server/pkg/itemsoptimizer"
	"my-collection/server/pkg/mixondemand"
	processor "my-collection/server/pkg/processor"
	"my-collection/server/pkg/relativasor"
	"my-collection/server/pkg/server"
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
)

var logger = logging.MustGetLogger("mycollection")

var (
	help          = flag.Bool("help", false, "Print help")
	rootDirectory = flag.String("root-directory", "", "Server root directory")
	listenAddress = flag.String("address", ":6969", "Server listen address")
)

func filesFilter(path string) bool {
	return utils.IsVideo(true, path)
}

func run() error {
	flag.Parse()
	if *help {
		flag.Usage()
		return nil
	}

	if err := utils.ConfigureLogger(); err != nil {
		return err
	}

	if err := relativasor.Init(*rootDirectory); err != nil {
		return err
	}

	dataDir := path.Join(relativasor.GetRootDirectory(), ".mycollection")
	if err := os.MkdirAll(dataDir, 0750); err != nil {
		return err
	}

	logger.Infof("Root directory is: %s", relativasor.GetRootDirectory())

	db, err := db.New(filepath.Join(dataDir, "db.sqlite"))
	if err != nil {
		return err
	}

	storage, err := storage.New(filepath.Join(dataDir, "storage"))
	if err != nil {
		return err
	}

	ctx, stopSignal := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stopSignal()

	processor, err := processor.New(db, storage)
	if err != nil {
		return err
	}
	processor.Continue()
	go processor.Run(ctx)

	if err := items.InitHighlights(db); err != nil {
		return err
	}
	if err := directories.Init(db); err != nil {
		return err
	}
	fsManager, err := fssync.NewFsManager(db, filesFilter, 60*time.Second)
	if err != nil {
		return err
	}
	go fsManager.Watch(ctx)

	automix, err := automix.New(db, db, db, 40)
	if err != nil {
		return err
	}
	go automix.Run(ctx)

	mixOnDemand, err := mixondemand.New(db, db, db, 20)
	if err != nil {
		return err
	}

	spectagger, err := spectagger.New(db, db, db)
	if err != nil {
		return err
	}
	go spectagger.Run(ctx)

	itemsoptimizer := itemsoptimizer.New(db, processor, 1080)
	go itemsoptimizer.Run(ctx)

	thumbnails := thumbnails.New(db, db, storage, 100, 100)
	go thumbnails.Run(ctx)

	return server.New(db, storage, fsManager, processor, spectagger, itemsoptimizer, thumbnails, mixOnDemand).Run(ctx, *listenAddress)
}

func main() {
	if err := run(); err != nil {
		utils.LogError("Error in main", err)
		os.Exit(1)
	}
}
