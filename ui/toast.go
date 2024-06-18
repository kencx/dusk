package ui

import (
	"net/http"

	"github.com/kencx/dusk/http/response"
	"github.com/kencx/dusk/ui/partials"
)

const (
	toastEvent = `{"onToast": ""}`
)

func SendToastMessage(rw http.ResponseWriter, r *http.Request, message string) {
	response.AddHxTriggerAfterSwap(rw, toastEvent)
	partials.ToastInfo(message, "", "").Render(r.Context(), rw)
}

func SendToastRawMessage(rw http.ResponseWriter, r *http.Request, rawMessage string) {
	response.AddHxTriggerAfterSwap(rw, toastEvent)
	partials.ToastRawInfo(rawMessage, "", "").Render(r.Context(), rw)
}
