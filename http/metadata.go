package http

import (
	"dusk"
	"errors"
	"net/http"

	"dusk/http/request"
	"dusk/http/response"
	"dusk/metadata"
	"dusk/util"
)

func (s *Server) ImportMetadata(rw http.ResponseWriter, r *http.Request) {
	id := HandleInt64("id", rw, r)
	if id == -1 {
		return
	}

	var input struct {
		Isbn string
	}

	err := request.Read(rw, r, &input)
	if err != nil {
		s.ErrorLog.Printf("err: %v", err)
		response.BadRequest(rw, r, err)
		return
	}

	m, err := metadata.Fetch(input.Isbn)
	if err != nil {
		s.ErrorLog.Printf("err: %v", err)
		response.InternalServerError(rw, r, err)
		return
	}

	b := m.ToBook()

	result, err := s.db.UpdateBook(id, b)
	if errors.Is(err, dusk.ErrDoesNotExist) {
		s.InfoLog.Printf("Book %d does not exist", id)
		response.NotFound(rw, r, err)
		return
	}
	if err != nil {
		s.ErrorLog.Printf("err: %v", err)
		response.InternalServerError(rw, r, err)
		return
	}

	body, err := util.ToJSON(response.Envelope{"books": result})
	if err != nil {
		s.ErrorLog.Printf("err: %v", err)
		response.InternalServerError(rw, r, err)
		return
	}

	s.InfoLog.Printf("Metadata for book %d fetched: %v", id, result)
	response.OK(rw, r, body)
}
