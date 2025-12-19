package web

import (
	"log/slog"
	"net/http"
)

type handler struct {
	logger *slog.Logger
}

func newHandler(l *slog.Logger) *handler {
	return &handler{
		logger: l,
	}
}

func (h *handler) handleGetResource(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("Hello, world!"))
}
