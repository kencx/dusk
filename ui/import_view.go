package ui

import (
	"dusk/metadata"
	"dusk/ui/pages"
	"dusk/validator"
	"net/http"
)

func (s *Handler) importView(rw http.ResponseWriter, r *http.Request) {
	// handle tabs
	if r.URL.Query().Has("tab") {
		tab := r.URL.Query().Get("tab")
		pages.Tab(tab).Render(r.Context(), rw)
		return
	}

	pages.Import("openlibrary", "").Render(r.Context(), rw)
}

func (s *Handler) importOpenLibrary(rw http.ResponseWriter, r *http.Request) {
	isbn := r.FormValue("openlibrary")
	m, err := metadata.Fetch(isbn)
	if err != nil {
		pages.Import("openlibrary", "Something went wrong").Render(r.Context(), rw)
		return
	}

	b := m.ToBook()

	v := validator.New()
	b.Validate(v)
	if !v.Valid() {
		pages.Import("openlibrary", "Something went wrong").Render(r.Context(), rw)
		return
	}

	_, err = s.db.CreateBook(b)
	if err != nil {
		pages.Import("openlibrary", "Something went wrong").Render(r.Context(), rw)
		return
	}

	pages.Import("openlibrary", "Book imported").Render(r.Context(), rw)
}
