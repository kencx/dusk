package http

import (
	"dusk"
	"dusk/http/request"
	"dusk/http/response"
	"dusk/util"
	"dusk/validator"
	"net/http"
)

func (s *Server) GetAuthor(rw http.ResponseWriter, r *http.Request) {
	id := HandleInt64("id", rw, r)
	if id == -1 {
		return
	}

	a, err := s.db.GetAuthor(id)
	if err == dusk.ErrDoesNotExist {
		s.InfoLog.Printf("Author %d does not exist", id)
		response.NotFound(rw, r, err)
		return

	} else if err != nil {
		s.ErrorLog.Printf("err: %v", err)
		response.InternalServerError(rw, r, err)
		return
	}

	res, err := util.ToJSON(response.Envelope{"authors": a})
	if err != nil {
		s.ErrorLog.Printf("err: %v", err)
		response.InternalServerError(rw, r, err)
		return
	}

	s.InfoLog.Printf("Author %d retrieved: %v", id, a)
	response.OK(rw, r, res)
}

func (s *Server) GetAllAuthors(rw http.ResponseWriter, r *http.Request) {

	a, err := s.db.GetAllAuthors()
	if err == dusk.ErrNoRows {
		s.InfoLog.Println("No authors retrieved")
		response.NoContent(rw, r)
		return

	} else if err != nil {
		s.ErrorLog.Printf("err: %v", err)
		response.InternalServerError(rw, r, err)
		return
	}

	res, err := util.ToJSON(response.Envelope{"authors": a})
	if err != nil {
		s.ErrorLog.Printf("err: %v", err)
		response.InternalServerError(rw, r, err)
		return
	}

	// s.InfoLog.Printf("%d authors retrieved: %v", len(a), a)
	response.OK(rw, r, res)
}

func (s *Server) AddAuthor(rw http.ResponseWriter, r *http.Request) {

	// marshal payload to struct
	var author dusk.Author
	err := request.Read(rw, r, &author)
	if err != nil {
		s.ErrorLog.Printf("err: %v", err)
		response.BadRequest(rw, r, err)
		return
	}

	// validate payload
	v := validator.New()
	author.Validate(v)
	if !v.Valid() {
		response.ValidationError(rw, r, v.Errors)
		return
	}

	result, err := s.db.CreateAuthor(&author)
	if err != nil {
		s.ErrorLog.Printf("err: %v", err)
		response.InternalServerError(rw, r, err)
		return
	}

	body, err := util.ToJSON(response.Envelope{"authors": result})
	if err != nil {
		s.ErrorLog.Printf("err: %v", err)
		response.InternalServerError(rw, r, err)
		return
	}
	s.InfoLog.Printf("New author created: %v", result)
	response.Created(rw, r, body)
}

func (s *Server) UpdateAuthor(rw http.ResponseWriter, r *http.Request) {
	id := HandleInt64("id", rw, r)
	if id == -1 {
		return
	}

	// marshal payload to struct
	var author dusk.Author
	err := request.Read(rw, r, &author)
	if err != nil {
		s.ErrorLog.Printf("err: %v", err)
		response.BadRequest(rw, r, err)
		return
	}

	// validate payload
	// PUT should require all fields
	v := validator.New()
	author.Validate(v)
	if !v.Valid() {
		response.ValidationError(rw, r, v.Errors)
		return
	}

	result, err := s.db.UpdateAuthor(id, &author)
	if err == dusk.ErrDoesNotExist {
		s.InfoLog.Printf("Author %d does not exist", id)
		response.InternalServerError(rw, r, err)
		return
	}
	if err != nil {
		s.ErrorLog.Printf("err: %v", err)
		response.InternalServerError(rw, r, err)
		return
	}

	body, err := util.ToJSON(response.Envelope{"authors": result})
	if err != nil {
		s.ErrorLog.Printf("err: %v", err)
		response.InternalServerError(rw, r, err)
		return
	}

	s.InfoLog.Printf("Author %d updated: %v", id, result)
	response.OK(rw, r, body)
}

func (s *Server) DeleteAuthor(rw http.ResponseWriter, r *http.Request) {
	id := HandleInt64("id", rw, r)
	if id == -1 {
		return
	}

	err := s.db.DeleteAuthor(id)
	if err == dusk.ErrDoesNotExist {
		s.InfoLog.Printf("Author %d does not exist", id)
		response.NotFound(rw, r, err)
		return
	}

	if err != nil {
		s.ErrorLog.Printf("err: %v", err)
		response.InternalServerError(rw, r, err)
		return
	}

	s.InfoLog.Printf("Author %d deleted", id)
	response.OK(rw, r, nil)
}
