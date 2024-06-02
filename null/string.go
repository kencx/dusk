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
	if s == "" {
		return NewString("", false)
	}

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

func (n *String) UnmarshalJSON(data []byte) error {
	if len(data) > 0 && data[0] == 'n' {
		n.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &n.String); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	n.Valid = true
	return nil
}

func (n String) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(n.String)
}

func (n String) Split(sep string) []string {
	if n.Valid {
		return strings.Split(n.String, sep)
	}
	return nil
}

func (n String) Equal(b String) bool {
	return (n.Valid == b.Valid && n.ValueOrZero() == b.ValueOrZero())
}
