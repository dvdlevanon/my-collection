package storage

import (
	"os"
	"path/filepath"

	"github.com/go-errors/errors"
	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("storage")

type Storage struct {
	rootDirectory string
}

func New(rootDirectory string) (*Storage, error) {
	if err := os.MkdirAll(rootDirectory, 0750); err != nil {
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
	return filepath.Join(s.rootDirectory, name)
}

func (s *Storage) GetFileForWriting(name string) (string, error) {
	path := s.GetFile(name)
	if err := os.MkdirAll(filepath.Dir(path), 0750); err != nil {
		return "", errors.Wrap(err, 0)
	}

	return path, nil
}
