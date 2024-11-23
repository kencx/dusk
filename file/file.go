package file

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/kencx/dusk"
	"github.com/kencx/dusk/file/epub"
	"github.com/kencx/dusk/null"
)

const (
	coverFilename = "cover"
	epubExt       = ".epub"
	azwExt        = ".azw"
	mobiExt       = ".mobi"
	pdfExt        = ".pdf"
	jpegExt       = ".jpeg"
	pngExt        = ".png"
	djvuExt       = ".djvu"
)

type Service struct {
	Directory string
}

func NewService(path string) (*Service, error) {
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return nil, err
	}

	return &Service{path}, nil
}

// Upload format for new book
func (s *Service) UploadNewBook(payload *Payload) (*dusk.Book, error) {
	switch payload.Extension {
	case epubExt:
		return s.uploadNewEpub(payload)
	default:
		return nil, errors.New("unsupported file format")
	}
}

// Upload new format for existing book
func (s *Service) UploadBookFormat(payload *Payload, book *dusk.Book) error {
	switch payload.Extension {
	case epubExt:
		return s.uploadEpub(payload, book)
	default:
		return s.uploadFormatFile(payload, book)
	}
}

// Upload book cover from payload
func (s *Service) UploadCoverFromPayload(payload *Payload, book *dusk.Book) error {
	if err := s.uploadCover(payload.File, payload.Extension, book); err != nil {
		return err
	}
	return nil
}

// Upload book cover from URL for existing book
func (s *Service) UploadCoverFromUrl(url string, book *dusk.Book) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("file: failed to fetch file from url: %w", err)
	}
	defer resp.Body.Close()

	ext := path.Ext(path.Base(resp.Request.URL.Path))
	if ext == "" {
		ext = ".jpeg"
	}

	if err := s.uploadCover(resp.Body, ext, book); err != nil {
		return err
	}
	return nil
}

// Upload EPUB format for new book
func (s *Service) uploadNewEpub(payload *Payload) (*dusk.Book, error) {
	ep, err := epub.NewFromReader(payload.File, payload.Size)
	if err != nil && !errors.Is(err, epub.ErrNoCovers) {
		return nil, fmt.Errorf("failed to parse epub file: %w", err)
	}

	book := ep.ToBook()
	errMap := book.Valid()
	if len(errMap) > 0 {
		return nil, errMap
	}

	if err := s.uploadFormatFile(payload, book); err != nil {
		return nil, err
	}
	if err := s.uploadCoverFromEpub(ep, book); err != nil {
		return nil, err
	}
	return book, nil
}

// Upload EPUB format for existing book
func (s *Service) uploadEpub(payload *Payload, book *dusk.Book) error {
	ep, err := epub.NewFromReader(payload.File, payload.Size)
	if err != nil && !errors.Is(err, epub.ErrNoCovers) {
		return fmt.Errorf("failed to parse epub file: %w", err)
	}

	if err := s.uploadFormatFile(payload, book); err != nil {
		return err
	}
	if err := s.uploadCoverFromEpub(ep, book); err != nil {
		return err
	}
	return nil
}

// Upload format file for book
func (s *Service) uploadFormatFile(payload *Payload, book *dusk.Book) error {
	bookDir, err := s.createBookDirectory(book)
	if err != nil {
		return err
	}

	filename := fmt.Sprintf("%s%s", book.SafeTitle(), payload.Extension)
	fullPath := filepath.Join(bookDir, filename)
	if err := s.upload(payload.File, fullPath); err != nil {
		return err
	}

	book.Formats = append(book.Formats, getRelativePath(fullPath))
	return nil
}

// Find and upload cover image in epub
func (s *Service) uploadCoverFromEpub(ep *epub.Epub, book *dusk.Book) error {
	if ep.CoverFile == "" {
		slog.Debug("[epub] No valid cover files found", slog.String("title", ep.Title))
		return epub.ErrNoCovers
	}

	coverFile, err := ep.Open(ep.CoverFile)
	if err != nil {
		return err
	}
	defer coverFile.Close()

	if err := s.uploadCover(coverFile, filepath.Ext(ep.CoverFile), book); err != nil {
		return err
	}
	return nil
}

func (s *Service) uploadCover(f io.Reader, extension string, book *dusk.Book) error {
	bookDir, err := s.createBookDirectory(book)
	if err != nil {
		return err
	}

	filename := fmt.Sprintf("%s%s", coverFilename, extension)
	fullPath := filepath.Join(bookDir, filename)
	if err := s.upload(f, fullPath); err != nil {
		return err
	}

	book.Cover = null.StringFrom(getRelativePath(fullPath))
	return nil
}

// Upload file to path
func (s *Service) upload(file io.Reader, path string) error {
	dest, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		return fmt.Errorf("file: failed to create file: %w", err)
	}
	defer dest.Close()

	if _, err = io.Copy(dest, file); err != nil {
		return fmt.Errorf("file: failed to copy file to dest: %w", err)
	}

	slog.Info("[file] New file uploaded", slog.String("path", path))
	return nil
}

// create or get book directory
func (s *Service) createBookDirectory(book *dusk.Book) (string, error) {
	bookDir := filepath.Join(s.Directory, strings.ToLower(book.SafeTitle()))
	if err := os.MkdirAll(bookDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create book directory: %w", err)
	}
	return bookDir, nil
}

// get last 2 elements of file path - parentDir/filename.ext
func getRelativePath(path string) string {
	parentDir := filepath.Base(filepath.Dir(path))
	return filepath.Join(parentDir, filepath.Base(path))
}
