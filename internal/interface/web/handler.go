package web

import (
	"fmt"
	"kite/internal/core"
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-chi/chi/v5"
)

type handler struct {
	logger  *slog.Logger
	service *core.Service
}

func newHandler(l *slog.Logger, s *core.Service) *handler {
	return &handler{
		logger:  l,
		service: s,
	}
}

func (h *handler) handleGetResource(w http.ResponseWriter, r *http.Request) {
	scope := chi.URLParam(r, "scope")
	name := chi.URLParam(r, "name")
	version := chi.URLParam(r, "version")
	path := chi.URLParam(r, "*")

	res, err := h.service.GetResource(r.Context(), scope, name, version, path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")

	if res.Redirect != nil {
		url := strings.Replace(r.URL.String(), url.PathEscape(version), url.PathEscape(res.Redirect.Version), 1)
		fmt.Printf("%s", url)
		http.Redirect(w, r, url, http.StatusMovedPermanently)
		return
	}

	http.ServeContent(w, r, res.Serve.Name, res.Serve.ModTime, res.Serve.Stream)
}
