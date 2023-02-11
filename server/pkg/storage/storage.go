package storage

import (
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/go-errors/errors"
	"github.com/google/uuid"
	"github.com/op/go-logging"
	cp "github.com/otiai10/copy"
)

const (
	TEMP_DIRECTORY     = "temp"
	STORAGE_URL_PREFIX = ".internal-storage"
)

var logger = logging.MustGetLogger("storage")

type Storage struct {
	rootDirectory string
}

func New(rootDirectory string) (*Storage, error) {
	if err := os.MkdirAll(rootDirectory, 0750); err != nil {
		return nil, errors.Wrap(err, 0)
	}

	if err := os.MkdirAll(path.Join(rootDirectory, TEMP_DIRECTORY), 0750); err != nil {
		return nil, errors.Wrap(err, 0)
	}

	execLocation, err := os.Executable()
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	storageTemplateDirectory := filepath.Join(filepath.Dir(execLocation), "storage-template")
	if err := cp.Copy(storageTemplateDirectory, rootDirectory, cp.Options{}); err != nil {
		return nil, errors.Wrap(err, 0)
	}

	logger.Infof("Srorage initialized in %s", rootDirectory)

	return &Storage{
		rootDirectory: rootDirectory,
	}, nil
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
