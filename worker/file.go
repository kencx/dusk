package worker

import (
	"os"
	"path/filepath"
)

type FileWorker struct {
	dataDir string
}

func NewFileWorker(dataDir string) (*FileWorker, error) {
	w := &FileWorker{dataDir}

	err := w.initDir()
	if err != nil {
		return nil, err
	}
	return w, nil
}

func (w *FileWorker) initDir() error {
	err := os.MkdirAll(w.dataDir, 0755)
	if err != nil {
		return err
	}
	return nil
}

func (w *FileWorker) GeneratePath(title string) (string, error) {
	sanitized := sanitize(title)

	fullPath := filepath.Join(w.dataDir, sanitized)
	err := os.Mkdir(fullPath, 0755)
	if err != nil {
		return "", err
	}
	return fullPath, nil
}

func sanitize(f string) string {
	return f
}
