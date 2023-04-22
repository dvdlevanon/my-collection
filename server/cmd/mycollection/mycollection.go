package main

import (
	"flag"
	"my-collection/server/pkg/automix"
	"my-collection/server/pkg/bl/directories"
	"my-collection/server/pkg/bl/items"
	"my-collection/server/pkg/db"
	"my-collection/server/pkg/fssync"
	processor "my-collection/server/pkg/processor"
	"my-collection/server/pkg/relativasor"
	"my-collection/server/pkg/server"
	"my-collection/server/pkg/storage"
	"my-collection/server/pkg/utils"
	"path/filepath"
	"time"

	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("mycollection")

var (
	help          = flag.Bool("help", false, "Print help")
	rootDirectory = flag.String("root-directory", "", "Server root directory")
	listenAddress = flag.String("address", ":8080", "Server listen address")
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

	logger.Infof("Root directory is: %s", relativasor.GetRootDirectory())

	db, err := db.New(relativasor.GetRootDirectory(), "db.sqlite")
	if err != nil {
		return err
	}

	storage, err := storage.New(filepath.Join(relativasor.GetRootDirectory(), ".storage"))
	if err != nil {
		return err
	}

	processor, err := processor.New(db, storage)
	if err != nil {
		return err
	}
	processor.Pause()
	go processor.Run()

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
	go fsManager.Watch()

	automix, err := automix.New(db, db, db, 40)
	if err != nil {
		return err
	}
	go automix.Run()

	return server.New(db, storage, fsManager, processor).Run(*listenAddress)
}

func main() {
	utils.LogError(run())
}
