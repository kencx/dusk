package file

import (
	"fmt"
	"mime"
	"mime/multipart"
	"path/filepath"
)

var (
	mimeExtMap = map[string]string{
		"application/epub+zip":           epubExt,
		"application/vnd.amazon.ebook":   azwExt,
		"application/x-mobipocket-ebook": mobiExt,
		"application/pdf":                pdfExt,
		"image/jpeg":                     jpegExt,
		"image/png":                      pngExt,
		"image/vnd.djvu":                 djvuExt,
	}
	defaultMime = "application/octet-stream"
)

type Payload struct {
	multipart.File
	Size      int64
	Filename  string
	MimeType  string
	Extension string
}

func NewPayload(file multipart.File, fh *multipart.FileHeader, mimetype string) (*Payload, error) {
	p := &Payload{
		File:     file,
		Size:     fh.Size,
		Filename: fh.Filename,
		MimeType: mimetype,
	}

	ext, err := extension(fh.Filename, mimetype)
	if err != nil {
		return nil, err
	}

	p.Extension = ext
	return p, nil
}

// set extension by given filename or mimetype
func extension(filename, mimetype string) (string, error) {
	// from filename
	ext := filepath.Ext(filename)
	if ext != "" {
		return ext, nil
	}

	// from known types
	exts, err := mime.ExtensionsByType(mimetype)
	if err != nil {
		return "", fmt.Errorf("failed to get extension from mimetype: %w", err)
	}
	if len(exts) > 0 {
		return exts[0], nil
	}

	// from custom types
	if e, ok := mimeExtMap[mimetype]; ok {
		return e, nil
	}

	return defaultMime, nil
}
