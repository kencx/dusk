package ui

import (
	"dusk"
	"dusk/integrations/openlibrary"
	"dusk/ui/views"
	"dusk/validator"
	"net/http"
	"regexp"
)

func (s *Handler) importPage(rw http.ResponseWriter, r *http.Request) {
	// handle tabs
	if r.URL.Query().Has("tab") {
		tab := views.TabView(r.URL.Query().Get("tab"))
		views.Tab(tab, nil).Render(r.Context(), rw)
		return
	}

	views.NewImport(views.OPENLIBRARY, nil, nil).Render(rw, r)
}

func (s *Handler) importOpenLibrary(rw http.ResponseWriter, r *http.Request) {
	ol := views.OPENLIBRARY
	isbn := r.FormValue("openlibrary")

	// input validation
	rx, err := regexp.Compile(`(\d+.*?)`)
	if err != nil {
		views.NewImport(ol, nil, err).Render(rw, r)
		return
	}

	if !rx.MatchString(isbn) {
		views.NewImport(ol, nil, views.ErrNotValidIsbn).Render(rw, r)
		return
	}

	metadata, err := openlibrary.FetchByIsbn(isbn)
	if err != nil {
		// TODO openlibrary request err
		views.NewImport(ol, nil, err).Render(rw, r)
		return
	}

	b := metadata.ToBook()

	v := validator.New()
	b.Validate(v)
	if !v.Valid() {
		views.NewImport(ol, nil, err).Render(rw, r)
		return
	}

	views.NewImport(ol, dusk.Books{b}, nil).Render(rw, r)
}
