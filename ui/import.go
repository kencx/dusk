package ui

import (
	"net/http"

	"github.com/kencx/dusk/http/request"
	"github.com/kencx/dusk/ui/views"
)

func (s *Handler) importIndex(rw http.ResponseWriter, r *http.Request) {
	tab := r.URL.Query().Get("tab")

	// default tab
	if tab == "" {
		views.NewImportIndex(s.base, "search", nil).Render(rw, r)
		return
	}

	if request.IsHtmxRequest(r) {
		views.NewImportIndex(s.base, tab, nil).Render(rw, r)
		return
	}

	views.ImportTabs.Select(tab).Render(r.Context(), rw)
}
