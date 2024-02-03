package ui

import (
	"dusk/metadata"
	"dusk/ui/views"
	"dusk/validator"
	"errors"
	"net/http"
)

const (
	OPENLIBRARY = "openlibrary"
	GOODREADS   = "goodreads"
	CALIBRE     = "calibre"
)

func (s *Handler) importPage(rw http.ResponseWriter, r *http.Request) {
	m := views.NewImportViewModel(OPENLIBRARY, nil)

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

	m := views.NewImportViewModel(OPENLIBRARY, nil)

	isbn := r.FormValue("openlibrary")
	metadata, err := metadata.Fetch(isbn)
	if err != nil {
		m.RenderError(rw, r, err)
		return
	}

	b := metadata.ToBook()

	v := validator.New()
	b.Validate(v)
	if !v.Valid() {
		m.RenderError(rw, r, errors.New("Something went wrong, please try again"))
		return
	}

	_, err = s.db.CreateBook(b)
	if err != nil {
		m.RenderError(rw, r, err)
		return
	}

	m.Render(rw, r)
}
