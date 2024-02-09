package epub

import (
	"archive/zip"
	"dusk"
	"path/filepath"
	"testing"

	"github.com/matryer/is"
)

const testData = "../testdata"

var (
	EPUB30_SPEC  = filepath.Join(testData, "epub30-spec.epub")
	NO_CONTAINER = filepath.Join(testData, "noContainer.epub")
	NO_ROOTFILE  = filepath.Join(testData, "noRootFile.epub")
	NO_COVERFILE = filepath.Join(testData, "noCover.epub")
	COVERFILE    = filepath.Join(testData, "diffCover.epub")
)

func TestNew(t *testing.T) {
	is := is.New(t)
	want := &Epub{
		Version: 3,
		metadata: metadata{
			Title:       "EPUB 3.0 Specification",
			Creator:     []string{"EPUB 3 Working Group"},
			Language:    "en",
			Identifiers: []string{"code.google.com.epub-samples.epub30-spec"},
		},

		RootFile:  "EPUB/package.opf",
		CoverFile: "EPUB/img/epub_logo_color.jpg",
	}

	got, err := New(EPUB30_SPEC)
	is.NoErr(err)
	is.Equal(got.Version, want.Version)
	is.Equal(got.metadata, want.metadata)
	is.Equal(got.RootFile, want.RootFile)
	is.Equal(got.CoverFile, want.CoverFile)
}

func TestToBook(t *testing.T) {
	is := is.New(t)
	want := &dusk.Book{
		Title:  "EPUB 3.0 Specification",
		Author: []string{"EPUB 3 Working Group"},
		ISBN:   "code.google.com.epub-samples.epub30-spec",
	}

	ep, err := New(EPUB30_SPEC)
	is.NoErr(err)

	got := ep.ToBook()
	is.Equal(got, want)
}

func TestNoContainerFile(t *testing.T) {
	is := is.New(t)

	rc, err := zip.OpenReader(NO_CONTAINER)
	is.NoErr(err)
	defer rc.Close()

	got := &Epub{ReadCloser: rc}
	err = got.getRootFile()
	is.Equal(err, ErrNotValidEpub)
}

func TestGetRootFile(t *testing.T) {
	is := is.New(t)
	want := "EPUB/package.opf"

	rc, err := zip.OpenReader(EPUB30_SPEC)
	is.NoErr(err)
	defer rc.Close()

	got := &Epub{ReadCloser: rc}
	err = got.getRootFile()
	is.NoErr(err)
	is.Equal(got.RootFile, want)
}

func TestNoRootFile(t *testing.T) {
	is := is.New(t)

	rc, err := zip.OpenReader(NO_ROOTFILE)
	is.NoErr(err)
	defer rc.Close()

	got := &Epub{ReadCloser: rc}
	err = got.getRootFile()
	is.Equal(err, ErrNoRootFiles)
}

func TestGetMetadata(t *testing.T) {
	is := is.New(t)
	version := 3
	want := metadata{
		Title:       "EPUB 3.0 Specification",
		Creator:     []string{"EPUB 3 Working Group"},
		Identifiers: []string{"code.google.com.epub-samples.epub30-spec"},
		Language:    "en",
	}

	rc, err := zip.OpenReader(EPUB30_SPEC)
	is.NoErr(err)
	defer rc.Close()

	got := &Epub{ReadCloser: rc}
	p, err := got.getPackage()
	is.NoErr(err)

	err = got.getMetadata(p)
	is.NoErr(err)
	is.Equal(got.metadata, want)
	is.Equal(got.Version, version)
}

func TestGetCover(t *testing.T) {
	is := is.New(t)
	want := "EPUB/img/epub_logo_color.jpg"

	rc, err := zip.OpenReader(EPUB30_SPEC)
	is.NoErr(err)
	defer rc.Close()

	got := &Epub{ReadCloser: rc}
	p, err := got.getPackage()
	is.NoErr(err)

	err = got.getCover(p)
	is.NoErr(err)
	is.Equal(got.CoverFile, want)
}

func TestNoCover(t *testing.T) {
	is := is.New(t)

	rc, err := zip.OpenReader(NO_COVERFILE)
	is.NoErr(err)
	defer rc.Close()

	got := &Epub{ReadCloser: rc}
	p, err := got.getPackage()
	is.NoErr(err)

	err = got.getCover(p)
	is.Equal(err, ErrNoCovers)
}

func TestGetCoverInMetadata(t *testing.T) {
	is := is.New(t)
	want := "EPUB/cover.jpg"

	rc, err := zip.OpenReader(COVERFILE)
	is.NoErr(err)
	defer rc.Close()

	got := &Epub{ReadCloser: rc}
	p, err := got.getPackage()
	is.NoErr(err)

	err = got.getCover(p)
	is.NoErr(err)
	is.Equal(got.CoverFile, want)
}
