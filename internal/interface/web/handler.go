package web

import (
	"fmt"
	"kite/internal/application/resource"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
)

type handler struct {
	logger        *slog.Logger
	getResourceUC *resource.GetResourceUC
}

func newHandler(l *slog.Logger, uc *resource.GetResourceUC) *handler {
	return &handler{
		logger:        l,
		getResourceUC: uc,
	}
}

func (h *handler) handleGetResource(w http.ResponseWriter, r *http.Request) {
	scope := chi.URLParam(r, "scope")
	name := chi.URLParam(r, "name")
	version := chi.URLParam(r, "version")
	path := chi.URLParam(r, "*")

	res, err := h.getResourceUC.Execute(r.Context(), resource.GetResourceUCRequest{
		Name:         name,
		Scope:        scope,
		VersionQuery: version,
		Path:         path,
	})
	if err != nil {
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")

	if res.Redirect != nil {
		url := fmt.Sprintf("/%s@%s/%s", url.PathEscape(res.Redirect.Name), url.PathEscape(res.Redirect.Version), res.Redirect.Path)

		if res.Redirect.Scope != "" {
			url = fmt.Sprintf("/@%s%s", res.Redirect.Scope, url)
		}

		w.Header().Set("Cache-Control", "public, max-age=60, stale-while-revalidate=30")

		http.Redirect(w, r, url, http.StatusFound)
		return
	}

	w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")

	http.ServeContent(w, r, res.Serve.Name, res.Serve.ModTime, res.Serve.Stream)
}
