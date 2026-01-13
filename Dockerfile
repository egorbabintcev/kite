FROM golang:1.25-alpine AS dev

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .

CMD ["go", "run", "./cmd/kite/main.go"]

# ---

FROM golang:1.25-alpine AS builder

WORKDIR /build

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /kite \
    ./cmd/kite/main.go

FROM alpine:3.20 AS prod

COPY --from=builder /kite /kite

VOLUME ["/var/lib/kite"]
EXPOSE 8000

ENTRYPOINT ["/kite"]
