FROM golang:1.25-alpine AS builder

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o /kite cmd/kite/main.go

FROM alpine:3
COPY --from=builder kite /bin/kite
ENTRYPOINT ["/bin/kite"]