package ui

import (
	"net/http"

	"github.com/kencx/dusk/ui/views"
)

func (s *Handler) importIndex(rw http.ResponseWriter, r *http.Request) {
	// handle htmx tabs
	if r.URL.Query().Has("tab") {
		tab := r.URL.Query().Get("tab")
		views.ImportTabs.Select(tab).Render(r.Context(), rw)
		return
	}
	views.NewImportIndex(s.baseView, "search", nil).Render(rw, r)
}
