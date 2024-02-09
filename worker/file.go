package worker

import (
	"dusk/epub"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/kennygrant/sanitize"
)

const (
	coverFilename = "cover.jpg"
)

type FileWorker struct {
	DataDir string
}

func NewFileWorker(dataDir string) (*FileWorker, error) {
	w := &FileWorker{dataDir}

	err := w.createDataDir()
	if err != nil {
		return nil, err
	}
	return w, nil
}

func (w *FileWorker) createDataDir() error {
	err := os.MkdirAll(w.DataDir, 0755)
	if err != nil {
		return err
	}
	return nil
}

func (w *FileWorker) UploadCoverFromFile(path, title string) (string, error) {
	f, err := epub.ExtractCoverFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to get cover file: %v", err)
	}
	defer f.Close()

	dest, err := w.UploadCover(f, title)
	if err != nil {
		return "", fmt.Errorf("failed to upload cover file: %v", err)
	}

	return dest, nil
}

func (w *FileWorker) UploadCover(cover io.Reader, title string) (string, error) {
	return w.UploadFile(cover, title, coverFilename)
}

func (w *FileWorker) UploadFile(file io.Reader, title, filename string) (string, error) {
	bookDir, err := w.CreateBookDir(title)
	if err != nil {
		return "", fmt.Errorf("failed to create book directory: %v", err)
	}

	path := filepath.Join(bookDir, filename)
	out, err := os.Create(path)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %v", err)
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		return "", fmt.Errorf("failed to copy file to dest: %v", err)
	}

	// TODO log upload file complete
	return path, nil
}

func (w *FileWorker) CreateBookDir(title string) (string, error) {
	sanitized := sanitize.Name(title)

	path := filepath.Join(w.DataDir, sanitized)
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return "", err
	}

	return path, nil
}

func (w *FileWorker) GetRelativePath(path string) string {
	parentDir := filepath.Base(filepath.Dir(path))
	return filepath.Join(parentDir, filepath.Base(path))
}
