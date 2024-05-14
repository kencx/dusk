package googlebooks

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/kencx/dusk/integration"
)

type GbMetadata struct {
	integration.Metadata
}

type QueryJson struct {
	TotalItems int `json:"totalItems"`
	Items      []struct {
		Id         string `json:"id"`
		SelfLink   string `json:"selfLink"`
		VolumeInfo struct {
			Title               string   `json:"title"`
			Subtitle            string   `json:"subtitle"`
			Authors             []string `json:"authors"`
			Publisher           string   `json:"publisher"`
			PublishDate         string   `json:"publishedDate"`
			Description         string   `json:"description"`
			IndustryIdentifiers []struct {
				Type       string `json:"type"`
				Identifier string `json:"identifier"`
			} `json:"IndustryIdentifiers"`
			NumberOfPages int `json:"pageCount"`
			ImageLinks    struct {
				SmallThumbNail string `json:"smallThumbnail"`
				ThumbNail      string `json:"thumbnail"`
			} `json:"imageLinks"`
			Language string `json:"language"`
			InfoLink string `json:"infoLink"`
		} `json:"volumeInfo"`
	} `json:"items"`
}

func (m *GbMetadata) UnmarshalJSON(buf []byte) error {
	var im QueryJson
	if err := json.Unmarshal(buf, &im); err != nil {
		return err
	}

	if len(im.Items) == 0 || im.TotalItems == 0 {
		return nil
	}

	if len(im.Items) > 1 || im.TotalItems > 1 {
		slog.Debug("[googlebooks] more than 1 item fetched for 1 ISBN")
	}

	link := im.Items[0].SelfLink
	vol := im.Items[0].VolumeInfo

	if vol.Title == "" || len(vol.Authors) == 0 {
		return ErrInvalidResult
	}

	m.Title = vol.Title
	m.Subtitle = vol.Subtitle
	m.Authors = vol.Authors
	m.NumberOfPages = vol.NumberOfPages
	m.Publishers = append(m.Publishers, vol.Publisher)
	m.PublishDate = vol.PublishDate

	cover, err := fetchCover(link)
	if err != nil {
		slog.Warn("[googlebooks] failed to fetch cover")
	}
	m.CoverUrl = cover

	for _, id := range vol.IndustryIdentifiers {
		switch id.Type {
		case "ISBN_10":
			m.Isbn10 = append(m.Isbn10, id.Identifier)
		case "ISBN_13":
			m.Isbn13 = append(m.Isbn13, id.Identifier)
		default:
			_, ok := m.Identifiers[id.Type]
			if !ok {
				m.Identifiers[id.Type] = []string{id.Identifier}
			} else {
				m.Identifiers[id.Type] = append(m.Identifiers[id.Type], id.Identifier)
			}
		}
	}
	return nil
}

func fetchCover(volumeLink string) (string, error) {
	type coverJson struct {
		VolumeInfo struct {
			ImageLinks struct {
				ThumbNail string `json:"thumbnail"`
				Small     string `json:"small"`
				Medium    string `json:"medium"`
				Large     string `json:"large"`
			} `json:"imageLinks"`
		} `json:"volumeInfo"`
	}
	var coverResp coverJson

	coverLink := fmt.Sprintf("%s?%s", volumeLink, coverFields)
	err := fetch(coverLink, &coverResp)
	if err != nil {
		return "", err
	}

	// TODO select appropriate cover size
	return coverResp.VolumeInfo.ImageLinks.Medium, nil
}
