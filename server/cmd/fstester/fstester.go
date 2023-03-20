package main

import "fmt"

func main() {
	fmt.Print("Done\n")
}

// package main

// import (
// 	"fmt"
// 	"my-collection/server/pkg/bl/directories"
// 	"my-collection/server/pkg/db"
// 	"my-collection/server/pkg/fssync"
// 	"my-collection/server/pkg/model"

// 	"k8s.io/utils/pointer"
// )

// func fs_tester(db *db.Database) {
// 	fs, err := fssync.NewFsManager(db, true)
// 	if err != nil {
// 		logError(err)
// 		return
// 	}

// 	err = directories.CreateOrUpdateDirectory(db, &model.Directory{
// 		Path:     "/home/david/work/my-projects/my-collection/server",
// 		Excluded: pointer.Bool(false),
// 	})
// 	if err != nil {
// 		logError(err)
// 		return
// 	}

// 	err = directories.IncludeDirectory(db, "output/storage-template/tags-image/none")
// 	if err != nil {
// 		logError(err)
// 		return
// 	}

// 	err = fs.Sync()
// 	if err != nil {
// 		logError(err)
// 		return
// 	}

// 	fmt.Printf("Done")
// }
