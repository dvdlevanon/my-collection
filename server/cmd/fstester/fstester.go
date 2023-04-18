package main

import (
	"flag"
	"fmt"
	"my-collection/server/pkg/bl/directories"
	"my-collection/server/pkg/db"
	"my-collection/server/pkg/fssync"
	"my-collection/server/pkg/model"
	"my-collection/server/pkg/relativasor"
	"my-collection/server/pkg/utils"
	"net/http"
	"time"

	"github.com/op/go-logging"
	"k8s.io/utils/pointer"
)

var logger = logging.MustGetLogger("fstester")

var (
	help          = flag.Bool("help", false, "Print help")
	rootDirectory = flag.String("root-directory", "", "Server root directory")
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

	db, err := db.New(relativasor.GetRootDirectory(), "test.sqlite")
	if err != nil {
		return err
	}

	if err := directories.Init(db); err != nil {
		return err
	}
	fs, err := fssync.NewFsManager(db, filesFilter, 1*time.Second)
	if err != nil {
		return err
	}

	rootDir := &model.Directory{
		Path:     relativasor.GetRootDirectory(),
		Excluded: pointer.Bool(false),
	}
	if err := directories.CreateOrUpdateDirectory(db, rootDir); err != nil {
		return err
	}

	go startControlServer(db)
	fs.Watch()
	return nil
}

func startControlServer(db *db.Database) {
	http.HandleFunc("/include", func(w http.ResponseWriter, r *http.Request) {
		dir := r.URL.Query().Get("dir")
		utils.LogError(directories.IncludeDirectory(db, dir))
		fmt.Fprintf(w, "Done %s\n", dir)
	})

	http.HandleFunc("/exclude", func(w http.ResponseWriter, r *http.Request) {
		dir := r.URL.Query().Get("dir")
		utils.LogError(directories.ExcludeDirectory(db, dir))
		fmt.Fprintf(w, "Done %s\n", dir)
	})

	go http.ListenAndServe(":9999", nil)
}

func main() {
	utils.LogError(run())
}
