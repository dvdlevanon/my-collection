package main

import (
	"flag"
	"my-collection/server/pkg/db"
	"my-collection/server/pkg/directories"
	"my-collection/server/pkg/gallery"
	"my-collection/server/pkg/server"
	"my-collection/server/pkg/storage"
	"os"
	"path/filepath"

	"github.com/go-errors/errors"
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

	logger.Infof("Root directory is: %s", rootdir)

	db, err := db.New(rootdir, "test.sqlite")
	if err != nil {
		return err
	}

	storage, err := storage.New(filepath.Join(rootdir, ".storage"))
	if err != nil {
		return err
	}

	gallery := gallery.New(db, storage, rootdir)
	directories := directories.New(gallery, storage)
	if err = directories.Init(); err != nil {
		return err
	}

	return server.New(gallery, storage, directories).Run(*listenAddress)
}

func logError(err error) {
	if err == nil {
		return
	}

	var e *errors.Error
	if errors.As(err, &e) {
		logger.Errorf("Error: %v", e.ErrorStack())
	} else {
		logger.Errorf("Error: %v", err)
	}
}

func main() {
	logError(run())
}
