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

// func FromJSON(body []byte, v interface{}) error {
// 	err := json.Unmarshal(body, v)
// 	if err != nil {
// 		return fmt.Errorf("unable to unmarshal: %w", err)
// 	}
// 	return nil
// }
