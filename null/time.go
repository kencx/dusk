package null

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

type Time struct {
	sql.NullTime
}

func NewTime(t time.Time, valid bool) Time {
	return Time{
		NullTime: sql.NullTime{
			Time:  t,
			Valid: valid,
		},
	}
}

func TimeFrom(t time.Time) Time {
	if t.IsZero() {
		return NewTime(t, false)
	}
	return NewTime(t, true)
}

func TimeFromPtr(t *time.Time) Time {
	if t == nil {
		return NewTime(time.Time{}, false)
	}
	if (*t).IsZero() {
		return NewTime(*t, false)
	}
	return NewTime(*t, true)
}

func (t Time) ValueOrZero() time.Time {
	if !t.Valid || t.Time.IsZero() {
		return time.Time{}
	}
	return t.Time
}

func (t Time) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return []byte("null"), nil
	}
	return t.Time.MarshalJSON()
}

func (t *Time) UnmarshalJSON(data []byte) error {
	if len(data) > 0 && data[0] == 'n' {
		t.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &t.Time); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	t.Valid = true
	return nil
}

func (t Time) Equal(b Time) bool {
	return (t.Valid == b.Valid && time.Time.Equal(t.ValueOrZero(), b.ValueOrZero()))
}
