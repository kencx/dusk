package ui

import (
	"dusk"
	"dusk/integrations/openlibrary"
	"dusk/ui/views"
	"dusk/validator"
	"fmt"
	"net/http"
	"regexp"
)

func (s *Handler) importPage(rw http.ResponseWriter, r *http.Request) {
	// handle tabs
	if r.URL.Query().Has("tab") {
		tab := views.TabView(r.URL.Query().Get("tab"))
		views.Tab(tab).Render(r.Context(), rw)
		return
	}

	views.Import{Tab: views.OPENLIBRARY}.Render(rw, r)
}

func (s *Handler) importOpenLibrary(rw http.ResponseWriter, r *http.Request) {
	isbn := r.FormValue("openlibrary")

	// isbn validation
	rx, err := regexp.Compile(`(\d+.*?)`)
	if err != nil {
		views.ImportResultsError(rw, r, err)
		return
	}

	if !rx.MatchString(isbn) {
		views.ImportResultsError(rw, r, views.ErrNotValidIsbn)
		return
	}

	metadata, err := openlibrary.FetchByIsbn(isbn)
	if err != nil {
		// TODO openlibrary request err
		views.ImportResultsError(rw, r, err)
		return
	}

	b := metadata.ToBook()

	v := validator.New()
	b.Validate(v)
	if !v.Valid() {
		views.ImportResultsError(rw, r, err)
		return
	}

	views.ImportResults(dusk.Books{b}, "", nil).Render(r.Context(), rw)
}

func (s *Handler) importAddResult(rw http.ResponseWriter, r *http.Request) {
	isbn := r.FormValue("result")

	// TODO We are fetching this endpoint and performing the same operations twice. It
	// would be good if we can cache the previously fetched data in importOpenLibrary on
	// the client side to send it here. This might require Alpine.js.

	metadata, err := openlibrary.FetchByIsbn(isbn)
	if err != nil {
		// TODO openlibrary request err
		views.ImportResultsError(rw, r, err)
		return
	}

	b := metadata.ToBook()

	v := validator.New()
	b.Validate(v)
	if !v.Valid() {
		views.ImportResultsError(rw, r, err)
		return
	}

	book, err := s.db.CreateBook(b)
	if err != nil {
		views.ImportResultsError(rw, r, err)
		return
	}

	rawMessage := fmt.Sprintf("<p><a href=\"/book/%d\">%s</a> added</p>", book.ID, book.Title)
	views.ImportResultsMessage(rw, r, dusk.Books{b}, rawMessage)
}
