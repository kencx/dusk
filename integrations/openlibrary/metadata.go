package openlibrary

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

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

func (m *Metadata) UnmarshalJSON(buf []byte) error {
	var im struct {
		Title         string `json:"title"`
		Subtitle      string `json:"subtitle,omitempty"`
		NumberOfPages int    `json:"number_of_pages,omitempty"`
		Authors       []struct {
			Key string
		} `json:"authors"`

		Isbn10      []string            `json:"isbn_10"`
		Isbn13      []string            `json:"isbn_13"`
		Identifiers map[string][]string `json:"identifiers,omitempty"`

		Works []struct {
			Key string
		} `json:"works"`
		Series []string `json:"series,omitempty"`

		Publishers  []string `json:"publishers,omitempty"`
		PublishDate string   `json:"publish_date,omitempty"`
		Covers      []int    `json:"covers"`
	}

	if err := json.Unmarshal(buf, &im); err != nil {
		return err
	}

	// fallback to works metadata if title or authors missing
	if im.Title == "" || len(im.Authors) == 0 || im.Authors == nil {
		if len(im.Works) == 0 {
			return ErrInvalidResult
		}

		var worksMetadata struct {
			Title   string `json:"title"`
			Authors []struct {
				Author struct {
					Key string `json:"key"`
				} `json:"author"`
			} `json:"authors"`
			Description string `json:"description"`
		}

		url := fmt.Sprintf(olEndpoint, im.Works[0].Key)
		if err := fetch(url, &worksMetadata); err != nil {
			return fmt.Errorf("failed to fetch by works: %w", err)
		}

		if im.Title == "" {
			im.Title = worksMetadata.Title
		}

		if len(im.Authors) == 0 || im.Authors == nil {
			for _, a := range worksMetadata.Authors {
				im.Authors = append(im.Authors, struct{ Key string }{a.Author.Key})
			}
		}
	}

	m.Title = im.Title
	m.Subtitle = im.Subtitle
	m.Isbn10 = im.Isbn10
	m.Isbn13 = im.Isbn13
	m.Identifiers = im.Identifiers
	m.NumberOfPages = im.NumberOfPages
	m.Series = im.Series
	m.Publishers = im.Publishers
	m.PublishDate = im.PublishDate

	var authors []string
	for _, a := range im.Authors {
		authorUrl := fmt.Sprintf(olEndpoint, a.Key)
		var author struct {
			Name string `json:"name"`
		}

		if err := fetch(authorUrl, &author); err != nil {
			return fmt.Errorf("failed to fetch author: %w", err)
		}

		authors = append(authors, author.Name)
	}
	m.Authors = authors

	if len(im.Covers) > 0 {
		m.CoverUrl = fmt.Sprintf(coverIdEndpoint, strconv.Itoa(im.Covers[0]), "M")
	}

	return nil
}

func (m *Metadata) ToBook() *dusk.Book {
	var (
		isbn, isbn13, publisher string
		identifiers             map[string]string = make(map[string]string)
		tags                    []string
		datePublished           time.Time
	)

	if len(m.Isbn10) > 0 {
		isbn = m.Isbn10[0]
	}
	if len(m.Isbn13) > 0 {
		isbn13 = m.Isbn13[0]
	}
	if len(m.Publishers) > 0 {
		publisher = m.Publishers[0]
	}
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

	// there is no standardized format so only handle year for now
	if len(m.PublishDate) == 4 {
		var err error
		datePublished, err = time.Parse("2006", m.PublishDate)
		if err != nil {
			slog.Warn("failed to parse publish date", slog.Any("err", err))
		}
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
