package ui

import (
	"dusk/metadata"
	"dusk/ui/views"
	"dusk/validator"
	"net/http"
)

const (
	OPENLIBRARY = "openlibrary"
	GOODREADS   = "goodreads"
	CALIBRE     = "calibre"
)

func (s *Handler) importPage(rw http.ResponseWriter, r *http.Request) {
	m := views.NewImport(OPENLIBRARY)

	// handle tabs
	if r.URL.Query().Has("tab") {
		tab := r.URL.Query().Get("tab")
		views.Tab(tab).Render(r.Context(), rw)
		return
	}

	m.Render(rw, r)
}

func (s *Handler) importOpenLibrary(rw http.ResponseWriter, r *http.Request) {
	// TODO add clearer error messages

	m := views.NewImport(OPENLIBRARY)

	isbn := r.FormValue("openlibrary")
	metadata, err := metadata.Fetch(isbn)
	if err != nil {
		m.Render(rw, r)
		return
	}

	b := metadata.ToBook()

	v := validator.New()
	b.Validate(v)
	if !v.Valid() {
		m.Render(rw, r)
		return
	}

	_, err = s.db.CreateBook(b)
	if err != nil {
		m.Render(rw, r)
		return
	}

	m.Render(rw, r)
}
