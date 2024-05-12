package openlibrary

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/kencx/dusk/integration"
)

type OlMetadata struct {
	integration.Metadata
}

func (m *OlMetadata) UnmarshalJSON(buf []byte) error {
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
