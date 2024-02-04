package openlibrary

import (
	"dusk"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

const (
	isbnEndpoint      = "https://openlibrary.org/isbn/%s.json"
	authorEndpoint    = "https://openlibrary.org%s.json"
	searchEndpoint    = "https://openlibrary.org/search.json?%s&fields=title,author_name,isbn,publisher,cover_i,publish_date,edition_count"
	coverIdEndpoint   = "https://covers.openlibrary.org/b/id/%s-%s.jpg"
	coverIsbnEndpoint = "https://covers.openlibrary.org/b/isbn/%s-%s.jpg"
)

type MetadataJson struct {
	Title         string `json:"title"`
	Subtitle      string `json:"subtitle,omitempty"`
	NumberOfPages int    `json:"number_of_pages"`
	Authors       []struct {
		Key string
	} `json:"authors"`

	Isbn10      []string            `json:"isbn_10,omitempty"`
	Isbn13      []string            `json:"isbn_13,omitempty"`
	Identifiers map[string][]string `json:"identifiers,omitempty"`

	Series []string `json:"series,omitempty"`

	Publishers  []string `json:"publishers"`
	PublishDate string   `json:"publish_date"`
	Covers      []int    `json:"covers"`
}

type QueryJson struct {
	start    int `json:"start"`
	numFound int `json:"num_found"`
	results  []struct {
		Title       string   `json:"title"`
		Authors     []string `json:"author_name"`
		Isbn        []string `json:"isbn"`
		Publishers  []string `json:"publisher"`
		PublishDate []string `json:"publish_date"`
		CoverId     int      `json:"cover_i`
	} `json:"docs"`
}

type Metadata struct {
	Title         string
	Subtitle      string
	Isbn10        []string
	Isbn13        []string
	Authors       []string
	NumberOfPages int
	Series        []string
	PublishDate   string
	Publishers    []string
	CoverUrl      string
}

type QueryResults []Metadata

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

func FetchByIsbn(isbn string) (*Metadata, error) {
	// TODO add error types

	url := fmt.Sprintf(isbnEndpoint, isbn)

	var mj MetadataJson
	err := fetch(url, &mj)
	if err != nil {
		return nil, err
	}

	res := &Metadata{
		Title:         mj.Title,
		Subtitle:      mj.Subtitle,
		Isbn10:        mj.Isbn10,
		Isbn13:        mj.Isbn13,
		NumberOfPages: mj.NumberOfPages,
		Series:        mj.Series,
		PublishDate:   mj.PublishDate,
		Publishers:    mj.Publishers,
	}

	var authors []string
	for _, a := range mj.Authors {
		authorUrl := fmt.Sprintf(authorEndpoint, a.Key)
		var author struct {
			Name string `json:"name"`
		}

		err := fetch(authorUrl, &author)
		if err != nil {
			return nil, err
		}
		authors = append(authors, author.Name)
	}
	res.Authors = authors

	if len(mj.Covers) > 0 {
		res.CoverUrl = fmt.Sprintf(coverIdEndpoint, strconv.Itoa(mj.Covers[0]), "M")
	}

	return res, nil
}

func FetchByQuery(query string) (*QueryResults, error) {
	url := fmt.Sprintf(searchEndpoint, query)

	var qj QueryJson
	err := fetch(url, &qj)
	if err != nil {
		return nil, err
	}

	var results QueryResults
	for _, res := range qj.results {
		m := &Metadata{
			Title:      res.Title,
			Authors:    res.Authors,
			Isbn10:     res.Isbn,
			Publishers: res.Publishers,
		}
		m.CoverUrl = fmt.Sprintf(coverIdEndpoint, strconv.Itoa(res.CoverId), "M")
		results = append(results, *m)
	}

	return &results, nil
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
