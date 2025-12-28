package storage

import (
	"context"
	"my-collection/server/pkg/utils"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/go-errors/errors"
	"github.com/google/uuid"
	cp "github.com/otiai10/copy"
)

const (
	TEMP_DIRECTORY     = "temp"
	STORAGE_URL_PREFIX = ".internal-storage"
)

var logger = utils.MustGetLogger("storage")

type Storage struct {
	rootDirectory string
}

func New(ctx context.Context, rootDirectory string) (*Storage, error) {
	if err := os.MkdirAll(rootDirectory, 0750); err != nil {
		return nil, errors.Wrap(err, 0)
	}

	if err := os.MkdirAll(path.Join(rootDirectory, TEMP_DIRECTORY), 0750); err != nil {
		return nil, errors.Wrap(err, 0)
	}

	cwd, err := os.Getwd()
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	storageTemplateDirectory := filepath.Join(cwd, "storage-template")
	if err := cp.Copy(storageTemplateDirectory, rootDirectory, cp.Options{}); err != nil {
		logger.Warningf(ctx, "Error copying storage template from %s to %s - %s", storageTemplateDirectory, rootDirectory, err)
	}

	logger.Infof(ctx, "Srorage initialized in %s", rootDirectory)

	return &Storage{
		rootDirectory: rootDirectory,
	}, nil
}

func (s *Storage) GetStorageDirectory(name string) string {
	return s.GetFile(name)
}

func (s *Storage) Put(name string, bytes []byte) error {
	return os.WriteFile(s.GetFile(name), bytes, 0750)
}

func (s *Storage) Get(name string) ([]byte, error) {
	return os.ReadFile(s.GetFile(name))
}

func (s *Storage) GetFile(name string) string {
	return filepath.Join(s.rootDirectory, strings.TrimPrefix(name, STORAGE_URL_PREFIX))
}

func (s *Storage) GetStorageUrl(name string) string {
	return filepath.Join(".internal-storage", name)
}

func (s *Storage) IsStorageUrl(name string) bool {
	return strings.HasPrefix(name, STORAGE_URL_PREFIX)
}

func (s *Storage) GetFileForWriting(name string) (string, error) {
	path := s.GetFile(name)
	if err := os.MkdirAll(filepath.Dir(path), 0750); err != nil {
		return "", errors.Wrap(err, 0)
	}

	return path, nil
}

func (s *Storage) GetTempFile() string {
	return filepath.Join(s.rootDirectory, TEMP_DIRECTORY, uuid.New().String())
}
