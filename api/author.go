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

func (s *Handler) GetAuthor(rw http.ResponseWriter, r *http.Request) {
	id := request.HandleInt64("id", rw, r)
	if id == -1 {
		return
	}

	a, err := s.db.GetAuthor(id)
	if err == dusk.ErrDoesNotExist {
		response.NotFound(rw, r, err)
		return

	} else if err != nil {
		response.InternalServerError(rw, r, err)
		return
	}

	res, err := util.ToJSON(response.Envelope{"authors": a})
	if err != nil {
		response.InternalServerError(rw, r, err)
		return
	}

	response.OK(rw, r, res)
}

func (s *Handler) GetAllAuthors(rw http.ResponseWriter, r *http.Request) {

	a, err := s.db.GetAllAuthors()
	if err == dusk.ErrNoRows {
		response.NoContent(rw, r)
		return

	} else if err != nil {
		response.InternalServerError(rw, r, err)
		return
	}

	res, err := util.ToJSON(response.Envelope{"authors": a})
	if err != nil {
		response.InternalServerError(rw, r, err)
		return
	}

	response.OK(rw, r, res)
}

func (s *Handler) AddAuthor(rw http.ResponseWriter, r *http.Request) {

	// marshal payload to struct
	var author dusk.Author
	err := request.ReadJSON(rw, r, &author)
	if err != nil {
		response.BadRequest(rw, r, err)
		return
	}

	errMap := validator.Validate(author)
	if len(errMap) > 0 {
		response.ValidationError(rw, r, errMap)
		return
	}

	result, err := s.db.CreateAuthor(&author)
	if err != nil {
		response.InternalServerError(rw, r, err)
		return
	}

	body, err := util.ToJSON(response.Envelope{"authors": result})
	if err != nil {
		response.InternalServerError(rw, r, err)
		return
	}
	response.Created(rw, r, body)
}

func (s *Handler) UpdateAuthor(rw http.ResponseWriter, r *http.Request) {
	id := request.HandleInt64("id", rw, r)
	if id == -1 {
		return
	}

	// marshal payload to struct
	var author dusk.Author
	err := request.ReadJSON(rw, r, &author)
	if err != nil {
		response.BadRequest(rw, r, err)
		return
	}

	errMap := validator.Validate(author)
	if len(errMap) > 0 {
		response.ValidationError(rw, r, errMap)
		return
	}

	result, err := s.db.UpdateAuthor(id, &author)
	if err == dusk.ErrDoesNotExist {
		response.InternalServerError(rw, r, err)
		return
	}
	if err != nil {
		response.InternalServerError(rw, r, err)
		return
	}

	body, err := util.ToJSON(response.Envelope{"authors": result})
	if err != nil {
		response.InternalServerError(rw, r, err)
		return
	}

	response.OK(rw, r, body)
}

func (s *Handler) DeleteAuthor(rw http.ResponseWriter, r *http.Request) {
	id := request.HandleInt64("id", rw, r)
	if id == -1 {
		return
	}

	err := s.db.DeleteAuthor(id)
	if err == dusk.ErrDoesNotExist {
		response.NotFound(rw, r, err)
		return
	}

	if err != nil {
		response.InternalServerError(rw, r, err)
		return
	}

	slog.Debug("Deleted author", slog.Int64("author_id", id))
	response.OK(rw, r, nil)
}
