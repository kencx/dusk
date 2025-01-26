package ui

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/araddon/dateparse"
	"github.com/kencx/dusk"
	"github.com/kencx/dusk/http/request"
	"github.com/kencx/dusk/http/response"
	"github.com/kencx/dusk/null"
	"github.com/kencx/dusk/page"
	"github.com/kencx/dusk/ui/partials"
	"github.com/kencx/dusk/ui/views"
	"github.com/kencx/dusk/validator"
)

// Perform FTS on library (with pagination)
func (s *Handler) bookSearch(rw http.ResponseWriter, r *http.Request) {
	filters := initBookFilters(r)
	if errMap := validator.Validate(filters.Base); errMap != nil {
		slog.Error("[ui] failed to validate query params", slog.Any("err", errMap.Error()))
		partials.BookSearchResults(page.Page[dusk.Book]{}, errors.New("validate error")).Render(r.Context(), rw)
		return
	}

	p, err := s.db.GetAllBooks(filters)
	if err != nil {
		slog.Error("failed to query books", slog.Any("err", err))
		partials.BookSearchResults(page.Page[dusk.Book]{}, err).Render(r.Context(), rw)
		return
	}
	partials.BookSearchResults(*p, nil).Render(r.Context(), rw)
}

// Render details of book
func (s *Handler) bookPage(rw http.ResponseWriter, r *http.Request) {
	id := request.FetchIdFromSlug(rw, r)
	if id == -1 {
		return
	}

	book, err := s.db.GetBook(id)
	if err != nil {
		slog.Error("[ui] failed to find book", slog.Int64("id", id), slog.Any("err", err))
		views.NewBook(s.base, nil, nil, nil, err).Render(rw, r)
		return
	}

	if r.URL.Query().Has("delete") {
		response.AddHxTriggerAfterSwap(rw, `{"openModal": ""}`)
		views.DeleteBookModal(book).Render(r.Context(), rw)
		return
	}

	authors, err := s.db.GetAuthorsFromBook(id)
	if err != nil {
		slog.Error("[ui] failed to fetch authors of book", slog.Int64("id", id), slog.Any("err", err))
		views.NewBook(s.base, nil, nil, nil, err).Render(rw, r)
		return
	}

	tags, err := s.db.GetTagsFromBook(id)
	if err != nil {
		slog.Error("[ui] failed to fetch tags of book", slog.Int64("id", id), slog.Any("err", err))
		views.NewBook(s.base, nil, nil, nil, err).Render(rw, r)
		return
	}

	views.NewBook(s.base, book, authors, tags, nil).Render(rw, r)
}

func (s *Handler) editBookForm(rw http.ResponseWriter, r *http.Request) {
	id := request.FetchIdFromSlug(rw, r)
	if id == -1 {
		return
	}

	book, err := s.db.GetBook(id)
	if err != nil {
		slog.Error("[ui] failed to render edit book form", slog.Int64("id", id), slog.Any("err", err))
		views.NewBookForm(s.base, nil, err).Render(rw, r)
		return
	}

	views.NewBookForm(s.base, book, nil).Render(rw, r)
}

func (s *Handler) updateBook(rw http.ResponseWriter, r *http.Request) {
	id := request.FetchIdFromSlug(rw, r)
	if id == -1 {
		return
	}

	book, err := s.db.GetBook(id)
	if err != nil {
		slog.Error("[ui] failed to get book", slog.Int64("id", id), slog.Any("err", err))
		return
	}

	if err = r.ParseForm(); err != nil {
		// TODO validation error
		slog.Error("[ui] failed to parse form", slog.Int64("id", id), slog.Any("err", err))
		return
	}

	// TODO validation error
	// only include values that have changed
	book = parseBookForm(r, book)
	if errMap := validator.Validate(book); errMap != nil {
		slog.Error("[ui] failed to validate book", slog.Int64("id", id), slog.String("err", errMap.Error()))
		views.NewBookForm(s.base, nil, errMap).Render(rw, r)
		return
	}

	new_book, err := s.db.UpdateBook(id, book)
	if err != nil {
		slog.Error("[ui] failed to update book", slog.Int64("id", id), slog.Any("err", err))
		views.NewBook(s.base, nil, nil, nil, err).Render(rw, r)
		return
	}
	// redirect to book page
	response.HxRedirect(rw, r, "/b/"+new_book.Slugify())
}

