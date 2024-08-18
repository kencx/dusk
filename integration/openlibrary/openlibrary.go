package openlibrary

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/kencx/dusk/integration"
)

const (
	olEndpoint   = "https://openlibrary.org%s.json"
	isbnEndpoint = "https://openlibrary.org/isbn/%s.json"

	searchEndpoint = "https://openlibrary.org/search.json?q=%s&%s&%s"
	searchFields   = "fields=key,title,editions,editions.title,editions.subtitle,editions.author_name,editions.isbn,editions.publisher,editions.cover_i,editions.publish_date,editions.number_of_pages_median"
	searchLimit    = "limit=5&offset=0"

	coverIdEndpoint   = "https://covers.openlibrary.org/b/id/%s-%s.jpg"
	coverIsbnEndpoint = "https://covers.openlibrary.org/b/isbn/%s-%s.jpg"

	clientTimeout = 5 * time.Second
)

type Fetcher struct{}

func (f *Fetcher) GetName() string {
	return "Openlibrary"
}

func (f *Fetcher) FetchByIsbn(isbn string) (*integration.Metadata, error) {
	url := fmt.Sprintf(isbnEndpoint, isbn)
	var m OlMetadata

	err := fetch(url, &m)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch by isbn: %w", err)
	}

	return &m.Metadata, nil
}

func (f *Fetcher) FetchByQuery(query string) (*integration.QueryResults, error) {
	query = url.QueryEscape(query)
	url := fmt.Sprintf(searchEndpoint, query, searchFields, searchLimit)
	var q OlQueryResults

	err := fetch(url, &q)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch by query: %w", err)
	}

	var res integration.QueryResults
	for _, qr := range q {
		res = append(res, qr)
	}

	return &res, nil
}

// TODO FetchByWork
func (f *Fetcher) FetchByWork() error {
	return nil
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
