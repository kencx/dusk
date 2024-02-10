package api

import (
	"dusk"
	"dusk/http/request"
	"dusk/http/response"
	"dusk/util"
	"dusk/validator"
	"fmt"
	"net/http"
	"path/filepath"
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
	if errMap != nil {
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

	file, err := request.ReadFile(r, "cover", "image/")
	if err != nil {
		response.BadRequest(rw, r, err)
		return
	}

	path, err := s.fw.UploadCover(file, b.Title)
	if err != nil {
		response.InternalServerError(rw, r, err)
		return
	}

	b.Cover = s.fw.GetRelativePath(path)
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

	// TODO
	file, err := request.ReadFile(r, "format", "application/")
	if err != nil {
		response.BadRequest(rw, r, err)
		return
	}

	path, err := s.fw.UploadFile(file, b.Title, fmt.Sprintf("%s.%s", b.Title, "epub"))
	if err != nil {
		response.InternalServerError(rw, r, err)
		return
	}

	if filepath.Ext(path) == ".epub" {
		coverPath, err := s.fw.UploadCoverFromFile(path, b.Title)
		if err != nil {
			// TODO delete uploaded file on err
			response.InternalServerError(rw, r, err)
			return
		}
		b.Cover = s.fw.GetRelativePath(coverPath)
	}

	b.Formats = append(b.Formats, s.fw.GetRelativePath(path))

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
	if errMap != nil {
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

	response.OK(rw, r, nil)
}
