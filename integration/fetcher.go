package integration

import (
	"errors"
	"log/slog"

	"github.com/kencx/dusk/filters"
	"github.com/kencx/dusk/page"
)

type Fetcher interface {
	FetchByIsbn(isbn string) (*page.Page[Metadata], error)
	FetchByQuery(filters *filters.Search, query string) (*page.Page[Metadata], error)
	GetName() string
}

type Fetchers []Fetcher

func (fs Fetchers) FetchByIsbn(isbn string) (*page.Page[Metadata], error) {
	var result *page.Page[Metadata]

	for _, f := range fs {
		m, err := f.FetchByIsbn(isbn)
		if err == nil {
			result = m
			break
		}
		slog.Warn("", slog.Any("err", err))
	}

	if result != nil {
		return result, nil
	} else {
		return nil, errors.New("failed to fetch from list of given fetchers")
	}
}

func (fs Fetchers) FetchByQuery(filters *filters.Search, query string) (*page.Page[Metadata], error) {
	var result *page.Page[Metadata]

	for _, f := range fs {
		m, err := f.FetchByQuery(filters, query)
		if err == nil {
			result = m
			break
		}
		slog.Warn("", slog.Any("err", err))
	}

	if result != nil {
		return result, nil
	} else {
		return nil, errors.New("failed to fetch from list of given fetchers")
	}
}
