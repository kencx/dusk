package googlebooks

import (
	"encoding/json"
	"log/slog"
	"strings"

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

		m := &integration.Metadata{
			Title:         vol.Title,
			Subtitle:      vol.Subtitle,
			Authors:       vol.Authors,
			NumberOfPages: vol.NumberOfPages,
			Publishers:    []string{vol.Publisher},
			PublishDate:   vol.PublishDate,
			Identifiers:   make(map[string][]string),
		}

		if vol.ImageLinks.ThumbNail != "" {
			m.CoverUrl = vol.ImageLinks.ThumbNail
		} else {
			m.CoverUrl = vol.ImageLinks.SmallThumbNail
		}

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
		*q = append(*q, m)
	}
	return nil
}