func (s *Handler) updateBookStatus(rw http.ResponseWriter, r *http.Request) {
	var readStatus dusk.ReadStatus
	switch r.FormValue("read-status") {
	case "unread":
		readStatus = dusk.Unread
	case "reading":
		readStatus = dusk.Reading
	case "read":
		readStatus = dusk.Read
	default:
		readStatus = dusk.Unread
	}

	id := request.FetchIdFromSlug(rw, r)
	if id == -1 {
		return
	}

	book, err := s.db.GetBook(id)
	if err != nil {
		slog.Error("[ui] failed to get book", slog.Int64("id", id), slog.Any("err", err))
		return
	}

	book.Status = readStatus
	if readStatus == dusk.Read {
		book.DateCompleted = null.TimeFrom(time.Now())
	}

	new_book, err := s.db.UpdateBook(id, book)
	if err != nil {
		slog.Error("[ui] failed to update book", slog.Int64("id", id), slog.Any("err", err))
		views.NewBook(s.base, nil, nil, nil, err).Render(rw, r)
		return
	}

	views.BookStatus(new_book).Render(r.Context(), rw)
}

func (s *Handler) deleteBook(rw http.ResponseWriter, r *http.Request) {
	id := request.FetchIdFromSlug(rw, r)
	if id == -1 {
		return
	}

	book, err := s.db.GetBook(id)
	if err != nil {
		slog.Error("[ui] failed to get book", slog.Int64("id", id), slog.Any("err", err))
		views.NewBook(s.base, nil, nil, nil, err).Render(rw, r)
		return
	}

	err = s.db.DeleteBook(id)
	if err != nil {
		slog.Error("[ui] failed to delete book", slog.Int64("id", id), slog.Any("err", err))
		views.NewBook(s.base, nil, nil, nil, err).Render(rw, r)
		return
	}

	// TODO: allow archive files instead of permenently delete
	err = s.fs.DeleteBook(book)
	if err != nil {
		slog.Warn("[ui] failed to delete book from filesystem", slog.Int64("id", id), slog.Any("err", err))
	}

	// redirect to index page
	response.HxRedirect(rw, r, "/")
}

func parseBookForm(r *http.Request, b *dusk.Book) *dusk.Book {
	if request.HasValue(r.Form, "title") {
		b.Title = r.FormValue("title")
	}

	if request.HasOptionalValue(r.Form, "subtitle") {
		b.Subtitle = null.StringFrom(r.FormValue("subtitle"))
	}

	if request.HasOptionalValue(r.Form, "author") {
		var authors []string

		for _, a := range strings.Split(r.FormValue("author"), ";") {
			author := strings.TrimSpace(a)
			if author != "" {
				authors = append(authors, author)
			}
		}

		b.Author = authors
	}

	if request.HasOptionalValue(r.Form, "tags") {
		var tags []string

		for _, t := range strings.Split(r.FormValue("tags"), ",") {
			tag := strings.TrimSpace(t)
			if tag != "" {
				tags = append(tags, tag)
			}
		}

		b.Tag = tags
	}

	if request.HasOptionalValue(r.Form, "series") {
		b.Series = null.StringFrom(r.FormValue("series"))
	}

	if request.HasOptionalValue(r.Form, "numOfPages") {
		pages, _ := strconv.Atoi(r.FormValue("numOfPages"))
		b.NumOfPages = pages
	}

	if request.HasOptionalValue(r.Form, "rating") {
		rating, _ := strconv.Atoi(r.FormValue("rating"))
		b.Rating = rating
	}

	if request.HasOptionalValue(r.Form, "publisher") {
		b.Publisher = null.StringFrom(r.FormValue("publisher"))
	}

	if request.HasOptionalValue(r.Form, "datePublished") {
		dp, _ := dateparse.ParseAny(r.FormValue("datePublished"))
		b.DatePublished = null.TimeFrom(dp)
	}

	if request.HasValue(r.Form, "read-status") {
		var status dusk.ReadStatus
		switch strings.ToLower(r.FormValue("read-status")) {
		case "unread":
			status = dusk.Unread
		case "reading":
			status = dusk.Reading
		case "read":
			status = dusk.Read
		default:
			status = dusk.Unread
		}
		b.Status = status
	}

	if request.HasOptionalValue(r.Form, "dateAdded") {
		dp, _ := dateparse.ParseAny(r.FormValue("dateAdded"))
		b.DateAdded = null.TimeFrom(dp)
	}

	if request.HasOptionalValue(r.Form, "dateStarted") {
		dp, _ := dateparse.ParseAny(r.FormValue("dateStarted"))
		b.DateStarted = null.TimeFrom(dp)
	}

	if request.HasOptionalValue(r.Form, "dateCompleted") {
		dp, _ := dateparse.ParseAny(r.FormValue("dateCompleted"))
		b.DateCompleted = null.TimeFrom(dp)
	}

	if request.HasOptionalValue(r.Form, "description") {
		b.Description = null.StringFrom(r.FormValue("description"))
	}

	if request.HasOptionalValue(r.Form, "notes") {
		b.Notes = null.StringFrom(r.FormValue("notes"))
	}
	return b
}
