package api

import (
	"log/slog"
	"net/http"

	"github.com/kencx/dusk"
	"github.com/kencx/dusk/http/request"
	"github.com/kencx/dusk/http/response"
	"github.com/kencx/dusk/util"
	"github.com/kencx/dusk/validator"
)

func (s *Handler) GetBook(rw http.ResponseWriter, r *http.Request) {
	id := request.HandleInt64("id", rw, r)
	if id == -1 {
		return
	}

	b, err := s.db.GetBook(id)
	if err == dusk.ErrDoesNotExist {
		response.NotFound(rw, r, err)
		return

	} else if err != nil {
		response.InternalServerError(rw, r, err)
		return
	}

	res, err := util.ToJSON(response.Envelope{"books": b})
	if err != nil {
		response.InternalServerError(rw, r, err)
		return
	}

	response.OK(rw, r, res)
}

func (s *Handler) GetAllBooks(rw http.ResponseWriter, r *http.Request) {
	b, err := s.db.GetAllBooks()
	if err == dusk.ErrNoRows {
		response.NoContent(rw, r)
		return

	} else if err != nil {
		response.InternalServerError(rw, r, err)
		return
	}

	res, err := util.ToJSON(response.Envelope{"books": b})
	if err != nil {
		response.InternalServerError(rw, r, err)
		return
	}

	response.OK(rw, r, res)
}

func (s *Handler) AddBook(rw http.ResponseWriter, r *http.Request) {

	// marshal payload to struct
	var book dusk.Book
	err := request.ReadJSON(rw, r, &book)
	if err != nil {
		response.BadRequest(rw, r, err)
		return
	}

	errMap := validator.Validate(book)
	if len(errMap) > 0 {
		response.ValidationError(rw, r, errMap)
		return
	}

	result, err := s.db.CreateBook(&book)
	if err != nil {
		response.BadRequest(rw, r, err)
		return
	}

	body, err := util.ToJSON(response.Envelope{"books": result})
	if err != nil {
		response.InternalServerError(rw, r, err)
		return
	}
	response.Created(rw, r, body)
}

func (s *Handler) AddBookCover(rw http.ResponseWriter, r *http.Request) {
	id := request.HandleInt64("id", rw, r)
	if id == -1 {
		return
	}

	b, err := s.db.GetBook(id)
	if err == dusk.ErrDoesNotExist {
		response.NotFound(rw, r, err)
		return
	}
	if err != nil {
		response.InternalServerError(rw, r, err)
		return
	}

	file, err := request.ReadFile(rw, r, "cover", "image/")
	if err != nil {
		response.BadRequest(rw, r, err)
		return
	}

	if err := s.fs.UploadBookCover(file, b); err != nil {
		slog.Error("[API] Failed to upload file", slog.Any("err", err))
		response.InternalServerError(rw, r, err)
		return
	}

	result, err := s.db.UpdateBook(id, b)
	if err != nil {
		// TODO delete uploaded file on err
		slog.Error("[API] Failed to update book", slog.Any("err", err))
		response.InternalServerError(rw, r, err)
		return
	}

	body, err := util.ToJSON(response.Envelope{"books": result})
	if err != nil {
		// TODO delete uploaded file on err
		response.InternalServerError(rw, r, err)
		return
	}

	response.OK(rw, r, body)
}

// TODO should adding format change metadata?
func (s *Handler) AddBookFormat(rw http.ResponseWriter, r *http.Request) {
	id := request.HandleInt64("id", rw, r)
	if id == -1 {
		return
	}

	b, err := s.db.GetBook(id)
	if err == dusk.ErrDoesNotExist {
		response.NotFound(rw, r, err)
		return
	}
	if err != nil {
		response.InternalServerError(rw, r, err)
		return
	}

	file, err := request.ReadFile(rw, r, "format", "application/")
	if err != nil {
		response.BadRequest(rw, r, err)
		return
	}

	if err := s.fs.UploadBookFormat(file, b); err != nil {
		slog.Error("[API] Failed to upload file", slog.Any("err", err))
		response.InternalServerError(rw, r, err)
		return
	}

	result, err := s.db.UpdateBook(id, b)
	if err != nil {
		// TODO delete uploaded file on err
		response.InternalServerError(rw, r, err)
		return
	}

	body, err := util.ToJSON(response.Envelope{"books": result})
	if err != nil {
		// TODO delete uploaded file on err
		response.InternalServerError(rw, r, err)
		return
	}

	response.OK(rw, r, body)
}

func (s *Handler) UpdateBook(rw http.ResponseWriter, r *http.Request) {
	id := request.HandleInt64("id", rw, r)
	if id == -1 {
		return
	}

	// marshal payload to struct
	var book dusk.Book
	err := request.ReadJSON(rw, r, &book)
	if err != nil {
		response.BadRequest(rw, r, err)
		return
	}

	// PUT should require all fields
	errMap := validator.Validate(book)
	if len(errMap) > 0 {
		response.ValidationError(rw, r, errMap)
		return
	}

	result, err := s.db.UpdateBook(id, &book)
	if err == dusk.ErrDoesNotExist {
		response.NotFound(rw, r, err)
		return
	}
	if err != nil {
		response.InternalServerError(rw, r, err)
		return
	}

	body, err := util.ToJSON(response.Envelope{"books": result})
	if err != nil {
		response.InternalServerError(rw, r, err)
		return
	}

	response.OK(rw, r, body)
}

func (s *Handler) DeleteBook(rw http.ResponseWriter, r *http.Request) {
	id := request.HandleInt64("id", rw, r)
	if id == -1 {
		return
	}

	err := s.db.DeleteBook(id)
	if err == dusk.ErrDoesNotExist {
		response.NotFound(rw, r, err)
		return
	}

	if err != nil {
		response.InternalServerError(rw, r, err)
		return
	}

	slog.Debug("Deleted book", slog.Int64("book_id", id))
	response.OK(rw, r, nil)
}
