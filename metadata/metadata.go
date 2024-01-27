package metadata

import (
	"dusk"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	isbnString        = "ISBN:%s"
	openLibraryApiURL = "https://openlibrary.org/api/books?bibkeys=%s&jscmd=data&format=json"
)

type Metadata struct {
	Title         string
	Subtitle      string
	Isbn10        []string
	Isbn13        []string
	Authors       []string
	Publishers    []string
	NumberOfPages int
	PublishDate   string
	CoverUrl      string
}

func (m *Metadata) ToBook() *dusk.Book {
	b := &dusk.Book{
		Title:      m.Title,
		Author:     m.Authors,
		NumOfPages: m.NumberOfPages,
	}

	if len(m.Isbn10) > 0 {
		b.ISBN = m.Isbn10[0]
	}
	return b
}

type Content map[string]MetadataJson

type MetadataJson struct {
	Title         string `json:"title"`
	Subtitle      string `json:"subtitle,omitempty"`
	NumberOfPages int    `json:"number_of_pages"`
	PublishDate   string `json:"publish_date"`

	Identifiers map[string][]string `json:"identifiers"`
	Authors     []struct {
		Url  string
		Name string
	} `json:"authors"`

	Publishers []struct {
		Name string
	} `json:"publishers"`

	CoverUrl map[string]string `json:"cover"`
}

func Fetch(isbn string) (*Metadata, error) {

	bibKey := fmt.Sprintf(isbnString, isbn)
	url := fmt.Sprintf(openLibraryApiURL, bibKey)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var content Content
	if err := json.Unmarshal(data, &content); err != nil {
		return nil, err
	}

	mj := content[bibKey]
	return mj.parse(), nil
}

func (m *MetadataJson) parse() *Metadata {
	res := &Metadata{
		Title:         m.Title,
		Subtitle:      m.Subtitle,
		NumberOfPages: m.NumberOfPages,
		PublishDate:   m.PublishDate,
	}

	if i, ok := m.Identifiers["isbn_10"]; ok {
		res.Isbn10 = i
	}

	if i, ok := m.Identifiers["isbn_13"]; ok {
		res.Isbn13 = i
	}

	var authors []string
	for _, a := range m.Authors {
		authors = append(authors, a.Name)
	}
	res.Authors = authors

	var publishers []string
	for _, p := range m.Publishers {
		publishers = append(publishers, p.Name)
	}
	res.Publishers = publishers
	res.CoverUrl = m.CoverUrl["small"]

	return res
}
