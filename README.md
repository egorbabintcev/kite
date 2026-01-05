# 🪁 Kite

A lightweight, self-hosted caching CDN for npm-compatible registries.

## Overview

Kite provides a simple HTTP API for fetching package files from an npm-compatible registry. It's mainly designed for use in online playgrounds and sandboxes (see [browser Import Map API](https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/script/type/importmap)).

<pre align="center">
http://localhost:8000/:package@:version/:file
</pre>

- `:package` - name of the package on npm (can be scoped)
- `:version` - version of the package 
- `:file` - path to a file in the package

For example:

- `/preact@10.26.4/dist/preact.min.js`
- `/react@18.3.1/umd/react.production.min.js`
- `/three@0.174.0/build/three.module.min.js`

You can also use any valid semver range or npm tag:

- `/preact@latest/dist/preact.min.js`
- `/react@^18/umd/react.production.min.js`

If you don't specify a version, the latest tag is used by default.

- `/preact/dist/preact.min.js`
- `/vue/dist/vue.esm-browser.prod.js`

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
    image: egorbabintsev/kite:latest
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

## Inspiration

This project was heavily inspired by [unpkg](https://unpkg.com). While unpkg provides a fantastic public CDN, Kite focuses on bringing that same simplicity to your private infrastructure and self-hosted environments.
