package storage

import (
	"fmt"
	"io"
	"log/slog"
	"os"
)

const basePath = "../../.data"

// Storage defines the logic of storage's operations with data.
type Storage struct {
	log *slog.Logger
}

func New(log *slog.Logger) Storage {
	return Storage{
		log: log,
	}
}

// getPath returns the full path for the current resource.
func (s Storage) getPath(path string) string {
	return fmt.Sprintf("%s/%s", basePath, path)
}

func (s Storage) readFile(source io.Reader, op string) ([]byte, error) {
	data := make([]byte, 0, 10000)

	for {
		buffer := make([]byte, 10000)
		n, err := source.Read(buffer)

		if n != 0 && (err == nil || err == io.EOF) {
			data = append(data, buffer[:n]...)
		} else if err != nil && err != io.EOF {
			s.log.Warn(fmt.Sprintf("error of the %v: %v", op, err))
			return nil, fmt.Errorf("error of the %s: %s", op, err)
		} else if err == io.EOF {
			break
		}
	}

	return data, nil
}

// CreateNameSpace creates a new namespace.
func (s Storage) CreateNameSpace(name string) error {
	const op = "storage.create-name-space"

	if len(name) == 0 {
		s.log.Warn(fmt.Sprintf("error of the %s: %s", op, ErrEmptyItem))
		return fmt.Errorf("error of the %s: %w", op, ErrEmptyItem)
	}

	if err := os.Mkdir(s.getPath(name), 0740); err != nil {
		s.log.Warn(fmt.Sprintf("error of the %s: %s", op, err))
		return fmt.Errorf("error of the %s: %s", op, err)
	}

	return nil
}

// DeleteNameSpace deletes the namespace.
func (s Storage) DeleteNameSpace(name string) error {
	const op = "storage.delete-name-space"

	if len(name) == 0 {
		s.log.Warn(fmt.Sprintf("error of the %s: %s", op, ErrEmptyItem))
		return fmt.Errorf("error of the %s: %w", op, ErrEmptyItem)
	}

	if err := os.RemoveAll(s.getPath(name)); err != nil {
		s.log.Warn(fmt.Sprintf("error of the %s: %s", op, err))
		return fmt.Errorf("error of the %s: %v", op, err)
	}

	return nil
}

// CreateDir creates a new directory with the set path (path that set from namespace).
func (s Storage) CreateDir(path string) error {
	const op = "storage.create-dir"

	if len(path) == 0 {
		s.log.Warn(fmt.Sprintf("error of the %s: %s", op, ErrEmptyItem))
		return fmt.Errorf("error of the %s: %w", op, ErrEmptyItem)
	}

	if err := os.Mkdir(s.getPath(path), 0740); err != nil {
		s.log.Warn(fmt.Sprintf("error of the %s: %s", op, err))
		return fmt.Errorf("error of the %s: %s", op, err)
	}

	return nil
}

// DeleteDir deletes the dir from the namespace.
func (s Storage) DeleteDir(path string) error {
	const op = "storage.delete-name-space"

	if len(path) == 0 {
		s.log.Warn(fmt.Sprintf("error of the %s: %s", op, ErrEmptyItem))
		return fmt.Errorf("error of the %s: %w", op, ErrEmptyItem)
	}

	if err := os.RemoveAll(s.getPath(path)); err != nil {
		s.log.Warn(fmt.Sprintf("error of the %s: %s", op, err))
		return fmt.Errorf("error of the %s: %v", op, err)
	}

	return nil
}

// CreateFile creates the file.
func (s Storage) CreateFile(path string, data []byte) error {
	const op = "storage.create-file"

	file, err := os.OpenFile(s.getPath(path), os.O_RDWR|os.O_CREATE, 0740)

	if err != nil {
		s.log.Error(fmt.Sprintf("error of the %s: %s", op, err))
		return fmt.Errorf("error of the %s: %s", op, err)
	}
	defer file.Close()

	if _, err := file.Write(data); err != nil {
		s.log.Error(fmt.Sprintf("error of the %s: %s", op, err))
		return fmt.Errorf("error of the %s: %s", op, err)
	}

	return nil
}

// DeleteFile deletes the file.
func (s Storage) DeleteFile(path string) error {
	const op = "storage.delete-file"

	if err := os.Remove(s.getPath(path)); err != nil {
		s.log.Error(fmt.Sprintf("error of the %s: %s", op, err))
		return fmt.Errorf("error of the %s: %s", op, err)
	}

	return nil
}

// GetFile returns the file's data.
func (s Storage) GetFile(path string) ([]byte, error) {
	const op = "storage.get-file"

	file, err := os.OpenFile(s.getPath(path), os.O_RDWR, 0)

	if err != nil {
		s.log.Error(fmt.Sprintf("error of the %s: %s", op, err))
		return nil, fmt.Errorf("error of the %s: %s", op, err)
	}
	defer file.Close()

	data, err := s.readFile(file, op)

	if err != nil {
		return nil, err
	}

	return data, nil
}

// CopyFile copies file.
func (s Storage) CopyFile(from string, to string) error {
	const op = "storage.copy-file"

	source, err := os.OpenFile(s.getPath(from), os.O_RDONLY, 0)

	if err != nil {
		s.log.Error(fmt.Sprintf("error of the %s: %s", op, err))
		return fmt.Errorf("error of the %s: %s", op, err)
	}
	defer source.Close()

	file, err := os.Create(s.getPath(to))

	if err != nil {
		s.log.Error(fmt.Sprintf("error of the %s: %s", op, err))
		return fmt.Errorf("error of the %s: %s", op, err)
	}
	defer file.Close()

	data, err := s.readFile(source, op)

	if err != nil {
		return err
	}

	if _, err := file.Write(data); err != nil {
		s.log.Error(fmt.Sprintf("error of the %s: %s", op, err))
		return fmt.Errorf("error of the %s: %s", op, err)
	}

	return nil
}
