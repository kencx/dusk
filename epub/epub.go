package epub

import (
	"archive/zip"
	"dusk"
	"dusk/util"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"slices"
)

const (
	CONTAINER = "META-INF/container.xml"
)

var (
	coverExtension = []string{".jpg", ".jpeg", ".png"}

	ErrNoRootFiles = errors.New("no rootfiles found")
	ErrNoCovers    = errors.New("no cover files found")
)

func Open(path string) (*zip.ReadCloser, error) {
	return zip.OpenReader(path)
}

func ExtractCover(path string) (io.ReadCloser, error) {
	rc, err := zip.OpenReader(path)
	if err != nil {
		return nil, fmt.Errorf("failed to unzip epub: %v", err)
	}
	defer rc.Close()

	cover, err := getCover(rc)
	if err != nil {
		return nil, err
	}

	return rc.Open(cover)
}

type Epub struct {
	*metadata

	RootFile, CoverFile string
}

type container struct {
	Container xml.Name `xml:"container"`
	Rootfiles []struct {
		Rootfile  xml.Name `xml:"rootfile"`
		FullPath  string   `xml:"full-path,attr"`
		MediaType string   `xml:"media-type,attr"`
	} `xml:"rootfiles>rootfile"`
}

type Package struct {
	Package  xml.Name `xml:"package"`
	Metadata metadata `xml:"metadata"`
}

type metadata struct {
	Title       string   `xml:"title"`
	Creator     []string `xml:"creator"`
	Identifiers []string `xml:"identifier"`
	Language    string   `xml:"language"`
	Description string   `xml:"description,omitempty"`
	Date        string   `xml:"date,omitempty"`
	Publisher   string   `xml:"publisher,omitempty"`
}

type identifiers struct{}

func New(path string) (*Epub, error) {
	rc, err := zip.OpenReader(path)
	if err != nil {
		return nil, fmt.Errorf("failed to unzip epub: %v", err)
	}
	defer rc.Close()

	return NewFromFile(rc)
}

func NewFromFile(rc *zip.ReadCloser) (*Epub, error) {
	rootFile, err := getRootFile(rc)
	if err != nil {
		return nil, fmt.Errorf("failed to get rootFile: %v", err)
	}

	m, err := getMetadata(rc, rootFile)
	if err != nil {
		return nil, fmt.Errorf("failed to extract metadata: %v", err)
	}

	cover, err := getCover(rc)
	if err != nil {
		return nil, fmt.Errorf("failed to extract cover: %v", err)
	}

	return &Epub{m, rootFile, cover}, nil
}

func (e *Epub) ToBook() *dusk.Book {
	return &dusk.Book{
		Title:  e.Title,
		Author: e.Creator,
		ISBN:   e.Identifiers[0],
	}
}

func getRootFile(rc *zip.ReadCloser) (string, error) {
	f, err := rc.Open(CONTAINER)
	if err != nil {
		return "", fmt.Errorf("failed to open container.xml: %v", err)
	}
	defer f.Close()

	c := &container{}
	err = util.UnmarshalXml(f, c)
	if err != nil {
		return "", err
	}

	if len(c.Rootfiles) >= 1 {
		return c.Rootfiles[0].FullPath, nil
	}

	return "", ErrNoRootFiles
}

func getMetadata(rc *zip.ReadCloser, rootFile string) (*metadata, error) {
	f, err := rc.Open(rootFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	dest := &Package{}
	err = util.UnmarshalXml(f, dest)
	if err != nil {
		return nil, err
	}

	return &dest.Metadata, nil
}

func getCover(rc *zip.ReadCloser) (string, error) {
	for _, f := range rc.File {
		if slices.Contains(coverExtension, filepath.Ext(f.Name)) {
			return f.Name, nil
		}
	}

	return "", ErrNoCovers
}
