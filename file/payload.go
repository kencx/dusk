package file

import (
	"mime/multipart"
	"path/filepath"
)

var mime2Ext = map[string]string{
	"application/epub+zip": ".epub",
	"application/pdf":      ".pdf",
	"image/jpeg":           ".jpeg",
	"image/png":            ".png",
}

type Payload struct {
	multipart.File
	Size     int64
	Filename string
	MimeType string
}

// set extension by given filename or mimetype
func (p *Payload) Extension() string {
	defaultExt := ".epub"
	ext := filepath.Ext(p.Filename)
	if ext != "" {
		return ext
	}

	if e, ok := mime2Ext[p.MimeType]; ok {
		return e
	}
	return defaultExt
}
