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

func (s *Handler) GetTag(rw http.ResponseWriter, r *http.Request) {
	id := request.HandleInt64("id", rw, r)
	if id == -1 {
		return
	}

	a, err := s.db.GetTag(id)
	if err == dusk.ErrDoesNotExist {
		response.NotFound(rw, r, err)
		return

	} else if err != nil {
		response.InternalServerError(rw, r, err)
		return
	}

	res, err := util.ToJSON(response.Envelope{"tags": a})
	if err != nil {
		response.InternalServerError(rw, r, err)
		return
	}

	response.OK(rw, r, res)
}

func (s *Handler) GetAllTags(rw http.ResponseWriter, r *http.Request) {

	a, err := s.db.GetAllTags(nil)
	if err == dusk.ErrNoRows {
		response.NoContent(rw, r)
		return

	} else if err != nil {
		response.InternalServerError(rw, r, err)
		return
	}

	res, err := util.ToJSON(response.Envelope{"tags": a})
	if err != nil {
		response.InternalServerError(rw, r, err)
		return
	}

	response.OK(rw, r, res)
}

func (s *Handler) AddTag(rw http.ResponseWriter, r *http.Request) {

	// marshal payload to struct
	var tag dusk.Tag
	err := request.ReadJSON(rw, r, &tag)
	if err != nil {
		response.BadRequest(rw, r, err)
		return
	}

	errMap := validator.Validate(tag)
	if len(errMap) > 0 {
		response.ValidationError(rw, r, errMap)
		return
	}

	result, err := s.db.CreateTag(&tag)
	if err != nil {
		response.InternalServerError(rw, r, err)
		return
	}

	body, err := util.ToJSON(response.Envelope{"tags": result})
	if err != nil {
		response.InternalServerError(rw, r, err)
		return
	}
	response.Created(rw, r, body)
}

func (s *Handler) UpdateTag(rw http.ResponseWriter, r *http.Request) {
	id := request.HandleInt64("id", rw, r)
	if id == -1 {
		return
	}

	// marshal payload to struct
	var tag dusk.Tag
	err := request.ReadJSON(rw, r, &tag)
	if err != nil {
		response.BadRequest(rw, r, err)
		return
	}

	// PUT should require all fields
	errMap := validator.Validate(tag)
	if len(errMap) > 0 {
		response.ValidationError(rw, r, errMap)
		return
	}

	result, err := s.db.UpdateTag(id, &tag)
	if err == dusk.ErrDoesNotExist {
		response.InternalServerError(rw, r, err)
		return
	}
	if err != nil {
		response.InternalServerError(rw, r, err)
		return
	}

	body, err := util.ToJSON(response.Envelope{"tags": result})
	if err != nil {
		response.InternalServerError(rw, r, err)
		return
	}

	response.OK(rw, r, body)
}

func (s *Handler) DeleteTag(rw http.ResponseWriter, r *http.Request) {
	id := request.HandleInt64("id", rw, r)
	if id == -1 {
		return
	}

	err := s.db.DeleteTag(id)
	if err == dusk.ErrDoesNotExist {
		response.NotFound(rw, r, err)
		return
	}

	if err != nil {
		response.InternalServerError(rw, r, err)
		return
	}

	slog.Debug("Deleted tag", slog.Int64("tag_id", id))
	response.OK(rw, r, nil)
}
