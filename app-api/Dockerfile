FROM golang:1.15.6-alpine3.12 AS build-env

ADD . /go/src/github.com/hyorimits/hello-ebpf/app-api
WORKDIR /go/src/github.com/hyorimits/hello-ebpf/app-api

RUN go build .

FROM alpine:3.12

ENV GOPATH=/go
COPY --from=build-env /go/src/github.com/hyorimits/hello-ebpf/app-api/app-api /app-api

ENTRYPOINT [ "./app-api","--host", "0.0.0.0", "--port", "8080", "--scheme", "http"]
