package null

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
)

type String struct {
	sql.NullString
}

func NewString(s string, valid bool) String {
	return String{
		NullString: sql.NullString{
			String: s,
			Valid:  valid,
		},
	}
}

func StringFrom(s string) String {
	return NewString(s, true)
}

func StringFromPtr(s *string) String {
	if s == nil {
		return NewString("", false)
	}
	return NewString(*s, true)
}

func (n String) ValueOrZero() string {
	if !n.Valid {
		return ""
	}
	return n.String
}

func (s *String) UnmarshalJSON(data []byte) error {
	if len(data) > 0 && data[0] == 'n' {
		s.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &s.String); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	s.Valid = true
	return nil
}

func (s String) MarshalJSON() ([]byte, error) {
	if !s.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(s.String)
}

func (n String) Split(sep string) []string {
	if n.Valid {
		return strings.Split(n.String, sep)
	}
	return nil
}
