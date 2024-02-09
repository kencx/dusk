package util

import (
	"encoding/xml"
	"fmt"
	"io"
)

func UnmarshalXml(f io.Reader, v interface{}) error {
	decoder := xml.NewDecoder(f)
	decoder.Entity = xml.HTMLEntity

	err := decoder.Decode(v)
	if err != nil {
		return fmt.Errorf("failed to decode xml: %v", err)
	}
	return nil
}
