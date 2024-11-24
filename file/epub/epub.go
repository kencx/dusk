package epub

import (
	"archive/zip"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/araddon/dateparse"
	"github.com/kencx/dusk"
	"github.com/kencx/dusk/util"
)

const (
	CONTAINER = "META-INF/container.xml"
)

var (
	coverExtension = []string{".jpg", ".jpeg", ".png"}

	ErrNotValidEpub = errors.New("not valid epub file")
	ErrNoRootFiles  = errors.New("no root files found")
	ErrNoCovers     = errors.New("no cover files found")
)

func ExtractCoverFile(path string) (io.ReadCloser, error) {
	ep, err := New(path)
	if err != nil {
		return nil, err
	}

	return ep.Open(ep.CoverFile)
}

type Epub struct {
	*zip.Reader
	Version int
	metadata

	// rel path from EPUB root
	RootFile  string
	CoverFile string
}

type container struct {
	Container xml.Name `xml:"container"`
	Rootfiles []struct {
		Rootfile  xml.Name `xml:"rootfile"`
		FullPath  string   `xml:"full-path,attr"`
		MediaType string   `xml:"media-type,attr"`
	} `xml:"rootfiles>rootfile"`
}

type contentPackage struct {
	Package  xml.Name `xml:"package"`
	Version  string   `xml:"version,attr"`
	Metadata metadata `xml:"metadata"`
	Manifest manifest `xml:"manifest"`
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

type manifest struct {
	Item []struct {
		Item       xml.Name `xml:"item"`
		Href       string   `xml:"href,attr"`
		Id         string   `xml:"id,attr"`
		MediaType  string   `xml:"media-type,attr"`
		Properties string   `xml:"properties,attr,omitempty"`
	} `xml:"item"`
}

func New(path string) (*Epub, error) {
	rc, err := zip.OpenReader(path)
	if err != nil {
		return nil, fmt.Errorf("epub: failed to unzip epub: %w", err)
	}
	defer rc.Close()

	return new(&rc.Reader)
}

func NewFromReader(r multipart.File, fileSize int64) (*Epub, error) {
	rc, err := zip.NewReader(r, fileSize)
	if err != nil {
		return nil, err
	}

	return new(rc)
}

func new(r *zip.Reader) (*Epub, error) {
	ep := &Epub{Reader: r}
	if err := ep.getRootFile(); err != nil {
		return nil, fmt.Errorf("epub: failed to extract rootFile: %w", err)
	}

	p, err := ep.getPackage()
	if err != nil {
		return nil, fmt.Errorf("epub: failed to extract package: %w", err)
	}

	err = ep.getMetadata(p)
	if err != nil {
		return nil, fmt.Errorf("epub: failed to extract metadata: %w", err)
	}

	err = ep.getCover(p)
	if err != nil {
		if errors.Is(err, ErrNoCovers) {
			slog.Debug("[epub] No valid cover files found", slog.String("title", ep.Title))
			return ep, ErrNoCovers
		}
		return nil, fmt.Errorf("epub: failed to extract cover: %w", err)
	}

	return ep, nil
}

func (e *Epub) ToBook() *dusk.Book {
	var (
		isbn10        []string
		isbn13        []string
		datePublished time.Time
		err           error
	)

	for _, id := range e.Identifiers {
		i, err := util.IsbnExtract(id)
		if err != nil {
			continue
		}

		if len(i) == 10 {
			isbn10 = append(isbn10, i)
		} else if len(i) == 13 {
			isbn13 = append(isbn13, i)
		}
	}

	datePublished, err = dateparse.ParseAny(e.Date)
	if err != nil {
		datePublished = time.Time{}
	}

	return dusk.NewBook(
		e.Title, "",
		e.Creator, nil, nil,
		isbn10, isbn13,
		0, 0, 0, 0,
		e.Publisher, "", e.Description, "", "",
		datePublished, time.Time{}, time.Time{},
	)
}

func (e *Epub) getRootFile() error {
	f, err := e.Open(CONTAINER)
	if errors.Is(err, os.ErrNotExist) {
		return ErrNotValidEpub
	} else if err != nil {
		return fmt.Errorf("failed to open container.xml: %w", err)
	}
	defer f.Close()

	c := &container{}
	err = util.UnmarshalXml(f, c)
	if err != nil {
		return err
	}

	for _, rootFile := range c.Rootfiles {
		if slices.ContainsFunc[[]*zip.File](e.File, func(z *zip.File) bool {
			return z.Name == rootFile.FullPath
		}) {
			e.RootFile = rootFile.FullPath
			return nil
		}
	}

	return ErrNoRootFiles
}

func (e *Epub) getMetadata(p *contentPackage) error {
	v, err := strconv.ParseFloat(p.Version, 64)
	if err != nil {
		return err
	}

	e.Version = int(v)
	e.metadata = p.Metadata
	return nil
}

func (e *Epub) getCover(p *contentPackage) error {
	for _, item := range p.Manifest.Item {
		if item.Properties == "cover-image" || strings.Contains(item.Id, "cover") {
			// handle only image files
			if strings.Contains(item.MediaType, "image/") || slices.Contains(coverExtension, filepath.Ext(item.Href)) {
				// The cover-image property returns a path that is relative to
				// the root file. Thus, we prefix it with with the root file's
				// parent directory to get the absolute path from the EPUB root.
				e.CoverFile = filepath.Join(filepath.Dir(e.RootFile), item.Href)
				return nil
			}
		}
	}

	// fallback to any image file
	for _, f := range e.File {
		if slices.Contains(coverExtension, filepath.Ext(f.Name)) {
			e.CoverFile = f.Name
			return nil
		}
	}
	return ErrNoCovers
}

func (e *Epub) getPackage() (*contentPackage, error) {
	if e.RootFile == "" {
		if err := e.getRootFile(); err != nil {
			return nil, err
		}
	}

	f, err := e.Open(e.RootFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open root file: %w", err)
	}
	defer f.Close()

	cp := &contentPackage{}
	err = util.UnmarshalXml(f, cp)
	if err != nil {
		return nil, err
	}

	return cp, nil
}
