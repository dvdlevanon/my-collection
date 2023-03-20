package main

import (
	"flag"
	"my-collection/server/pkg/db"
	"my-collection/server/pkg/fssync"
	processor "my-collection/server/pkg/processor"
	"my-collection/server/pkg/relativasor"
	"my-collection/server/pkg/server"
	"my-collection/server/pkg/storage"
	"my-collection/server/pkg/utils"
	"os"
	"path/filepath"

	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("main")

var (
	help          = flag.Bool("help", false, "Print help")
	rootDirectory = flag.String("root-directory", "", "Server root directory")
	listenAddress = flag.String("address", ":8080", "Server listen address")
)

func configureLogger() error {
	logFormat := `[%{time:2006-01-02 15:04:05.000}] %{color}%{level:-7s}%{color:reset} %{message} [%{module} - %{shortfile}]`
	formatter, err := logging.NewStringFormatter(logFormat)
	if err != nil {
		return err
	}

	logging.SetBackend(logging.NewLogBackend(os.Stdout, "", 0))
	logging.SetFormatter(formatter)

	logger.Debugf("Logger initialized with format %v", logFormat)
	return nil
}

func getRootDirectory() (string, error) {
	if *rootDirectory != "" {
		return *rootDirectory, nil
	}

	path, err := os.Getwd()

	if err != nil {
		return "", err
	}

	return path, nil
}

func run() error {
	flag.Parse()
	if *help {
		flag.Usage()
		return nil
	}

	if err := configureLogger(); err != nil {
		return err
	}

	rootdir, err := getRootDirectory()
	if err != nil {
		return err
	}

	relativasor.Init(rootdir)

	logger.Infof("Root directory is: %s", rootdir)

	db, err := db.New(rootdir, "db.sqlite")
	if err != nil {
		return err
	}

	storage, err := storage.New(filepath.Join(rootdir, ".storage"))
	if err != nil {
		return err
	}

	processor, err := processor.New(db, storage)
	if err != nil {
		return err
	}
	processor.Pause()
	go processor.Run()

	fsManager, err := fssync.NewFsManager(db, true)
	if err != nil {
		return err
	}
	go fsManager.Watch()

	return server.New(db, storage, fsManager, processor).Run(*listenAddress)
}

func main() {
	utils.LogError(run())
}
