package ui

import (
	"dusk/http/request"
	"dusk/ui/views"
	"fmt"
	"log"
	"net/http"
)

func (s *Handler) upload(rw http.ResponseWriter, r *http.Request) {
	f, err := request.ReadFile(rw, r, "upload", "application/")
	if err != nil {
		log.Printf("failed to read file: %v", err)
		views.ImportResultsError(rw, r, err)
		return
	}

	b, err := s.fw.UploadBook(f)
	if err != nil {
		log.Printf("failed to upload cover: %v", err)
		views.ImportResultsError(rw, r, err)
		return
	}

	res, err := s.db.CreateBook(b)
	if err != nil {
		log.Printf("failed to create book: %v", err)
		views.ImportResultsError(rw, r, err)
		return
	}

	rawMessage := fmt.Sprintf("<p><a href=\"/book/%d\">%s</a> added</p>", res.ID, res.Title)
	views.ImportResultsMessage(rw, r, nil, rawMessage)
}
