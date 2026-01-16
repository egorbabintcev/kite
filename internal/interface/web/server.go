package web

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	server   *http.Server
	logger   *slog.Logger
	handlers []RouteHandler
}

type RouteHandler interface {
	Route(r chi.Router)
}

func NewServer(l *slog.Logger, hs ...RouteHandler) *Server {
	l = l.With(slog.String("component", "web_server"))

	return &Server{
		logger:   l,
		handlers: hs,
	}
}

func (s *Server) Start(addr string) error {
	s.logger.Info(fmt.Sprintf("Starting http server at %s", addr))

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	for _, handler := range s.handlers {
		handler.Route(r)
	}

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
