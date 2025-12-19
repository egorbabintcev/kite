package main

import (
	"context"
	"kyte/internal/core"
	"kyte/internal/infrastructure/cache"
	"kyte/internal/infrastructure/registry"
	"kyte/internal/interface/web"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	logger.Info("Starting kyte...")

	cacheDir := os.Getenv("KYTE_CACHE_DIR")
	if cacheDir == "" {
		cacheDir = "/opt/kyte/cache/packages"
	}

	registryUrl := os.Getenv("KYTE_REGISTRY_URL")
	if registryUrl == "" {
		registryUrl = "https://registry.npmjs.org"
	}

	cache := cache.NewFS(cacheDir)
	registry := registry.NewHttpClient(registryUrl)
	core := core.NewService(logger, cache, registry)
	srv := web.NewServer(logger, core)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan

		srv.Stop(context.Background())
	}()

	srv.Start(":8000")
}
