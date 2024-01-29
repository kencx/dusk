package ui

import (
	"dusk/metadata"
	"dusk/ui/pages"
	"dusk/validator"
	"errors"
	"net/http"
)

const (
	OPENLIBRARY = "openlibrary"
	GOODREADS   = "goodreads"
	CALIBRE     = "calibre"
)

func (s *Handler) importView(rw http.ResponseWriter, r *http.Request) {
	m := pages.NewImportViewModel(OPENLIBRARY, nil)

	// handle tabs
	if r.URL.Query().Has("tab") {
		tab := r.URL.Query().Get("tab")
		pages.Tab(tab).Render(r.Context(), rw)
		return
	}

	m.Render(rw, r)
}

func (s *Handler) importOpenLibrary(rw http.ResponseWriter, r *http.Request) {
	// TODO add clearer error messages

	m := pages.NewImportViewModel(OPENLIBRARY, nil)

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
