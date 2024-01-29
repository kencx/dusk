package pages

import "net/http"

type View interface {
	Render(rw http.ResponseWriter, r *http.Request)
	RenderError(rw http.ResponseWriter, r *http.Request, err error)
}

type ViewModel struct {
	Message      string
	ErrorMessage string
}

func NewViewModel(err error) ViewModel {
	m := ViewModel{}
	if err != nil {
		m.ErrorMessage = "Something went wrong, please try again"
	}
	return m
}
