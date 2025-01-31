package integration

import (
	"errors"
	"log/slog"
	"time"

	"github.com/araddon/dateparse"
	"github.com/kencx/dusk"
)

var ErrInvalidMetadata = errors.New("invalid metadata")

type Metadata struct {
	Title         string
	Subtitle      string
	Isbn10        []string
	Isbn13        []string
	Identifiers   map[string][]string
	Authors       []string
	NumberOfPages int
	Series        []string
	PublishDate   string
	Publishers    []string
	CoverUrl      string
}

func (m *Metadata) ToBook() *dusk.Book {
	var (
		isbn, isbn13      []string
		publisher, series string
		identifiers       = make(map[string]string)
		tags              []string
		datePublished     time.Time
	)

	isbn = m.Isbn10
	isbn13 = m.Isbn13
	publisher = GetFirst(m.Publishers)

	if len(m.Series) > 0 {
		series = m.Series[0]
	}

	if len(m.Identifiers) > 0 {
		for k, v := range m.Identifiers {
			if len(v) > 0 {
				identifiers[k] = v[0]
			}
		}
	}

	datePublished, err := dateparse.ParseAny(m.PublishDate)
	if err != nil {
		slog.Warn("failed to parse publish date", slog.Any("err", err))
	}

	return dusk.NewBook(
		m.Title, m.Subtitle,
		m.Authors, tags, nil,
		isbn, isbn13,
		m.NumberOfPages, 0, 0, 0,
		publisher, series, "", "", m.CoverUrl,
		datePublished, time.Time{}, time.Time{}, time.Time{},
	)
}

func GetFirst(sl []string) string {
	if len(sl) > 0 {
		return sl[0]
	}
	return ""
}
