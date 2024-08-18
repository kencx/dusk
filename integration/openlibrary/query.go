package openlibrary

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/kencx/dusk"
	"github.com/kencx/dusk/integration"
)

type QueryJson struct {
	Start    int `json:"start"`
	NumFound int `json:"numFound"`
	Results  []struct {
		Key      string `json:"key"`
		Title    string `json:"title"`
		Editions struct {
			Start    int `json:"start"`
			NumFound int `json:"numFound"`
			Results  []struct {
				Title         string   `json:"title"`
				Subtitle      string   `json:"subtitle,omitempty"`
				Authors       []string `json:"author_name"`
				CoverId       int      `json:"cover_i"`
				NumberOfPages int      `json:"number_of_pages_median"`

				// isbn10 and isbn13 of various editions
				Isbn []string `json:"isbn"`

				// multiple publishers of various editions
				Publishers []string `json:"publisher"`

				// both year and full dates
				PublishDate []string `json:"publish_date"`
			} `json:"docs"`
		} `json:"editions"`
	} `json:"docs"`
}

type OlQueryResults []*integration.Metadata

func (q *OlQueryResults) UnmarshalJSON(buf []byte) error {
	var qj QueryJson

	if err := json.Unmarshal(buf, &qj); err != nil {
		return err
	}

	if len(qj.Results) == 0 {
		return dusk.ErrNoRows
	}

	slog.Debug(fmt.Sprintf("[openlibrary] Found %d results", qj.NumFound))
	for _, work := range qj.Results {
		for _, r := range work.Editions.Results {
			m := &integration.Metadata{
				Title:         r.Title,
				Subtitle:      r.Subtitle,
				Authors:       r.Authors,
				Publishers:    r.Publishers,
				NumberOfPages: r.NumberOfPages,
				CoverUrl:      fmt.Sprintf(coverIdEndpoint, strconv.Itoa(r.CoverId), "M"),
			}

			m.PublishDate = integration.GetFirst(r.PublishDate)

			if len(r.Isbn) > 0 {
				for _, i := range r.Isbn {
					if len(i) == 10 {
						m.Isbn10 = append(m.Isbn10, i)
					} else {
						m.Isbn13 = append(m.Isbn13, i)
					}
				}
			}

			// fallback to works with key
			if len(m.Authors) == 0 || m.Authors == nil || m.Title == "" {
				slog.Debug("[openlibrary] result has no title or authors, falling back to works")

				if work.Key == "" {
					slog.Debug("[openlibrary] no work key found, skipping...")
					continue
				}

				var worksMetadata struct {
					Title   string `json:"title"`
					Authors []struct {
						Author struct {
							Key string `json:"key"`
						} `json:"author"`
					} `json:"authors"`
				}

				url := fmt.Sprintf(olEndpoint, work.Key)
				if err := fetch(url, &worksMetadata); err != nil {
					slog.Debug("[openlibrary] failed to fetch by works", slog.Any("err", err))
					continue
				}

				if m.Title == "" {
					m.Title = worksMetadata.Title
				}

				if len(m.Authors) == 0 || m.Authors == nil {
					for _, a := range worksMetadata.Authors {
						authorUrl := fmt.Sprintf(olEndpoint, a.Author.Key)
						var author struct {
							Name string `json:"name"`
						}

						if err := fetch(authorUrl, &author); err != nil {
							slog.Debug("[openlibrary] failed to fetch by author", slog.Any("err", err))
							continue
						}

						m.Authors = append(m.Authors, author.Name)
					}
				}
			}
			*q = append(*q, m)
		}
	}
	return nil
}
