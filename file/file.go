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

	"github.com/kencx/dusk"
	"github.com/kencx/dusk/file/epub"
	"github.com/kencx/dusk/null"
)

const (
	coverFilename = "cover"
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

// Upload new format for new book
func (s *Service) UploadBook(payload *Payload) (*dusk.Book, error) {
	switch payload.Extension {
	case ".epub":
		return s.UploadNewEpub(payload)
	default:
		return nil, errors.New("unsupported file format")
	}
}

// Upload new format for existing book
func (s *Service) UploadBookFormat(payload *Payload, book *dusk.Book) error {
	switch payload.Extension {
	case ".epub":
		return s.UploadEpub(payload, book)
	default:
		return s.UploadOtherFormat(payload, book)
	}
}

// Upload EPUB format for new book
func (s *Service) UploadNewEpub(payload *Payload) (*dusk.Book, error) {
	ep, err := epub.NewFromReader(payload.File, payload.Size)
	if err != nil && !errors.Is(err, epub.ErrNoCovers) {
		return nil, fmt.Errorf("failed to parse epub file: %w", err)
	}

	book := ep.ToBook()
	errMap := book.Valid()
	if len(errMap) > 0 {
		return nil, errMap
	}

	if err := s.UploadOtherFormat(payload, book); err != nil {
		return nil, err
	}

	// find and upload cover image in epub
	if ep.CoverFile != "" {
		coverFile, err := ep.Open(ep.CoverFile)
		if err != nil {
			return nil, err
		}
		defer coverFile.Close()

		bookDir := filepath.Join(s.Directory, book.SafeTitle())
		filename := fmt.Sprintf("%s%s", coverFilename, filepath.Ext(ep.CoverFile))
		fullPath := filepath.Join(bookDir, filename)
		if err = s.UploadFile(coverFile, fullPath); err != nil {
			return nil, err
		}
		book.Cover = null.StringFrom(getRelativePath(fullPath))
	}

	return book, nil
}

// Upload EPUB format for existing book
func (s *Service) UploadEpub(payload *Payload, book *dusk.Book) error {
	ep, err := epub.NewFromReader(payload.File, payload.Size)
	if err != nil && !errors.Is(err, epub.ErrNoCovers) {
		return fmt.Errorf("failed to parse epub file: %w", err)
	}

	if err := s.UploadOtherFormat(payload, book); err != nil {
		return err
	}

	if ep.CoverFile != "" {
		coverFile, err := ep.Open(ep.CoverFile)
		if err != nil {
			return err
		}
		defer coverFile.Close()

		bookDir := filepath.Join(s.Directory, book.SafeTitle())
		filename := fmt.Sprintf("%s%s", coverFilename, filepath.Ext(ep.CoverFile))
		fullPath := filepath.Join(bookDir, filename)
		if err = s.UploadFile(coverFile, fullPath); err != nil {
			return err
		}
		book.Cover = null.StringFrom(getRelativePath(fullPath))
	}

	return nil
}

// Upload other format for existing book
func (s *Service) UploadOtherFormat(payload *Payload, book *dusk.Book) error {
	bookDir, err := s.createBookDirectory(book)
	if err != nil {
		return err
	}

	filename := fmt.Sprintf("%s%s", book.SafeTitle(), payload.Extension)
	fullPath := filepath.Join(bookDir, filename)
	if err := s.UploadFile(payload.File, fullPath); err != nil {
		return err
	}

	book.Formats = append(book.Formats, getRelativePath(fullPath))
	return nil
}

// Upload book cover for existing book
func (s *Service) UploadBookCover(payload *Payload, book *dusk.Book) error {
	bookDir, err := s.createBookDirectory(book)
	if err != nil {
		return err
	}

	filename := fmt.Sprintf("%s%s", coverFilename, payload.Extension)
	fullPath := filepath.Join(bookDir, filename)
	if err := s.UploadFile(payload.File, fullPath); err != nil {
		return err
	}

	book.Cover = null.StringFrom(getRelativePath(fullPath))
	return nil
}

// Upload book cover from URL for existing book
func (s *Service) UploadBookCoverFromUrl(url string, book *dusk.Book) error {
	bookDir, err := s.createBookDirectory(book)
	if err != nil {
		return err
	}

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("file: failed to fetch file from url: %w", err)
	}
	defer resp.Body.Close()

	ext := path.Ext(path.Base(resp.Request.URL.Path))
	filename := fmt.Sprintf("%s%s", coverFilename, ext)
	fullPath := filepath.Join(bookDir, filename)
	if err := s.UploadFile(resp.Body, fullPath); err != nil {
		return err
	}

	book.Cover = null.StringFrom(getRelativePath(fullPath))
	return nil
}

// Upload file to path
func (s *Service) UploadFile(file io.Reader, path string) error {
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
	bookDir := filepath.Join(s.Directory, book.SafeTitle())
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
