package request

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
)

const (
	maxBytes      = 1_048_576
	maxUploadSize = 1024 * 1024 * 300
)

var (
	syntaxError           *json.SyntaxError
	unmarshalTypeError    *json.UnmarshalTypeError
	invalidUnmarshalError *json.InvalidUnmarshalError
)

func ReadJSON(rw http.ResponseWriter, r *http.Request, dest interface{}) error {

	// limit request body
	r.Body = http.MaxBytesReader(rw, r.Body, int64(maxBytes))

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(dest)

	if err != nil {
		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON at character %d", syntaxError.Offset)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")

		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type at character %d", unmarshalTypeError.Offset)

		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)

		case err.Error() == "http: request body too large":
			return fmt.Errorf("body must not be larger than %d bytes", maxBytes)

		// panic when decoding to non-nil pointer
		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		default:
			return err
		}
	}
	return nil
}

type Payload struct {
	multipart.File
	Size     int64
	Filename string
}

func ReadFile(rw http.ResponseWriter, r *http.Request, key, mimetype string) (*Payload, error) {
	r.Body = http.MaxBytesReader(rw, r.Body, maxUploadSize)

	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		return nil, fmt.Errorf("file exceeds max upload size: %v", err)
	}

	file, fh, err := r.FormFile(key)
	if err != nil {
		return nil, fmt.Errorf("failed to parse form data: %v", err)
	}
	defer file.Close()

	// checking mimetype
	// create buffer to store file header
	fileHeader := make([]byte, 512)
	if _, err := file.Read(fileHeader); err != nil {
		return nil, fmt.Errorf("failed to read file header: %v", err)
	}

	// set position back to start
	if _, err := file.Seek(0, 0); err != nil {
		return nil, err
	}

	if !strings.HasPrefix(http.DetectContentType(fileHeader), mimetype) {
		return nil, fmt.Errorf("incorrect mimetype, must be %s", mimetype)
	}

	return &Payload{
		File:     file,
		Size:     fh.Size,
		Filename: fh.Filename,
	}, nil
}
