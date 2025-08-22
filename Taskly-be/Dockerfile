FROM golang:alpine AS builder

WORKDIR /build

COPY . .

RUN go mod download

RUN go build -o Taskly.com ./cmd/server

FROM scratch

COPY ./sql/schema /sql/schema

COPY --from=builder /build/Taskly.com /

ENTRYPOINT [ "/Taskly.com", "sql/schema" ]