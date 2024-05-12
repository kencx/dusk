package integration

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/araddon/dateparse"
	"github.com/kencx/dusk"
)

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

type QueryResults []*Metadata

func (m *Metadata) ToBook() *dusk.Book {
	var (
		isbn, isbn13, publisher string
		identifiers             = make(map[string]string)
		tags                    []string
		datePublished           time.Time
	)

	isbn = GetFirst(m.Isbn10)
	isbn13 = GetFirst(m.Isbn13)
	publisher = GetFirst(m.Publishers)

	if len(m.Series) > 0 {
		series := strings.ReplaceAll(m.Series[0], ",", "")
		tags = append(tags, fmt.Sprintf("series.%s", series))
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
		isbn, isbn13,
		m.NumberOfPages, 0, 0,
		publisher, "", "", m.CoverUrl,
		m.Authors, tags, nil,
		identifiers,
		datePublished, time.Time{}, time.Time{},
	)
}

func GetFirst(sl []string) string {
	if len(sl) > 0 {
		return sl[0]
	}
	return ""
}
