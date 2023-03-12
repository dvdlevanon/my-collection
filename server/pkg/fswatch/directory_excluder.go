package fswatch

import (
	"my-collection/server/pkg/bl/directories"
	"my-collection/server/pkg/bl/tags"
	"my-collection/server/pkg/model"
)

func newDirectoryExcluder(directory *model.Directory, tr model.TagReader) (*directoryExcluder, error) {
	title := directories.DirectoryNameToTag(directory.Path)
	tag, err := tags.GetChildTag(tr, directories.DIRECTORIES_TAG_ID, title)
	if err != nil {
		return nil, err
	}

	return &directoryExcluder{
		directory: directory,
		tag:       tag,
	}, nil
}

type directoryExcluder struct {
	directory *model.Directory
	tag       *model.Tag
}

func (d *directoryExcluder) process() error {

	return nil
}
