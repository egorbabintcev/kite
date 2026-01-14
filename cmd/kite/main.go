package main

import (
	"context"
	"kite/internal/application/resource"
	"kite/internal/infrastructure/cache"
	"kite/internal/infrastructure/registry"
	"kite/internal/infrastructure/repository"
	"kite/internal/interface/web"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	logger.Info("Starting kite...")

	cacheDir := "/var/lib/kite/cache/packages"

	registryUrl := os.Getenv("KITE_REGISTRY_URL")
	if registryUrl == "" {
		registryUrl = "https://registry.npmjs.org"
	}

	cache := cache.NewFS(cacheDir)
	registry := registry.NewHttpClient(registryUrl)
	resolutionRepo := repository.NewPackageRepository(registry)
	servingRepo := repository.NewPackageArtifactRepository(registry, cache)
	getResourceUC := resource.NewGetResourceUC(resolutionRepo, servingRepo)
	srv := web.NewServer(logger, getResourceUC)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan

		srv.Stop(context.Background())
	}()

	srv.Start(":8000")
}
