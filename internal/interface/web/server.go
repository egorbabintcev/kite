package web

import (
	"context"
	"fmt"
	"kyte/internal/core"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	server  *http.Server
	logger  *slog.Logger
	service *core.Service
}

func NewServer(l *slog.Logger, s *core.Service) *Server {
	l = l.With(slog.String("component", "web_server"))

	return &Server{
		logger:  l,
		service: s,
	}
}

func (s *Server) Start(addr string) error {
	s.logger.Info(fmt.Sprintf("Starting http server at %s", addr))

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	h := newHandler(s.logger, s.service)

	r.Get("/{name}", h.handleGetResource)
	r.Get("/{name}/*", h.handleGetResource)
	r.Get("/{name}@{version}", h.handleGetResource)
	r.Get("/{name}@{version}/*", h.handleGetResource)

	r.Get("/@{scope}/{name}", h.handleGetResource)
	r.Get("/@{scope}/{name}/*", h.handleGetResource)
	r.Get("/@{scope}/{name}@{version}", h.handleGetResource)
	r.Get("/@{scope}/{name}@{version}/*", h.handleGetResource)

	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	s.server = srv

	return srv.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) {
	if s.server == nil {
		s.logger.Error("Nothing to stop, http server was not started")
		return
	}

	s.logger.Info("Gracefully shutting down http server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		s.logger.Error(fmt.Sprintf("Error gracefully shutting down http server: %v", err))
	}

	if err := s.server.Close(); err != nil {
		s.logger.Error(fmt.Sprintf("Error forcibly shutting down http server: %s", err))
	}
}
