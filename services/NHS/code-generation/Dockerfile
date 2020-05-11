FROM golang:buster AS builder

RUN mkdir /build
COPY . /build

RUN apt-get update && \
    apt-get -y install protobuf-compiler && \
    go get github.com/julienschmidt/httprouter && \
    go get github.com/dgrijalva/jwt-go && \
    go get github.com/spf13/viper && \
    cd /build && \
    make

FROM alpine

RUN apk add --no-cache libc6-compat

COPY ec256-key /ec256-key
COPY config.ini /service/etc/exposed/config.ini
COPY --from=builder /build/authcode /service/bin/authcode

EXPOSE 8080/tcp
ENTRYPOINT ["/service/bin/authcode"]
