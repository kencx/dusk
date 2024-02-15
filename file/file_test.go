package file

import (
	"errors"
	"io"
	"os"
	"path"
	"testing"

	"github.com/kencx/dusk"
	"github.com/matryer/is"
)

var (
	testFileService *Service
)

func TestUploadFile(t *testing.T) {
	is := is.New(t)
	tempDir := t.TempDir()
	testFileService, err := NewService(tempDir)
	is.NoErr(err)

	f, err := os.Open("../testdata/epub30-spec.epub")
	is.NoErr(err)
	defer f.Close()

	err = testFileService.UploadFile(f, path.Join(tempDir, "foo.epub"))
	is.NoErr(err)

	// check uploaded file
	_, err = os.Stat(path.Join(tempDir, "foo.epub"))
	is.NoErr(err)
}

func TestUploadFileAlreadyExists(t *testing.T) {
	is := is.New(t)
	tempDir := t.TempDir()
	testFileService, err := NewService(tempDir)
	is.NoErr(err)

	// upload file
	f, err := os.Open("../testdata/epub30-spec.epub")
	is.NoErr(err)
	defer f.Close()

	dest, err := os.OpenFile(path.Join(tempDir, "foo.epub"), os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
	is.NoErr(err)
	defer dest.Close()

	_, err = io.Copy(dest, f)
	is.NoErr(err)

	// upload file again
	err = testFileService.UploadFile(f, path.Join(tempDir, "foo.epub"))
	is.True(errors.Is(err, os.ErrExist))
}

func TestCreateBookDirectory(t *testing.T) {
	is := is.New(t)
	tempDir := t.TempDir()
	testFileService, err := NewService(tempDir)
	is.NoErr(err)

	book := &dusk.Book{
		Title: "foobar test",
	}

	want := path.Join(tempDir, "foobar-test")
	got, err := testFileService.createBookDirectory(book)
	is.NoErr(err)
	is.Equal(got, want)
}

func TestCreateBookDirectoryAlreadyExists(t *testing.T) {
	is := is.New(t)

	// create directory
	tempDir := t.TempDir()
	book := &dusk.Book{
		Title: "foobar test",
	}
	tempPath := path.Join(tempDir, book.SafeTitle())
	err := os.MkdirAll(tempPath, 0755)
	is.NoErr(err)

	testFileService, err := NewService(tempDir)
	is.NoErr(err)

	// recreate directory
	want := path.Join(tempDir, "foobar-test")
	got, err := testFileService.createBookDirectory(book)
	is.NoErr(err)
	is.Equal(got, want)
}

func TestGetRelativePath(t *testing.T) {
	is := is.New(t)

	want := "bar/filename.txt"
	got := getRelativePath("tmp/foo/bar/filename.txt")
	is.Equal(got, want)
}
