package cache

import (
	"io"
	"os"
	"path/filepath"
	"sync"
)

type FS struct {
	root  string
	locks sync.Map
}

func NewFS(root string) *FS {
	return &FS{root: root}
}

func (s *FS) Get(path string) (io.ReadSeekCloser, os.FileInfo, error) {
	fullPath := filepath.Join(s.root, path)

	if wg, ok := s.locks.Load(fullPath); ok {
		wg.(*sync.WaitGroup).Wait()
	}

	file, err := os.Open(fullPath)
	if err != nil {
		return nil, nil, err
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, nil, err
	}

	return file, info, nil
}

func (s *FS) Put(path string, data []byte) error {
	fullPath := filepath.Join(s.root, path)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	actual, loaded := s.locks.LoadOrStore(fullPath, wg)

	if loaded {
		actual.(*sync.WaitGroup).Wait()
		return nil
	}

	defer func() {
		wg.Done()
		s.locks.Delete(fullPath)
	}()

	if s.Exists(path) {
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return err
	}

	tmpFile, err := os.CreateTemp(filepath.Dir(fullPath), "tmp-")
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write(data); err != nil {
		tmpFile.Close()
		return err
	}
	tmpFile.Close()

	return os.Rename(tmpFile.Name(), fullPath)
}

func (s *FS) Exists(path string) bool {
	_, err := os.Stat(filepath.Join(s.root, path))
	return err == nil
}
