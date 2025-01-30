package googlebooks

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/kencx/dusk/filters"
	"github.com/kencx/dusk/integration"
	"github.com/kencx/dusk/page"
)

const (
	isbnEndpoint = "https://www.googleapis.com/books/v1/volumes?q=isbn:%s"

	// https://developers.google.com/books/docs/v1/performance#partial-response
	searchEndpoint = "https://www.googleapis.com/books/v1/volumes?q=%s&%s&%s"
	searchFields   = "fields=totalItems,items(id,selfLink,volumeInfo(title,subtitle,authors,publisher,publishedDate,description,industryIdentifiers,pageCount,imageLinks,language,infoLink))"
	// maxResults must be 0 < x < 40
	searchLimit = "startIndex=%d&maxResults=%d"

	coverFields = "fields=volumeInfo(imageLinks)"

	clientTimeout = 5 * time.Second
)

type Fetcher struct{}

func (f *Fetcher) GetName() string {
	return "Googlebooks"
}

func (f *Fetcher) FetchByIsbn(isbn string) (*page.Page[integration.Metadata], error) {
	url := fmt.Sprintf(isbnEndpoint, isbn)
	var m GbMetadata

	err := fetch(url, &m)
	if err != nil {
		return nil, fmt.Errorf("[googlebooks] failed to fetch by isbn: %w", err)
	}

	final := page.Single(nil, m.Metadata)
	return final, nil
}

func (f *Fetcher) FetchByQuery(filters *filters.Search, query string) (*page.Page[integration.Metadata], error) {
	query = url.QueryEscape(query)
	searchPage := fmt.Sprintf(searchLimit, filters.AfterId, 30)
	url := fmt.Sprintf(searchEndpoint, query, searchFields, searchPage)
	var results GbQueryResults

	slog.Debug("[googlebooks] Fetching query", slog.String("url", url))

	err := fetch(url, &results)
	if err != nil {
		return nil, fmt.Errorf("[googlebooks] failed to fetch by query: %w", err)
	}

	final := page.New(results.TotalCount, filters.AfterId, filters.AfterId+len(results.Items), &filters.Base, results.Items)
	if filters.Search != "" {
		final.QueryParams.Add("q", filters.Search)
	}
	return final, nil
}

func FetchCover(volumeLink string) (string, error) {
	var coverJson struct {
		VolumeInfo struct {
			ImageLinks struct {
				ThumbNail string `json:"thumbnail"`
				Small     string `json:"small"`
				Medium    string `json:"medium"`
				Large     string `json:"large"`
			} `json:"imageLinks"`
		} `json:"volumeInfo"`
	}

	coverLink := fmt.Sprintf("%s?%s", volumeLink, coverFields)
	err := fetch(coverLink, &coverJson)
	if err != nil {
		return "", err
	}

	image := coverJson.VolumeInfo.ImageLinks
	if image.Medium != "" {
		return image.Medium, nil
	} else if image.Small != "" {
		return image.Small, nil
	} else if image.Large != "" {
		return image.Large, nil
	} else {
		return "", nil
	}
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

func (m *GbMetadata) getIdentifiers(vol Volume) {
	for _, id := range vol.IndustryIdentifiers {
		switch id.Type {
		case "ISBN_10":
			m.Isbn10 = append(m.Isbn10, id.Identifier)
		case "ISBN_13":
			m.Isbn13 = append(m.Isbn13, id.Identifier)
		case "OTHER":
			temp := strings.Split(id.Identifier, ":")
			if len(temp) == 2 {
				t, id := temp[0], temp[1]

				_, ok := m.Identifiers[t]
				if !ok {
					m.Identifiers[t] = []string{id}
				} else {
					m.Identifiers[t] = append(m.Identifiers[t], id)
				}
			}
		default:
			_, ok := m.Identifiers[id.Type]
			if !ok {
				m.Identifiers[id.Type] = []string{id.Identifier}
			} else {
				m.Identifiers[id.Type] = append(m.Identifiers[id.Type], id.Identifier)
			}
		}
	}
}
