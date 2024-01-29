package api

import (
	"dusk"
	"dusk/http/request"
	"dusk/http/response"
	"dusk/util"
	"dusk/validator"
	"net/http"
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
	err := request.Read(rw, r, &book)
	if err != nil {
		response.BadRequest(rw, r, err)
		return
	}

	v := validator.New()
	book.Validate(v)
	if !v.Valid() {
		response.ValidationError(rw, r, v.Errors)
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

func (s *Handler) UpdateBook(rw http.ResponseWriter, r *http.Request) {
	id := request.HandleInt64("id", rw, r)
	if id == -1 {
		return
	}

	// marshal payload to struct
	var book dusk.Book
	err := request.Read(rw, r, &book)
	if err != nil {
		response.BadRequest(rw, r, err)
		return
	}

	// validate payload
	// PUT should require all fields
	v := validator.New()
	book.Validate(v)
	if !v.Valid() {
		response.ValidationError(rw, r, v.Errors)
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

	response.OK(rw, r, nil)
}
