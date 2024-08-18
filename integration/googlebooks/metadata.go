package googlebooks

import (
	"encoding/json"
	"log/slog"

	"github.com/kencx/dusk"
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
		VolumeInfo Volume `json:"volumeInfo"`
	} `json:"items"`
}

type Volume struct {
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
}

func (m *GbMetadata) UnmarshalJSON(buf []byte) error {
	var im QueryJson
	if err := json.Unmarshal(buf, &im); err != nil {
		return err
	}

	if len(im.Items) == 0 || im.TotalItems == 0 {
		return dusk.ErrNoRows
	}

	if len(im.Items) > 1 || im.TotalItems > 1 {
		slog.Debug("[googlebooks] more than 1 item fetched for 1 ISBN")
	}

	vol := im.Items[0].VolumeInfo

	if vol.Title == "" || len(vol.Authors) == 0 {
		return integration.ErrInvalidMetadata
	}

	m.Title = vol.Title
	m.Subtitle = vol.Subtitle
	m.Authors = vol.Authors
	m.NumberOfPages = vol.NumberOfPages
	m.Publishers = append(m.Publishers, vol.Publisher)
	m.PublishDate = vol.PublishDate
	m.Identifiers = make(map[string][]string)

	cover, err := FetchCover(im.Items[0].SelfLink)
	if err != nil {
		slog.Warn("[googlebooks] failed to fetch cover")
	}
	m.CoverUrl = cover

	m.getIdentifiers(vol)

	return nil
}
