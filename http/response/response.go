package response

import (
	"dusk/util"
	"net/http"
)

var contentType = "application/json"

type Envelope map[string]interface{}

type response struct {
	rw         http.ResponseWriter
	r          *http.Request
	statusCode int
	headers    map[string]string
	body       []byte
}

func new(rw http.ResponseWriter, r *http.Request) *response {
	return &response{
		rw:         rw,
		r:          r,
		statusCode: http.StatusOK,
		headers:    map[string]string{"Content-Type": contentType},
	}
}

func OK(rw http.ResponseWriter, r *http.Request, body []byte) {
	res := new(rw, r)
	res.statusCode = http.StatusOK
	res.body = body
	res.write()
}

func NoContent(rw http.ResponseWriter, r *http.Request) {
	res := new(rw, r)
	res.statusCode = http.StatusNoContent
	res.write()
}

func Created(rw http.ResponseWriter, r *http.Request, body []byte) {
	res := new(rw, r)
	res.statusCode = http.StatusCreated
	res.body = body
	res.write()
}

func Custom(rw http.ResponseWriter, r *http.Request, code int, headers map[string]string, body []byte) {
	res := new(rw, r)
	res.statusCode = code
	res.headers = headers
	res.body = body
	res.write()
}

func newError(rw http.ResponseWriter, r *http.Request, err interface{}) *response {
	res := new(rw, r)
	res.statusCode = http.StatusBadRequest

	switch t := err.(type) {
	case error:
		res.body, err = util.ToJSON(Envelope{"error": t.Error()})
		if err != nil {
			res.body, err = util.ToJSON(Envelope{"error": "something went wrong"})
			res.statusCode = http.StatusInternalServerError
		}
	default:
		res.body, err = util.ToJSON(Envelope{"error": err})
		if err != nil {
			res.body, err = util.ToJSON(Envelope{"error": "something went wrong"})
			res.statusCode = http.StatusInternalServerError
		}
	}
	return res
}

func BadRequest(rw http.ResponseWriter, r *http.Request, err error) {
	res := newError(rw, r, err)
	res.write()
}

func InternalServerError(rw http.ResponseWriter, r *http.Request, err error) {
	res := newError(rw, r, err)
	res.statusCode = http.StatusInternalServerError
	res.write()
}

func NotFound(rw http.ResponseWriter, r *http.Request, err error) {
	res := newError(rw, r, err)
	res.statusCode = http.StatusNotFound
	res.write()
}

func Unauthorized(rw http.ResponseWriter, r *http.Request, err error) {
	res := newError(rw, r, err)
	res.statusCode = http.StatusUnauthorized
	res.headers["WWW-Authenticate"] = `Basic realm="Restricted"`
	res.write()
}

func ValidationError(rw http.ResponseWriter, r *http.Request, err map[string]string) {
	res := newError(rw, r, err)
	res.statusCode = http.StatusUnprocessableEntity
	res.write()
}

func (r *response) write() {
	for k, v := range r.headers {
		r.rw.Header().Set(k, v)
	}

	r.rw.WriteHeader(r.statusCode)
	r.rw.Write(r.body)
}
