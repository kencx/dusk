package file

import (
	"os"
	"path/filepath"

	"github.com/kencx/dusk"
)

// The full flow of deleting/archiving a book is:
// 1. Delete book from database
// 2. Delete book directory from filesystem
// 3. If step 1 fails, worker should identify orphan directories and flag them
func (s *Service) DeleteBook(book *dusk.Book) error {
	bookDir, err := s.getBookDirectory(book)
	if err != nil {
		return err
	}

	if err := os.RemoveAll(bookDir); err != nil {
		return err
	}
	return nil
}

func (s *Service) ArchiveBook(book *dusk.Book) error {
	bookDir, err := s.getBookDirectory(book)
	if err != nil {
		return err
	}

	if err := os.Rename(bookDir, filepath.Join(s.Directory, s.Archive, book.SafeTitle())); err != nil {
		return err
	}
	return nil
}
