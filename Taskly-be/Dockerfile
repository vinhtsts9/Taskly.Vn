FROM golang:alpine AS builder

WORKDIR /build

COPY . .

RUN go mod download

RUN go build -o social.network.com ./cmd/server

FROM scratch

COPY ./configs /configs

COPY --from=builder /build/social.network.com /

ENTRYPOINT [ "/social.network.com", "configs/local.yaml" ]