package main

import (
	"context"
	application "kite/internal/application/get_package_artifact"
	"kite/internal/infrastructure/cache"
	"kite/internal/infrastructure/registry"
	"kite/internal/infrastructure/repository"
	"kite/internal/interface/web"
	"kite/internal/interface/web/handlers"
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

	packageRepo := repository.NewPackageRepository(registry)
	artifactRepo := repository.NewPackageArtifactRepository(registry, cache)

	getPackageArtifact := application.NewGetPackageArtifactUC(packageRepo, artifactRepo)

	srv := web.NewServer(logger,
		handlers.NewGetPackageArtifactHandler(logger, getPackageArtifact),
	)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan

		srv.Stop(context.Background())
	}()

	srv.Start(":8000")
}
