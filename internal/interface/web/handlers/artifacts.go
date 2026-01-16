package handlers

import (
	"fmt"
	application "kite/internal/application/get_package_artifact"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
)

type GetPackageArtifactHandler struct {
	logger *slog.Logger
	uc     application.GetPackageArtifact
}

func NewGetPackageArtifactHandler(l *slog.Logger, uc application.GetPackageArtifact) *GetPackageArtifactHandler {
	return &GetPackageArtifactHandler{
		logger: l,
		uc:     uc,
	}
}

func (h *GetPackageArtifactHandler) Route(r chi.Router) {
	r.Get("/{name}", h.handle)
	r.Get("/{name}/*", h.handle)
	r.Get("/{name}@{version}", h.handle)
	r.Get("/{name}@{version}/*", h.handle)

	r.Get("/@{scope}/{name}", h.handle)
	r.Get("/@{scope}/{name}/*", h.handle)
	r.Get("/@{scope}/{name}@{version}", h.handle)
	r.Get("/@{scope}/{name}@{version}/*", h.handle)
}

func (h *GetPackageArtifactHandler) handle(w http.ResponseWriter, r *http.Request) {
	scope := chi.URLParam(r, "scope")
	name := chi.URLParam(r, "name")
	version := chi.URLParam(r, "version")
	path := chi.URLParam(r, "*")

	res, err := h.uc.Execute(r.Context(), application.GetPackageArtifactRequest{
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
