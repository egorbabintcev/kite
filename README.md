# 🪁 Kite

A lightweight, self-hosted caching CDN for npm-compatible registries.

## Overview

Kite provides a simple HTTP API for fetching package files from an npm-compatible registry. 

It's mainly designed for use in online playgrounds and sandboxes (see [browser Import Map API](https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/script/type/importmap)).

## Inspiration

This project was heavily inspired by [unpkg](https://unpkg.com). While unpkg provides a fantastic public CDN, Kite focuses on bringing that same simplicity to your private infrastructure and self-hosted environments.

## Quick start

**In Docker container**

```sh
docker run -d \
  --name yourapp \
  -p 8000:8000 \
  egorbabintsev/kite:latest
```

**Or with Docker Compose**

```yaml
services:
  kite:
    container_name: kite
    ports:
      - "8000:8000"
```

```sh
docker-compose up -d
```

Kite API is now available on http://localhost:8000

## Configuration

You can configure Kite using following env variables:

- `KITE_REGISTRY_URL` - base URL for npm-compatible registry, that will be used for downloading package archives. Defaults to `https://registry.npmjs.com`