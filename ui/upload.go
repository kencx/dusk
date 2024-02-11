package ui

import (
	"dusk/epub"
	"dusk/http/request"
	"dusk/ui/views"
	"dusk/validator"
	"errors"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
)

func (s *Handler) upload(rw http.ResponseWriter, r *http.Request) {
	f, err := request.ReadFile(rw, r, "upload", "application/")
	if err != nil {
		log.Printf("failed to read file: %v", err)
		views.ImportResultsError(rw, r, err)
		return
	}

	if filepath.Ext(f.Filename) == ".epub" {
		s.uploadEpub(rw, r, f)
	}

}

func (s *Handler) uploadEpub(rw http.ResponseWriter, r *http.Request, f *request.Payload) {
	ep, err := epub.NewFromReader(f.File, f.Size)
	if err != nil {
		log.Printf("failed to create epub: %v", err)
		views.ImportResultsError(rw, r, err)
		return
	}

	b := ep.ToBook()
	errMap := validator.Validate(b)
	if errMap != nil {
		views.ImportResultsError(rw, r, errors.New("TODO"))
		return
	}

	path, err := s.fw.UploadFile(f.File, b.Title, fmt.Sprintf("%s.%s", b.Title, "epub"))
	if err != nil {
		log.Printf("failed to upload file: %v", err)
		return
	}

	cf, err := ep.Open(ep.CoverFile)
	if err != nil {
		log.Printf("failed to open cover: %v", err)
		return
	}

	coverPath, err := s.fw.UploadFile(cf, b.Title, "cover.jpg")
	if err != nil {
		log.Printf("failed to upload cover: %v", err)
		views.ImportResultsError(rw, r, err)
		return
	}
	b.Cover = s.fw.GetRelativePath(coverPath)
	b.Formats = append(b.Formats, s.fw.GetRelativePath(path))

	res, err := s.db.CreateBook(b)
	if err != nil {
		log.Printf("failed to create book: %v", err)
		views.ImportResultsError(rw, r, err)
		return
	}

	rawMessage := fmt.Sprintf("<p><a href=\"/book/%d\">%s</a> added</p>", res.ID, res.Title)
	views.ImportResultsMessage(rw, r, nil, rawMessage)
}
