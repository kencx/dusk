package googlebooks

import (
	"encoding/json"
	"log/slog"

	"github.com/kencx/dusk/integration"
)

type GbQueryResults []*integration.Metadata

func (q *GbQueryResults) UnmarshalJSON(buf []byte) error {
	var qj QueryJson

	if err := json.Unmarshal(buf, &qj); err != nil {
		return err
	}

	for _, item := range qj.Items {
		vol := item.VolumeInfo

		if vol.Title == "" || len(vol.Authors) == 0 {
			slog.Debug("volume has no title or authors, skipping...")
			continue
		}

		m := &GbMetadata{
			Metadata: integration.Metadata{
				Title:         vol.Title,
				Subtitle:      vol.Subtitle,
				Authors:       vol.Authors,
				NumberOfPages: vol.NumberOfPages,
				Publishers:    []string{vol.Publisher},
				PublishDate:   vol.PublishDate,
				Identifiers:   make(map[string][]string),
			},
		}

		// when querying, only get thumbnails
		if vol.ImageLinks.ThumbNail != "" {
			m.CoverUrl = vol.ImageLinks.ThumbNail
		} else {
			m.CoverUrl = vol.ImageLinks.SmallThumbNail
		}

		m.getIdentifiers(vol)
		*q = append(*q, &m.Metadata)
	}
	return nil
}
