package main

import (
	"flag"
	"fmt"
	"my-collection/server/pkg/bl/items"
	"my-collection/server/pkg/db"
	"my-collection/server/pkg/model"
	"my-collection/server/pkg/utils"
	"os"

	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("fstester")

var (
	help   = flag.Bool("help", false, "Print help")
	dbFile = flag.String("db-file", "", "A path to the sqlite database file")
)

func run() error {
	flag.Parse()
	if *help {
		flag.Usage()
		return nil
	}

	if err := utils.ConfigureLogger(); err != nil {
		return err
	}

	db, err := db.New(*dbFile)
	if err != nil {
		return err
	}

	allItems, err := db.GetAllItems()
	if err != nil {
		return err
	}

	itemsMap := make(map[uint64]model.Item, len(*allItems))
	for _, item := range *allItems {
		itemsMap[item.Id] = item
	}

	for _, item := range itemsMap {
		if items.IsHighlight(&item) {
			parent, ok := itemsMap[*item.HighlightParentItemId]
			if !ok {
				logger.Infof("Missing highlight parent %v", item)
			}

			if parent.Url != item.Url {
				fmt.Printf("Highlight URL not match %s != %s\n", parent.Url, item.Url)
			}
		}

		if items.IsSubItem(&item) {
			parent, ok := itemsMap[*item.MainItemId]
			if !ok {
				logger.Infof("Missing main item %v", item)
			}

			if parent.Url != item.Url {
				fmt.Printf("Subitem URL not match %s != %s\n", parent.Url, item.Url)
			}
		}
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		utils.LogError("Error in main", err)
		os.Exit(1)
	}
}
