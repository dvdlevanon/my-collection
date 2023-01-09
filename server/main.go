package main

import (
	"flag"
	"fmt"
	"my-collection/server/pkg/db"
	"my-collection/server/pkg/gallery"
	"my-collection/server/pkg/server"

	"github.com/go-errors/errors"
)

var help = flag.Bool("help", false, "Print help")
var rootDirectory = flag.String("root-directory", "", "Server root directory")

func run() error {
	db, err := db.New("test.sqlite")

	if err != nil {
		return err
	}

	gallery := gallery.New(db, *rootDirectory)
	server.New(gallery).Run()
	return nil
}

func logError(err error) {
	if err == nil {
		return
	}

	var e *errors.Error
	if errors.As(err, &e) {
		fmt.Printf("Error: %v", e.ErrorStack())
	} else {
		fmt.Printf("Error: %v", err)
	}
}

func main() {
	flag.Parse()

	if *help {
		flag.Usage()
		return
	}

	logError(run())
}
