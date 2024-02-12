package util

import (
	"encoding/json"
	"fmt"
)

func ToJSON(v interface{}) ([]byte, error) {
	res, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return []byte(""), fmt.Errorf("unable to marshal: %w", err)
	}
	return res, nil
}
