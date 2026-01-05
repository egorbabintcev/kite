FROM golang:1.25-alpine AS dev

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .

CMD ["go", "run", "./cmd/kite/main.go"]
