package main

import (
	"fmt"
	"my-collection/server/pkg/db"
	"my-collection/server/pkg/server"

	"github.com/go-errors/errors"
)

func run() error {
	db, err := db.New("test.sqlite")

	if err != nil {
		return err
	}

	server.New(db).Run()
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
	logError(run())
}
