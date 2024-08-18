package integration

import (
	"errors"
	"log/slog"
)

type Fetcher interface {
	FetchByIsbn(isbn string) (*Metadata, error)
	FetchByQuery(query string) (*QueryResults, error)
	GetName() string
}

type Fetchers []Fetcher

func (fs Fetchers) FetchByIsbn(isbn string) (*Metadata, error) {
	var result *Metadata

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

func (fs Fetchers) FetchByQuery(query string) (*QueryResults, error) {
	var result *QueryResults

	for _, f := range fs {
		m, err := f.FetchByQuery(query)
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
