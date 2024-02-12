package file

import (
	"dusk"
	"dusk/file/epub"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"os"
	"path/filepath"
)

const (
	coverFilename = "cover"
)

type Payload struct {
	multipart.File
	Size     int64
	Filename string
	MimeType string
}

type Worker struct {
	DataDir string
}

func NewWorker(dataDir string) (*Worker, error) {
	err := os.MkdirAll(dataDir, 0755)
	if err != nil {
		return nil, err
	}

	return &Worker{dataDir}, nil
}

// Upload new format for new book
func (w *Worker) UploadBook(payload *Payload) (*dusk.Book, error) {
	switch payload.Extension() {
	case ".epub":
		return w.UploadNewEpub(payload)
	default:
		return nil, errors.New("unsupported file format")
	}
}

// Upload new format for existing book
func (w *Worker) UploadBookFormat(payload *Payload, book *dusk.Book) error {
	switch payload.Extension() {
	case ".epub":
		return w.UploadEpub(payload, book)
	default:
		return w.UploadOtherFormat(payload, book)
	}
}

// Upload EPUB format for new book
func (w *Worker) UploadNewEpub(payload *Payload) (*dusk.Book, error) {
	ep, err := epub.NewFromReader(payload.File, payload.Size)
	if err != nil && !errors.Is(err, epub.ErrNoCovers) {
		return nil, fmt.Errorf("failed to parse epub file: %w", err)
	}

	book := ep.ToBook()
	errMap := book.Valid()
	if len(errMap) > 0 {
		return nil, errMap
	}

	if err := w.UploadOtherFormat(payload, book); err != nil {
		return nil, err
	}

	if ep.CoverFile != "" {
		coverFile, err := ep.Open(ep.CoverFile)
		if err != nil {
			return nil, err
		}
		defer coverFile.Close()

		bookDir := filepath.Join(w.DataDir, book.SafeTitle())
		filename := fmt.Sprintf("%s%s", coverFilename, filepath.Ext(ep.CoverFile))
		fullPath := filepath.Join(bookDir, filename)
		if err = w.UploadFile(coverFile, fullPath); err != nil {
			return nil, err
		}
		book.Cover.String = getRelativePath(fullPath)
	}

	return book, nil
}

// Upload EPUB format for existing book
func (w *Worker) UploadEpub(payload *Payload, book *dusk.Book) error {
	ep, err := epub.NewFromReader(payload.File, payload.Size)
	if err != nil && !errors.Is(err, epub.ErrNoCovers) {
		return fmt.Errorf("failed to parse epub file: %w", err)
	}

	if err := w.UploadOtherFormat(payload, book); err != nil {
		return err
	}

	if ep.CoverFile != "" {
		coverFile, err := ep.Open(ep.CoverFile)
		if err != nil {
			return err
		}
		defer coverFile.Close()

		bookDir := filepath.Join(w.DataDir, book.SafeTitle())
		filename := fmt.Sprintf("%s%s", coverFilename, filepath.Ext(ep.CoverFile))
		fullPath := filepath.Join(bookDir, filename)
		if err = w.UploadFile(coverFile, fullPath); err != nil {
			return err
		}
		book.Cover.String = getRelativePath(fullPath)
	}

	return nil
}

// Upload other format for existing book
func (w *Worker) UploadOtherFormat(payload *Payload, book *dusk.Book) error {
	// check or create book folder
	bookDir := filepath.Join(w.DataDir, book.SafeTitle())
	if err := os.MkdirAll(bookDir, 0755); err != nil {
		return fmt.Errorf("failed to create book directory: %w", err)
	}

	filename := fmt.Sprintf("%s%s", book.SafeTitle(), payload.Extension())
	fullPath := filepath.Join(bookDir, filename)
	if err := w.UploadFile(payload.File, fullPath); err != nil {
		return err
	}

	book.Formats = append(book.Formats, getRelativePath(fullPath))
	return nil
}

// Upload book cover for existing book
func (w *Worker) UploadBookCover(payload *Payload, book *dusk.Book) error {
	// check or create book folder
	bookDir := filepath.Join(w.DataDir, book.SafeTitle())
	if err := os.MkdirAll(bookDir, 0755); err != nil {
		return fmt.Errorf("failed to create book directory: %w", err)
	}

	filename := fmt.Sprintf("%s%s", coverFilename, payload.Extension())
	fullPath := filepath.Join(bookDir, filename)
	if err := w.UploadFile(payload.File, fullPath); err != nil {
		return err
	}

	book.Cover.String = getRelativePath(fullPath)
	return nil
}

// Upload file to path
func (w *Worker) UploadFile(file io.Reader, path string) error {
	dest, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("file: failed to create file: %w", err)
	}
	defer dest.Close()

	if _, err = io.Copy(dest, file); err != nil {
		return fmt.Errorf("file: failed to create file: %w", err)
	}

	slog.Info("[file] New file uploaded", slog.String("path", path))
	return nil
}

func getRelativePath(path string) string {
	parentDir := filepath.Base(filepath.Dir(path))
	return filepath.Join(parentDir, filepath.Base(path))
}
