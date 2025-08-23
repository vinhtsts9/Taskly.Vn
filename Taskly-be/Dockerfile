# ---------- Build stage ----------
FROM golang:alpine AS builder
WORKDIR /build

# cài git + chứng chỉ CA
RUN apk --no-cache add git ca-certificates

COPY . .

RUN go mod download
RUN go build -o Taskly.com ./cmd/server

# ---------- Runtime stage ----------
FROM scratch

# copy certs từ builder sang scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY ./sql/schema /sql/schema
COPY --from=builder /build/Taskly.com /

ENTRYPOINT ["/Taskly.com", "sql/schema"]
