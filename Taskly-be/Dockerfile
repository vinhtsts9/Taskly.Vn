# Build
FROM golang:1.21-alpine AS builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ENV CGO_ENABLED=0 GOOS=linux
RUN go build -ldflags="-s -w" -o /app ./cmd/server

# Run (distroless static có sẵn CA)
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /app /app
USER nonroot:nonroot
ENTRYPOINT ["/app"]
