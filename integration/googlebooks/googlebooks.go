package googlebooks

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/kencx/dusk/integration"
)

const (
	isbnEndpoint = "https://www.googleapis.com/books/v1/volumes?q=isbn:%s"

	coverFields = "fields=volumeInfo(imageLinks)"

	// https://developers.google.com/books/docs/v1/performance#partial-response
	searchEndpoint = "https://www.googleapis.com/books/v1/volumes?q=%s&%s&%s"
	searchFields   = "fields=totalItems,items(id,selfLink,volumeInfo(title,subtitle,authors,publisher,publishedDate,description,industryIdentifiers,pageCount,imageLinks,language,infoLink))"
	searchLimit    = "searchIndex=0&maxResults=10"

	clientTimeout = 5 * time.Second
)

var (
	ErrInvalidResult = errors.New("invalid googlebooks result")
)

type Fetcher struct{}

func (f *Fetcher) FetchByIsbn(isbn string) (*integration.Metadata, error) {
	url := fmt.Sprintf(isbnEndpoint, isbn)
	var m GbMetadata

	err := fetch(url, &m)
	if err != nil {
		return nil, fmt.Errorf("[googlebooks] failed to fetch by isbn: %w", err)
	}

	return &m.Metadata, nil
}

func (f *Fetcher) FetchByQuery(query string) (*integration.QueryResults, error) {
	query = url.QueryEscape(query)
	url := fmt.Sprintf(searchEndpoint, query, searchFields, searchLimit)
	var q GbQueryResults

	err := fetch(url, &q)
	if err != nil {
		return nil, fmt.Errorf("[googlebooks] failed to fetch by query: %w", err)
	}

	var res integration.QueryResults
	for _, qr := range q {
		res = append(res, qr)
	}

	return &res, nil
}

func fetch(url string, dest interface{}) error {
	client := http.Client{
		Timeout: clientTimeout,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	// To receive a gzip-encoded response, Google Books API expects
	// the following headers:
	//   1. Accept-Encoding: gzip
	//   2. User-Agent must contain the string gzip
	// See https://developers.google.com/books/docs/v1/performance

	// The request header "Accept-Encoding: gzip" is automatically
	// set, and the response body is automatically decompressed when
	// DisableCompression is true. However, if the user explicitly
	// adds the header manually, the response body is not automatically
	// uncompressed.
	req.Header.Add("User-Agent", "dusk (gzip)")

	resp, err := client.Do(req)
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
