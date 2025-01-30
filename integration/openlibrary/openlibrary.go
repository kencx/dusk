package openlibrary

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/kencx/dusk/filters"
	"github.com/kencx/dusk/integration"
	"github.com/kencx/dusk/page"
)

const (
	olEndpoint   = "https://openlibrary.org%s.json"
	isbnEndpoint = "https://openlibrary.org/isbn/%s.json"

	searchEndpoint = "https://openlibrary.org/search.json?q=%s&%s&%s"
	searchFields   = "fields=key,title,editions,editions.title,editions.subtitle,editions.author_name,editions.isbn,editions.publisher,editions.cover_i,editions.publish_date,editions.number_of_pages_median"
	searchLimit    = "limit=%d&offset=%d"

	coverIdEndpoint   = "https://covers.openlibrary.org/b/id/%s-%s.jpg"
	coverIsbnEndpoint = "https://covers.openlibrary.org/b/isbn/%s-%s.jpg"

	clientTimeout = 5 * time.Second
)

type Fetcher struct{}

func (f *Fetcher) GetName() string {
	return "Openlibrary"
}

func (f *Fetcher) FetchByIsbn(isbn string) (*page.Page[integration.Metadata], error) {
	url := fmt.Sprintf(isbnEndpoint, isbn)
	var m OlMetadata

	err := fetch(url, &m)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch by isbn: %w", err)
	}

	final := page.Single(nil, m.Metadata)
	return final, nil
}

func (f *Fetcher) FetchByQuery(filters *filters.Search, query string) (*page.Page[integration.Metadata], error) {
	query = url.QueryEscape(query)
	searchPage := fmt.Sprintf(searchLimit, 30, filters.AfterId)
	url := fmt.Sprintf(searchEndpoint, query, searchFields, searchPage)
	var results OlQueryResults

	slog.Debug("[openlibrary] Fetching query", slog.String("url", url))

	err := fetch(url, &results)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch by query: %w", err)
	}

	final := page.New(results.TotalCount, filters.AfterId, filters.AfterId+len(results.Items), &filters.Base, results.Items)
	if filters.Search != "" {
		final.QueryParams.Add("q", filters.Search)
	}
	return final, nil
}

// TODO FetchByWork
func (f *Fetcher) FetchByWork() (*page.Page[integration.Metadata], error) {
	return nil, nil
}

func fetch(url string, dest interface{}) error {
	client := http.Client{
		Timeout: clientTimeout,
	}

	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, &dest); err != nil {
		return err
	}
	return nil
}
