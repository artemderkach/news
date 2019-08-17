FROM golang:1.12.8-alpine3.10 AS builder
# WORKDIR /go

WORKDIR /build/news
COPY . .

ENV CGO_ENABLED=0
ENV GOFLAGS "-mod=vendor"

RUN go test ./...
RUN go build 

FROM alpine:3.10

COPY --from=builder /build/news/news /usr/bin/news

CMD ["/usr/bin/news"]
