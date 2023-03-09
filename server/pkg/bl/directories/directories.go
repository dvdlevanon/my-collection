package directories

import (
	"my-collection/server/pkg/model"
	"my-collection/server/pkg/relativasor"
	"path/filepath"
	"strings"

	"github.com/op/go-logging"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"k8s.io/utils/pointer"
)

var logger = logging.MustGetLogger("directories")

func ExcludeDirectory(drw model.DirectoryReaderWriter, path string) error {
	path = relativasor.GetRelativePath(path)

	directory, err := drw.GetDirectory(path)
	if err != nil {
		return err
	}

	if *directory.Excluded {
		return nil
	}

	directory.Excluded = pointer.Bool(true)
	return drw.CreateOrUpdateDirectory(directory)
}

func DirectoryNameToTag(path string) string {
	caser := cases.Title(language.English)
	return caser.String(strings.ReplaceAll(strings.ReplaceAll(filepath.Base(path), "-", " "), "_", " "))
}

func TagTitleToDirectory(title string) string {
	return strings.ToLower(strings.ReplaceAll(title, " ", "-"))
}

func tagExists(tag *model.Tag, tags []*model.Tag) bool {
	for _, t := range tags {
		if tag.Id == t.Id {
			return true
		}
	}

	return false
}

func SetDirectoryTags(drw model.DirectoryReaderWriter, directory *model.Directory) error {
	existingDirectory, err := drw.GetDirectory("path = ?", directory.Path)
	if err != nil {
		logger.Errorf("Error getting exising directory %s %t", directory.Path, err)
		return err
	}

	for _, tag := range existingDirectory.Tags {
		if tagExists(tag, directory.Tags) {
			continue
		}

		if err := drw.RemoveTagFromDirectory(directory.Path, tag.Id); err != nil {
			logger.Warningf("Unable to remove tag %d from directory %s - %t",
				directory.Path, tag.Id, err)
		}
	}

	return drw.CreateOrUpdateDirectory(directory)
}
