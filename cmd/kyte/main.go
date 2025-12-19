package main

import (
	"context"
	"kyte/internal/interface/web"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	logger.Info("Starting kyte...")

	srv := web.NewServer(logger)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan

		srv.Stop(context.Background())
	}()

	srv.Start(":8000")
}
