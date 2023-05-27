package http

import (
	"dusk"
	"dusk/http/request"
	"dusk/http/response"
	"dusk/util"
	"dusk/validator"
	"net/http"
)

func (s *Server) GetBook(rw http.ResponseWriter, r *http.Request) {
	id := HandleInt64("id", rw, r)
	if id == -1 {
		return
	}

	b, err := s.db.GetBook(id)
	if err == dusk.ErrDoesNotExist {
		s.InfoLog.Printf("Book %d does not exist", id)
		response.NotFound(rw, r, err)
		return

	} else if err != nil {
		s.ErrorLog.Printf("err: %v", err)
		response.InternalServerError(rw, r, err)
		return
	}

	res, err := util.ToJSON(response.Envelope{"books": b})
	if err != nil {
		s.ErrorLog.Printf("err: %v", err)
		response.InternalServerError(rw, r, err)
		return
	}

	s.InfoLog.Printf("Book %d retrieved: %v", id, b)
	response.OK(rw, r, res)
}

func (s *Server) GetAllBooks(rw http.ResponseWriter, r *http.Request) {
	b, err := s.db.GetAllBooks()
	if err == dusk.ErrNoRows {
		s.InfoLog.Println("No books retrieved")
		response.NoContent(rw, r)
		return

	} else if err != nil {
		s.ErrorLog.Printf("err: %v", err)
		response.InternalServerError(rw, r, err)
		return
	}

	res, err := util.ToJSON(response.Envelope{"books": b})
	if err != nil {
		s.ErrorLog.Printf("err: %v", err)
		response.InternalServerError(rw, r, err)
		return
	}

	s.InfoLog.Printf("%d books retrieved: %v", len(b), b)
	response.OK(rw, r, res)
}

func (s *Server) AddBook(rw http.ResponseWriter, r *http.Request) {

	// marshal payload to struct
	var book dusk.Book
	err := request.Read(rw, r, &book)
	if err != nil {
		s.ErrorLog.Printf("err: %v", err)
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
		s.ErrorLog.Printf("err: %v", err)
		response.BadRequest(rw, r, err)
		return
	}

	body, err := util.ToJSON(response.Envelope{"books": result})
	if err != nil {
		s.ErrorLog.Println(err)
		response.InternalServerError(rw, r, err)
		return
	}
	s.InfoLog.Printf("New book created: %v", result)
	response.Created(rw, r, body)
}

func (s *Server) UpdateBook(rw http.ResponseWriter, r *http.Request) {
	id := HandleInt64("id", rw, r)
	if id == -1 {
		return
	}

	// marshal payload to struct
	var book dusk.Book
	err := request.Read(rw, r, &book)
	if err != nil {
		s.ErrorLog.Printf("err: %v", err)
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

	s.InfoLog.Printf("Book %d updated: %v", id, result)
	response.OK(rw, r, body)
}

func (s *Server) DeleteBook(rw http.ResponseWriter, r *http.Request) {
	id := HandleInt64("id", rw, r)
	if id == -1 {
		return
	}

	err := s.db.DeleteBook(id)
	if err == dusk.ErrDoesNotExist {
		s.InfoLog.Printf("Book %d does not exist", id)
		response.NotFound(rw, r, err)
		return
	}

	if err != nil {
		s.ErrorLog.Printf("err: %v", err)
		response.InternalServerError(rw, r, err)
		return
	}

	s.InfoLog.Printf("Book %d deleted", id)
	response.OK(rw, r, nil)
}
