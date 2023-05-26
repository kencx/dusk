package util

import (
	"database/sql"
	"encoding/json"
	"strings"
)

type NullString struct {
	sql.NullString
}

func (n NullString) Split() []string {
	if n.Valid {
		return strings.Split(n.String, ",")
	}
	return nil
}

func (n NullString) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.String)
	}
	return json.Marshal(nil)
}

func (n NullString) UnmarshalJSON(data []byte) error {
	var s *string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	if s != nil {
		n.Valid = true
		n.String = *s
	} else {
		n.Valid = false
	}
	return nil
}
