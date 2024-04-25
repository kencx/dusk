package openlibrary

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	olEndpoint   = "https://openlibrary.org%s.json"
	isbnEndpoint = "https://openlibrary.org/isbn/%s.json"

	coverIdEndpoint   = "https://covers.openlibrary.org/b/id/%s-%s.jpg"
	coverIsbnEndpoint = "https://covers.openlibrary.org/b/isbn/%s-%s.jpg"

	searchEndpoint = "https://openlibrary.org/search.json?q=%s&%s&%s"
	searchFields   = "fields=key,title,editions,editions.title,editions.subtitle,editions.author_name,editions.isbn,editions.publisher,editions.cover_i,editions.publish_date,editions.number_of_pages_median"
	searchLimit    = "limit=5&offset=0"
)

var (
	ErrInvalidResult = errors.New("invalid openlibrary result")
)

func FetchByIsbn(isbn string) (*Metadata, error) {
	url := fmt.Sprintf(isbnEndpoint, isbn)
	var m Metadata

	err := fetch(url, &m)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch by isbn: %w", err)
	}

	return &m, nil
}

func FetchByQuery(query string) (*QueryResults, error) {
	query = url.QueryEscape(query)
	url := fmt.Sprintf(searchEndpoint, query, searchFields, searchLimit)
	var q QueryResults

	err := fetch(url, &q)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch by query: %w", err)
	}

	return &q, nil
}

func fetch(url string, dest interface{}) error {
	resp, err := http.Get(url)
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
