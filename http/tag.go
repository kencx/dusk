package http

import (
	"dusk"
	"dusk/http/request"
	"dusk/http/response"
	"dusk/util"
	"net/http"
)

func (s *Server) GetTag(rw http.ResponseWriter, r *http.Request) {
	id := HandleInt64("id", rw, r)
	if id == -1 {
		return
	}

	a, err := s.db.GetTag(id)
	if err == dusk.ErrDoesNotExist {
		s.InfoLog.Printf("Tag %d does not exist", id)
		response.NotFound(rw, r, err)
		return

	} else if err != nil {
		s.ErrorLog.Printf("err: %v", err)
		response.InternalServerError(rw, r, err)
		return
	}

	res, err := util.ToJSON(response.Envelope{"tags": a})
	if err != nil {
		s.ErrorLog.Printf("err: %v", err)
		response.InternalServerError(rw, r, err)
		return
	}

	s.InfoLog.Printf("Tag %d retrieved: %v", id, a)
	response.OK(rw, r, res)
}

func (s *Server) GetAllTags(rw http.ResponseWriter, r *http.Request) {

	a, err := s.db.GetAllTags()
	if err == dusk.ErrNoRows {
		s.InfoLog.Println("No tags retrieved")
		response.NoContent(rw, r)
		return

	} else if err != nil {
		s.ErrorLog.Printf("err: %v", err)
		response.InternalServerError(rw, r, err)
		return
	}

	res, err := util.ToJSON(response.Envelope{"tags": a})
	if err != nil {
		s.ErrorLog.Printf("err: %v", err)
		response.InternalServerError(rw, r, err)
		return
	}

	// s.InfoLog.Printf("%d tags retrieved: %v", len(a), a)
	response.OK(rw, r, res)
}

func (s *Server) AddTag(rw http.ResponseWriter, r *http.Request) {

	// marshal payload to struct
	var tag dusk.Tag
	err := request.Read(rw, r, &tag)
	if err != nil {
		s.ErrorLog.Printf("err: %v", err)
		response.BadRequest(rw, r, err)
		return
	}

	result, err := s.db.CreateTag(&tag)
	if err != nil {
		s.ErrorLog.Printf("err: %v", err)
		response.InternalServerError(rw, r, err)
		return
	}

	body, err := util.ToJSON(response.Envelope{"tags": result})
	if err != nil {
		s.ErrorLog.Printf("err: %v", err)
		response.InternalServerError(rw, r, err)
		return
	}
	s.InfoLog.Printf("New tag created: %v", result)
	response.Created(rw, r, body)
}

func (s *Server) UpdateTag(rw http.ResponseWriter, r *http.Request) {
	id := HandleInt64("id", rw, r)
	if id == -1 {
		return
	}

	// marshal payload to struct
	var tag dusk.Tag
	err := request.Read(rw, r, &tag)
	if err != nil {
		s.ErrorLog.Printf("err: %v", err)
		response.BadRequest(rw, r, err)
		return
	}


	result, err := s.db.UpdateTag(id, &tag)
	if err == dusk.ErrDoesNotExist {
		s.InfoLog.Printf("Tag %d does not exist", id)
		response.InternalServerError(rw, r, err)
		return
	}
	if err != nil {
		s.ErrorLog.Printf("err: %v", err)
		response.InternalServerError(rw, r, err)
		return
	}

	body, err := util.ToJSON(response.Envelope{"tags": result})
	if err != nil {
		s.ErrorLog.Printf("err: %v", err)
		response.InternalServerError(rw, r, err)
		return
	}

	s.InfoLog.Printf("Tag %d updated: %v", id, result)
	response.OK(rw, r, body)
}

func (s *Server) DeleteTag(rw http.ResponseWriter, r *http.Request) {
	id := HandleInt64("id", rw, r)
	if id == -1 {
		return
	}

	err := s.db.DeleteTag(id)
	if err == dusk.ErrDoesNotExist {
		s.InfoLog.Printf("Tag %d does not exist", id)
		response.NotFound(rw, r, err)
		return
	}

	if err != nil {
		s.ErrorLog.Printf("err: %v", err)
		response.InternalServerError(rw, r, err)
		return
	}

	s.InfoLog.Printf("Tag %d deleted", id)
	response.OK(rw, r, nil)
}
