package ui

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	ol "github.com/kencx/dusk/integrations/openlibrary"
	"github.com/kencx/dusk/ui/partials"
	"github.com/kencx/dusk/ui/views"
	"github.com/kencx/dusk/util"
	"github.com/kencx/dusk/validator"
)

func (s *Handler) importPage(rw http.ResponseWriter, r *http.Request) {
	// handle htmx tabs
	if r.URL.Query().Has("tab") {
		tab := partials.TabName(r.URL.Query().Get("tab"))
		views.ImportTabs.Select(tab).Render(r.Context(), rw)
		return
	}

	views.Search{DefaultTab: views.OPENLIBRARY}.Render(rw, r)
}

func (s *Handler) importOpenLibrary(rw http.ResponseWriter, r *http.Request) {
	isbn := r.FormValue("openlibrary")

	// isbn validation
	rx, err := regexp.Compile(`(\d+.*?)`)
	if err != nil {
		views.ImportResultsError(rw, r, err)
		return
	}

	if !validator.Matches(isbn, rx) {
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

	errMap := validator.Validate(b)
	if len(errMap) > 0 {
		views.ImportResultsError(rw, r, errors.New("TODO"))
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

	errMap := validator.Validate(b)
	if len(errMap) > 0 {
		views.ImportResultsError(rw, r, errors.New("TODO"))
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
