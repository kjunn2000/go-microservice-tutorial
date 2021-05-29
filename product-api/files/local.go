package files

import (
	"io"
	"os"
	"path/filepath"

	"golang.org/x/xerrors"
)

type Local struct {
	maxFileSize int
	basePath    string
}

func NewLocal(basePath string, maxFileSize int) (*Local, error) {

	path, err := filepath.Abs(basePath)
	if err != nil {
		return nil, err
	}
	return &Local{basePath: path}, nil
}

type Storage interface {
	Save(path string, file io.Reader) error
}

func (l *Local) Save(path string, contents io.Reader) error {

	fp := l.fullPath(path)

	d := filepath.Dir(fp)
	err := os.MkdirAll(d, os.ModePerm)
	if err != nil {
		return xerrors.Errorf("Unable to create directory: %w", err)
	}

	_, err = os.Stat(fp)
	if err == nil {
		err = os.Remove(fp)
		if err != nil {
			return xerrors.Errorf("Unable to delete file : %w", err)
		}
	} else if !os.IsNotExist(err) {
		return xerrors.Errorf("Unabel to get file info: %w", err)
	}

	f, err := os.Create(fp)
	if err != nil {
		return xerrors.Errorf("Unable to create file : %w", err)
	}
	defer f.Close()

	_, err = io.Copy(f, contents)
	if err != nil {
		return xerrors.Errorf("Unable to write to file: %w", err)
	}

	return nil
}

func (l *Local) Get(path string) (*os.File, error) {
	fp := l.fullPath(path)

	f, err := os.Open(fp)
	if err != nil {
		return nil, xerrors.Errorf("Unable to open file: %w", err)
	}
	return f, nil
}

func (l *Local) fullPath(path string) string {
	return filepath.Join(l.basePath, path)
}
